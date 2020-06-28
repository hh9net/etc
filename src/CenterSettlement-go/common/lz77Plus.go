package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
)

const BLOCKSIZE = 64 * 1024

type encoder struct {
	data    []byte
	indexed *heap
	offset  int
}

func newEncoder() *encoder {
	return &encoder{
		data:    make([]byte, BLOCKSIZE+16),
		indexed: newHeap(),
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

	e.indexed.Reset()
	indexed := e.indexed

	var prev []byte
	for i := 0; i < last && e.offset/8 <= total; i++ {
		b := i + 2
		k := source[i:b]
		length, offset := indexed.Seek(source, i, k)

		if length == 0 {
			// 输出单个非匹配字符 0(1bit) + char(8bit)
			e.appendBit(0)
			e.appendWord(uint16(source[i]), 8)
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
				b = i + 2
				k = source[i:b]
			}
		}
		if len(prev) != 0 {
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
	data       []int
	duplicated int
	size       int
	table      []int
}

func newHeap() *heap {
	return &heap{
		data:  make([]int, 0, 2*BLOCKSIZE),
		size:  1,
		table: make([]int, BLOCKSIZE),
	}
}

func (h *heap) Reset() {
	h.data = h.data[:0]
	h.duplicated = 0
	h.size = 1
	for i := range h.table {
		h.table[i] = 0
	}
}

func (h *heap) Seek(source []byte, start int, k []byte) (length, offset int) {
	i := int(k[0])*256 + int(k[1])
	sl := len(source)

	for j := h.table[i]; j > 0; {
		j -= 2
		a := 2 + start
		b := 2 + h.data[j]
		c := 2
		for a < sl && b < h.size && source[a] == source[b] {
			a++
			b++
			c++
		}
		if c > length {
			length = c
			offset = h.data[j]
		}

		j = h.data[j+1]
	}
	return
}

func (h *heap) Size() int {
	return h.size
}

func (h *heap) Store(k []byte, offset int) {
	h.size++

	if k[0] == k[1] {
		if h.duplicated+1 == offset {
			h.duplicated++
			return
		}
		h.duplicated = offset
	}

	i := int(k[0])*256 + int(k[1])
	h.data = append(h.data, offset, h.table[i])
	h.table[i] = len(h.data)
}

type orderer struct {
	*sync.WaitGroup

	C   chan task
	i   int
	out io.Writer
}

type task struct {
	id   int
	data []byte
}

func newOrderer(out io.Writer) *orderer {
	ord := &orderer{out: out}
	ord.C = make(chan task, runtime.NumCPU())
	ord.WaitGroup = &sync.WaitGroup{}
	return ord
}

func (ord *orderer) Start() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	for i := 0; i < runtime.NumCPU(); i++ {
		ord.Add(1)
		go ord.Go(cond)
	}
}

func (ord *orderer) Go(c *sync.Cond) {
	data := make([]byte, 4)
	enc := newEncoder()

	for task := range ord.C {
		n := len(task.data)
		enc.Reset()
		result := enc.Encode(task.data)

		binary.BigEndian.PutUint16(data[:2], uint16(n))
		if len(result) < n {
			binary.BigEndian.PutUint16(data[2:], uint16(len(result)))
		} else {
			copy(data[2:], data[:2])
			result = task.data
		}

		c.L.Lock()
		for ord.i != task.id {
			c.Wait()
		}

		ord.out.Write(data)
		ord.out.Write(result)
		ord.i++
		c.Broadcast()
		c.L.Unlock()
	}
	ord.Done()
}

func CompressPlus(input io.Reader, output io.Writer) error {
	ord := newOrderer(output)
	ord.Start()

	var i int
	for {
		block := make([]byte, BLOCKSIZE)
		n, err := input.Read(block)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		ord.C <- task{i, block[:n]}
		i++
	}
	close(ord.C)
	ord.Wait()
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

func DecompressPlus(input io.Reader, output io.Writer) error {
	data := make([]byte, 2)
	dst := make([]byte, BLOCKSIZE)
	src := make([]byte, BLOCKSIZE)

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
	pwd := "./receivezipfile/" + fname //fname：zip文件 go run main.go
	origin, err := os.Open(pwd)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer origin.Close()
	fn := strings.Split(fname, ".lz77")
	outpwd := "./receivexml/" + fn[0] //xml文件 go run main.go
	out, cerr := os.Create(outpwd)
	if cerr != nil {
		log.Fatalln(cerr)
		return err
	}
	defer out.Close()
	if unzerr := DecompressPlus(origin, out); unzerr != nil {
		log.Fatalln(unzerr)
		return err
	}
	log.Printf("该%s压缩文件已经成功解压，解压为%s：", pwd, outpwd)
	return nil
}

func ZipLz77(fname string) error {
	pwd := "./generatexml/" + fname //go run main.go
	origin, oerr := os.Open(pwd)
	if oerr != nil {
		log.Fatalln(oerr)
		return oerr
	}
	defer origin.Close()
	outpwd := "./sendzipxml/" + fname + ".lz77" //go run  main.go
	out, cerr := os.Create(outpwd)
	if cerr != nil {
		log.Fatalln(cerr)
		return cerr
	}
	defer out.Close()

	if zerr := CompressPlus(origin, out); zerr != nil {
		log.Fatalln(zerr)
		return zerr
	}
	return nil
}
