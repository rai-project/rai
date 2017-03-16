package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/rai_config.yml": {
		local:   "../rai_config.yml",
		size:    986,
		modtime: 1489706740,
		compressed: `
H4sIAAAAAAAA/1xTSXOrPBC88yt0dzkOeMFQ5YNX4jwbHC94uVACZBCbhBaw/eu/wkm+l5fjTGt6unpa
kFJTAaCAOTIBg1gBoELMJxyZ4AozjhQAQuTL6LuENW8GYBAgzr0U3T0cmmA0Gs93XW00atmbxMb1VOzt
an2IjJ7T2biWswlUezC/TW766v7o7q2tXJbutVyQs9za80yN0ks9UgDgKGBIeH/Zf1CT4rQtqSsy5J+O
rdDuQ1El9GLJu/02pI90PNPFTQ/xNDFYMbmfjNbkls13ropiQvSZPhFkNp1d9a5cnd3yNPuoQqfsv5OV
MR4vwiVTAGAowqQwgeRtBLloq40kwRs9RBbCBKqmqd1XTX8d9gZfGCPZ/9Z9114oGRSYFB5HASlCbgK1
nytcEIYa9ygjFQ4RMwHvKgD4kCNPsswEsRDU7HR49wXm8EEKWPOXgORPwzMTUOlnOGgzBEMlyDAqxO/j
SZoRGHq+DFIkTHDFGeIvDOI2ZSRBgfhi+8YlRyyEAj5bOAu9ZuDJ5T1rBUoR/6u46bw2qSA5xMXz7Tf3
E/ra8Cnv33D0jXk0aMWL2TUl6DJNP6rdgbzyaMrfAqlOjDWxorGd1vfHexkJ9/rB9KU1wOcWrK0h9D82
h7yAcFy+s/PVFzfn75rP3PxYpVEnZnGXbMe8RQaxetHTYN2rByu6SFYDvZVo1WS/tujOUBMxcXpa/cdK
0uVhPxXlQabru/an03fzx0rV+udAbooau/4F16c9y9nQ6UTH4Vtn5ljz4zRZ9sMygG/xkq8X7vy471u2
X6udnNz8Hl/0RgqVPpd+YyMqQkpwIZ5/CIA2+IR+X8gcdHWjcR1yXhP208NdcW6djHu1p5W7PZJWXJbq
JJF0WFXnuZgdtgF3tHo0Uv4LAAD//6gRXSvaAwAA
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
