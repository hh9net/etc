// Copyright 2017 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ini_test

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"testing"
)

func TestEmpty(t *testing.T) {
	Convey("Create an empty object", t, func() {
		f := ini.Empty()
		So(f, ShouldNotBeNil)

		// Should only have the default section
		So(len(f.Sections()), ShouldEqual, 1)

		// Default section should not contain any key
		So(len(f.Section("").Keys()), ShouldBeZeroValue)
	})
}

func TestFile_NewSection(t *testing.T) {
	Convey("Create a new section", t, func() {
		f := ini.Empty()
		So(f, ShouldNotBeNil)

		sec, err := f.NewSection("author")
		So(err, ShouldBeNil)
		So(sec, ShouldNotBeNil)
		So(sec.Name(), ShouldEqual, "author")

		So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "author"})

		Convey("With duplicated name", func() {
			sec, err := f.NewSection("author")
			So(err, ShouldBeNil)
			So(sec, ShouldNotBeNil)

			// Does nothing if section already exists
			So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "author"})
		})

		Convey("With empty string", func() {
			_, err := f.NewSection("")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestFile_NonUniqueSection(t *testing.T) {
	Convey("Read and write non-unique sections", t, func() {
		f, err := ini.LoadSources(ini.LoadOptions{
			AllowNonUniqueSections: true,
		}, []byte(`[Interface]
Address = 192.168.2.1
PrivateKey = <server's privatekey>
ListenPort = 51820

[Peer]
PublicKey = <client's publickey>
AllowedIPs = 192.168.2.2/32

[Peer]
PublicKey = <client2's publickey>
AllowedIPs = 192.168.2.3/32`))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		sec, err := f.NewSection("Peer")
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		_, _ = sec.NewKey("PublicKey", "<client3's publickey>")
		_, _ = sec.NewKey("AllowedIPs", "192.168.2.4/32")

		var buf bytes.Buffer
		_, err = f.WriteTo(&buf)
		So(err, ShouldBeNil)
		str := buf.String()
		So(str, ShouldEqual, `[Interface]
Address    = 192.168.2.1
PrivateKey = <server's privatekey>
ListenPort = 51820

[Peer]
PublicKey  = <client's publickey>
AllowedIPs = 192.168.2.2/32

[Peer]
PublicKey  = <client2's publickey>
AllowedIPs = 192.168.2.3/32

[Peer]
PublicKey  = <client3's publickey>
AllowedIPs = 192.168.2.4/32

`)
	})

	Convey("Delete non-unique section", t, func() {
		f, err := ini.LoadSources(ini.LoadOptions{
			AllowNonUniqueSections: true,
		}, []byte(`[Interface]
Address    = 192.168.2.1
PrivateKey = <server's privatekey>
ListenPort = 51820

[Peer]
PublicKey  = <client's publickey>
AllowedIPs = 192.168.2.2/32

[Peer]
PublicKey  = <client2's publickey>
AllowedIPs = 192.168.2.3/32

[Peer]
PublicKey  = <client3's publickey>
AllowedIPs = 192.168.2.4/32

`))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		err = f.DeleteSectionWithIndex("Peer", 1)
		So(err, ShouldBeNil)

		var buf bytes.Buffer
		_, err = f.WriteTo(&buf)
		So(err, ShouldBeNil)
		str := buf.String()
		So(str, ShouldEqual, `[Interface]
Address    = 192.168.2.1
PrivateKey = <server's privatekey>
ListenPort = 51820

[Peer]
PublicKey  = <client's publickey>
AllowedIPs = 192.168.2.2/32

[Peer]
PublicKey  = <client3's publickey>
AllowedIPs = 192.168.2.4/32

`)
	})

	Convey("Delete all sections", t, func() {
		f := ini.Empty(ini.LoadOptions{
			AllowNonUniqueSections: true,
		})
		So(f, ShouldNotBeNil)

		_ = f.NewSections("Interface", "Peer", "Peer")
		So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "Interface", "Peer", "Peer"})
		f.DeleteSection("Peer")
		So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "Interface"})
	})
}

func TestFile_NewRawSection(t *testing.T) {
	Convey("Create a new raw section", t, func() {
		f := ini.Empty()
		So(f, ShouldNotBeNil)

		sec, err := f.NewRawSection("comments", `1111111111111111111000000000000000001110000
111111111111111111100000000000111000000000`)
		So(err, ShouldBeNil)
		So(sec, ShouldNotBeNil)
		So(sec.Name(), ShouldEqual, "comments")

		So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "comments"})
		So(f.Section("comments").Body(), ShouldEqual, `1111111111111111111000000000000000001110000
111111111111111111100000000000111000000000`)

		Convey("With duplicated name", func() {
			sec, err := f.NewRawSection("comments", `1111111111111111111000000000000000001110000`)
			So(err, ShouldBeNil)
			So(sec, ShouldNotBeNil)
			So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "comments"})

			// Overwrite previous existed section
			So(f.Section("comments").Body(), ShouldEqual, `1111111111111111111000000000000000001110000`)
		})

		Convey("With empty string", func() {
			_, err := f.NewRawSection("", "")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestFile_NewSections(t *testing.T) {
	Convey("Create new sections", t, func() {
		f := ini.Empty()
		So(f, ShouldNotBeNil)

		So(f.NewSections("package", "author"), ShouldBeNil)
		So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "package", "author"})

		Convey("With duplicated name", func() {
			So(f.NewSections("author", "features"), ShouldBeNil)

			// Ignore section already exists
			So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "package", "author", "features"})
		})

		Convey("With empty string", func() {
			So(f.NewSections("", ""), ShouldNotBeNil)
		})
	})
}

func TestFile_GetSection(t *testing.T) {
	Convey("Get a section", t, func() {
		f, err := ini.Load(fullConf)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		sec, err := f.GetSection("author")
		So(err, ShouldBeNil)
		So(sec, ShouldNotBeNil)
		So(sec.Name(), ShouldEqual, "author")

		Convey("Section not exists", func() {
			_, err := f.GetSection("404")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestFile_Section(t *testing.T) {
	Convey("Get a section", t, func() {
		f, err := ini.Load(fullConf)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		sec := f.Section("author")
		So(sec, ShouldNotBeNil)
		So(sec.Name(), ShouldEqual, "author")

		Convey("Section not exists", func() {
			sec := f.Section("404")
			So(sec, ShouldNotBeNil)
			So(sec.Name(), ShouldEqual, "404")
		})
	})

	Convey("Get default section in lower case with insensitive load", t, func() {
		f, err := ini.InsensitiveLoad([]byte(`
[default]
NAME = ini
VERSION = v1`))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		So(f.Section("").Key("name").String(), ShouldEqual, "ini")
		So(f.Section("").Key("version").String(), ShouldEqual, "v1")
	})
}

func TestFile_Sections(t *testing.T) {
	Convey("Get all sections", t, func() {
		f, err := ini.Load(fullConf)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		secs := f.Sections()
		names := []string{ini.DefaultSection, "author", "package", "package.sub", "features", "types", "array", "note", "comments", "string escapes", "advance"}
		So(len(secs), ShouldEqual, len(names))
		for i, name := range names {
			So(secs[i].Name(), ShouldEqual, name)
		}
	})
}

func TestFile_ChildSections(t *testing.T) {
	Convey("Get child sections by parent name", t, func() {
		f, err := ini.Load([]byte(`
[node]
[node.biz1]
[node.biz2]
[node.biz3]
[node.bizN]
`))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		children := f.ChildSections("node")
		names := []string{"node.biz1", "node.biz2", "node.biz3", "node.bizN"}
		So(len(children), ShouldEqual, len(names))
		for i, name := range names {
			So(children[i].Name(), ShouldEqual, name)
		}
	})
}

func TestFile_SectionStrings(t *testing.T) {
	Convey("Get all section names", t, func() {
		f, err := ini.Load(fullConf)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		So(f.SectionStrings(), ShouldResemble, []string{ini.DefaultSection, "author", "package", "package.sub", "features", "types", "array", "note", "comments", "string escapes", "advance"})
	})
}

func TestFile_DeleteSection(t *testing.T) {
	Convey("Delete a section", t, func() {
		f := ini.Empty()
		So(f, ShouldNotBeNil)

		_ = f.NewSections("author", "package", "features")
		f.DeleteSection("features")
		f.DeleteSection("")
		So(f.SectionStrings(), ShouldResemble, []string{"author", "package"})
	})
}

func TestFile_Append(t *testing.T) {
	Convey("Append a data source", t, func() {
		f := ini.Empty()
		So(f, ShouldNotBeNil)

		So(f.Append(minimalConf, []byte(`
[author]
NAME = Unknwon`)), ShouldBeNil)

		Convey("With bad input", func() {
			So(f.Append(123), ShouldNotBeNil)
			So(f.Append(minimalConf, 123), ShouldNotBeNil)
		})
	})
}

func TestFile_WriteTo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping testing on Windows")
	}

	Convey("Write content to somewhere", t, func() {
		f, err := ini.Load(fullConf)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		f.Section("author").Comment = `Information about package author
# Bio can be written in multiple lines.`
		f.Section("author").Key("NAME").Comment = "This is author name"
		_, _ = f.Section("note").NewBooleanKey("boolean_key")
		_, _ = f.Section("note").NewKey("more", "notes")

		var buf bytes.Buffer
		_, err = f.WriteTo(&buf)
		So(err, ShouldBeNil)

		golden := "testdata/TestFile_WriteTo.golden"
		if *update {
			So(ioutil.WriteFile(golden, buf.Bytes(), 0644), ShouldBeNil)
		}

		expected, err := ioutil.ReadFile(golden)
		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, string(expected))
	})

	Convey("Support multiline comments", t, func() {
		f, err := ini.Load([]byte(`
# 
# general.domain
# 
# Domain name of XX system.
domain      = mydomain.com
`))
		So(err, ShouldBeNil)

		f.Section("").Key("test").Comment = "Multiline\nComment"

		var buf bytes.Buffer
		_, err = f.WriteTo(&buf)
		So(err, ShouldBeNil)

		So(buf.String(), ShouldEqual, `# 
# general.domain
# 
# Domain name of XX system.
domain = mydomain.com
; Multiline
; Comment
test   = 

`)

	})
}

func TestFile_SaveTo(t *testing.T) {
	Convey("Write content to somewhere", t, func() {
		f, err := ini.Load(fullConf)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		So(f.SaveTo("testdata/conf_out.ini"), ShouldBeNil)
		So(f.SaveToIndent("testdata/conf_out.ini", "\t"), ShouldBeNil)
	})
}

func TestFile_WriteToWithOutputDelimiter(t *testing.T) {
	Convey("Write content to somewhere using a custom output delimiter", t, func() {
		f, err := ini.LoadSources(ini.LoadOptions{
			KeyValueDelimiterOnWrite: "->",
		}, []byte(`[Others]
Cities = HangZhou|Boston
Visits = 1993-10-07T20:17:05Z, 1993-10-07T20:17:05Z
Years = 1993,1994
Numbers = 10010,10086
Ages = 18,19
Populations = 12345678,98765432
Coordinates = 192.168,10.11
Flags       = true,false
Note = Hello world!`))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		var actual bytes.Buffer
		var expected = []byte(`[Others]
Cities      -> HangZhou|Boston
Visits      -> 1993-10-07T20:17:05Z, 1993-10-07T20:17:05Z
Years       -> 1993,1994
Numbers     -> 10010,10086
Ages        -> 18,19
Populations -> 12345678,98765432
Coordinates -> 192.168,10.11
Flags       -> true,false
Note        -> Hello world!

`)
		_, err = f.WriteTo(&actual)
		So(err, ShouldBeNil)

		So(bytes.Equal(expected, actual.Bytes()), ShouldBeTrue)
	})
}

// Inspired by https://github.com/go-ini/ini/issues/207
func TestReloadAfterShadowLoad(t *testing.T) {
	Convey("Reload file after ShadowLoad", t, func() {
		f, err := ini.ShadowLoad([]byte(`
[slice]
v = 1
v = 2
v = 3
`))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)

		So(f.Section("slice").Key("v").ValueWithShadows(), ShouldResemble, []string{"1", "2", "3"})

		So(f.Reload(), ShouldBeNil)
		So(f.Section("slice").Key("v").ValueWithShadows(), ShouldResemble, []string{"1", "2", "3"})
	})
}
