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

	"/rai_config.yaml": {
		local:   "../rai_config.yaml",
		size:    800,
		modtime: 1489236601,
		compressed: `
H4sIAAAAAAAA/2ySSZO6PBCH73yK3C1HwQWl6n9wRWcU1FFcLlQIEcKWkEXUT/+WjFNzeY+dTj956teB
jFkaAAXMsQU4JBoAN8wDKrAFrjATWAMgxIGKfktYidcARAgL4af44ZPQAg1nkzikmsi9c1sfomHXbW08
290g3enP7uO7uXo8O3t7p5aldy3n9Kx2zizTo/RS/dMAEBhxLP0/qAVocdqVzJMZDk7HRuj0oLwl7GKr
h7MYsGc6mpryboZkkgx5MX6cho3xPZt9ezqOKTWn5ljS6WR6NTtqdfbK03R7C92y90lXw9FoHi65BgDH
EaGFBZRoYihkU3+ZSPHSoKqQFtANQ++0DbM96PbfPU6zd1BCUo5fUTBObyTE3AKiowEQQIF9xTMLxFIy
q9USnQ+YwyctYCU+EM3r9DILMBVkBDU5hqGGMoIL+aIpllEY+oFCKZYWuJIMiw8OSZNxmmAk34TfvhKY
h1DC+ohkof8aqAX9utb+7KCScfu1TppDUtR3fpl1603GBeIPJnHo/0jV++0NZ1G/Ec+n15TiyyTd3r4P
tC2iiVggpY+Ha2pHIyetHs/PMpLedcvNpd0n5was7AEMtptDXkA4Kj/5+RrIu/t/7/x8AgsYzI153KG7
kWjQfqxfzBStu1V/xebJqm82EuM23q9t9j3UEzl2u0b1ZSfp8rCfyPKg0vXD+Gr1vPy50o3eGalNUREv
uJDqtOc5H7it6DhYtKauPTtOkmUvLBFcxEuxnnuz475nO0Glt3J6D7pi3v2n/RcAAP//sGS6JiADAAA=
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
