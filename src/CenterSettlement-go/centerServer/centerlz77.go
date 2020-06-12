package centerServer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

const BLOCKSIZE = 64 * 1024

type encoder struct {
	data   []byte
	offset int
}

func newEncoder() *encoder {
	return &encoder{
		data: make([]byte, BLOCKSIZE+16),
	}
}

func (e *encoder) appendBit(val uint16) {
	idx, pos := e.offset/8, e.offset%8
	e.offset++
	if val > 0 {
		e.data[idx] |= (0x80 >> pos)
	}
}

func (e *encoder) appendWord(val uint16, bits int) {
	for i := 16 - bits; i < 16; i++ {
		e.appendBit(val & (0x8000 >> i))
	}
}

func (e *encoder) bytes() []byte {
	idx, pos := e.offset/8, e.offset%8
	if pos > 0 {
		idx++
	}
	return e.data[:idx]
}

func (e *encoder) Encode(source []byte) []byte {
	total := len(source)
	last := total - 1
	indexed := newHeap()

	var prev string
	for i := 0; i < last && e.offset/8 <= total; i++ {
		a := source[i]
		b := source[i+1]
		k := fmt.Sprintf("%c%c", a, b)
		length, offset := indexed.Seek(source, i, k)

		if length == 0 {
			// 输出单个非匹配字符 0(1bit) + char(8bit)
			e.appendBit(0)
			e.appendWord(uint16(a), 8)
		} else {
			// 输出匹配术语 flag(1bit) + len(γ编码) + offset(最大16bit)
			e.appendBit(1)
			gamma := length - 1
			q := int(math.Floor(math.Log2(float64(gamma))))
			// 输出q个1
			if q > 0 {
				e.appendWord(0xffff, q)
			}
			// 输出一个0
			e.appendBit(0)
			// 输出余数, q位
			if q > 0 {
				r := gamma - (1 << q)
				e.appendWord(uint16(r), q)
			}
			// offset
			// 在窗口不满64k大小时，不需要16位存储偏移
			bits := int(math.Ceil(math.Log2(float64(indexed.Size()))))
			e.appendWord(uint16(offset), bits)

			// scroll
			for j := 1; j < length; j++ {
				indexed.Store(prev, i-1)
				prev = k
				i++

				if i >= last {
					break
				}
				a = source[i]
				b = source[i+1]
				k = fmt.Sprintf("%c%c", a, b)
			}
		}
		if prev != "" {
			indexed.Store(prev, i-1)
		}
		prev = k
	}
	// 输出最后一个字符
	if indexed.Size() < total {
		e.appendBit(0)
		e.appendWord(uint16(source[last]), 8)
	}

	return e.bytes()
}

func (e *encoder) Reset() {
	for i := range e.data {
		e.data[i] = 0
	}
	e.offset = 0
}

type heap struct {
	data       map[string][]int
	duplicated int
	size       int
}

func newHeap() *heap {
	return &heap{
		data: make(map[string][]int),
		size: 1,
	}
}

func (h *heap) Seek(source []byte, start int, k string) (length, offset int) {
	list := h.data[k]
	tail := len(list)

	for tail > 0 {
		tail--

		a := start
		b := list[tail]
		var c int
		for a < len(source) && b < h.size && source[a] == source[b] {
			a++
			b++
			c++
		}
		if c > length {
			length = c
			offset = list[tail]
		}
	}
	return
}

func (h *heap) Size() int {
	return h.size
}

func (h *heap) Store(k string, offset int) {
	h.size++

	if k[0] == k[1] {
		if h.duplicated+1 == offset {
			h.duplicated++
			return
		}
		h.duplicated = offset
	}
	h.data[k] = append(h.data[k], offset)
}

func Compress(input io.Reader, output io.Writer) error {
	block := make([]byte, BLOCKSIZE)
	data := make([]byte, 2)
	enc := newEncoder()

	for {
		n, err := input.Read(block)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		binary.BigEndian.PutUint16(data, uint16(n))
		output.Write(data)

		enc.Reset()
		result := enc.Encode(block[:n])

		if len(result) < n {
			binary.BigEndian.PutUint16(data, uint16(len(result)))
			output.Write(data)
			output.Write(result)
		} else {
			output.Write(data)
			output.Write(block[:n])
		}
	}
	return nil
}

type decoder struct {
	data   []byte
	offset int
}

func newDecoder(block []byte) *decoder {
	return &decoder{
		data: block,
	}
}

func (d *decoder) readBit() (result uint16, n int) {
	idx, pos := d.offset/8, d.offset%8
	if idx >= len(d.data) {
		return
	}

	d.offset++
	if (d.data[idx] & (0x80 >> pos)) > 0 {
		result = 1
	}
	n = 1
	return
}

func (d *decoder) readWord(bits int) (uint16, int) {
	var result uint16
	for i := 0; i < bits; i++ {
		if i > 0 {
			result <<= 1
		}

		bit, n := d.readBit()
		if n == 0 {
			return result, i
		}
		result |= bit
	}
	return result, bits
}

func (d *decoder) Decode(dest []byte) error {
	b := bytes.NewBuffer(dest)
	b.Reset()

	for b.Len() < len(dest) {
		flag, n := d.readBit()
		if n < 1 {
			return errors.New("incomplete: flag")
		}

		if flag == 0 {
			w, n := d.readWord(8)
			if n < 8 {
				return errors.New("incomplete: normal")
			}
			b.WriteByte(byte(w))
		} else {
			var q int
			for {
				val, n := d.readBit()
				if n < 1 {
					return errors.New("incomplete: q")
				}
				if val == 0 {
					break
				}
				q++
			}

			length := 2
			if q > 0 {
				w, n := d.readWord(q)
				if n < q {
					return errors.New("incomplete: length")
				}
				length = 1 << q
				length += int(w)
				length += 1
			}

			// offset
			bits := int(math.Ceil(math.Log2(float64(b.Len()))))
			w, n := d.readWord(bits)
			if n < bits {
				return errors.New("incomplete: offset")
			}
			offset := int(w)
			b.Write(dest[offset : offset+length])
		}
	}
	return nil
}

func (d *decoder) Reset(block []byte) {
	d.data = block
	d.offset = 0
}

func Decompress(input io.Reader, output io.Writer) error {
	data := make([]byte, 2)
	src := make([]byte, BLOCKSIZE)
	dst := make([]byte, BLOCKSIZE)
	dec := newDecoder(nil)

	for {
		if _, err := input.Read(data); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("a: %w", err)
		}
		a := int(binary.BigEndian.Uint16(data))
		if a == 0 {
			a = BLOCKSIZE
		}

		if _, err := input.Read(data); err != nil {
			return fmt.Errorf("b: %w", err)
		}
		b := int(binary.BigEndian.Uint16(data))
		if b == 0 {
			b = BLOCKSIZE
		}

		block := src[:b]
		_, err := input.Read(block)
		if err != nil {
			return fmt.Errorf("block: %w", err)
		}
		if b < a {
			dec.Reset(block)
			block = dst[:a]
			if err := dec.Decode(block); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}
		}

		if _, err = output.Write(block); err != nil {
			return fmt.Errorf("writing: %w", err)
		}
	}
	return nil
}

func UnZipLz77(fname string) error {

	//pwd := "./receivezipfile/" + fname
	pwd := "./" + fname
	origin, err := os.Open(pwd)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer origin.Close()

	fn := strings.Split(fname, ".lz77")

	outpwd := "../centeryuanshi/" + fn[0]
	out, cerr := os.Create(outpwd)
	if cerr != nil {
		log.Fatalln(cerr)
		return err
	}
	defer out.Close()

	if unzerr := Decompress(origin, out); unzerr != nil {
		log.Fatalln(unzerr)
		return err
	}
	log.Println("原始交易消息解压成功")
	return nil
	//log.Println(Decompress(origin, os.Stdout))
}

func ZipLz77(fname string) error {
	pwd := "CenterSettlement-go/generatexml/" + fname
	origin, oerr := os.Open(pwd)
	if oerr != nil {
		log.Fatalln(oerr)
		return oerr
	}
	defer origin.Close()
	outpwd := "CenterSettlement-go/sendzipxml/" + fname + ".lz77"
	out, cerr := os.Create(outpwd)
	if cerr != nil {
		log.Fatalln(cerr)
		return cerr
	}
	defer out.Close()

	if zerr := Compress(origin, out); zerr != nil {
		log.Fatalln(zerr)
		return zerr
	}
	return nil
}
