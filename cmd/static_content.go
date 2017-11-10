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
		size:    1274,
		modtime: 1510329373,
		compressed: `
H4sIAAAAAAAA/2xTSZOiTBC98yvqbtAKKAgRHnCj7XFfaPVCFFAiaxW1gPrrv8Bepme+OWZl1Hsv38uE
hFgSAAXMkQUojCUAKkR9zJAFLjBjSAIgRL6IvkpYs+YDDALEmJeiuxeHFhgM7MlOUweD1nKdLON6xPfL
anGIzO6qvXad1TpQlvrkNrwZ8/tD2ztbMSvdSznFJ7FdTjIlSs/1QAKAoYAi7v1G/wGNi+O2JC7PkH98
b4XLHuRVQs6OuC9f++SR2mOD34wwHiUmLYb3o9ka3rLJzlXQFWNjbAw5Ho/GF0MT85NbHsebKlyVvTc8
N217Gs6oBABFUYwLCwgmI8i4rDSSOGv0YFFwCyiqqmgd1ej0u/pnj+Ls27qv2gsFhTzGhcdQgIuQWUDp
5RLjmKLGPUJxFYeIWoBpEgA+ZMgTNLPAlXNitdtMe4E5fOAC1uwlwPnT8MwCRPhZHMgUwVAKshgV/O/w
BMkwDD1fBCniFrjEGWIvFMYyoThBAf9E++oLhmgIOXw+xVnoNR+eWN6zbmYSfh5zj6JSxBTlqODPBQBA
BhQRTPkLCS8SAAn2vVIggbxvPR7MQ73biC9gdudx8BnpwZaVjtlTDaPTlxWpEdBY0MCiIiQ4/sGBAlXu
dWVFN2SlY8ia3kxABEey8j+Pmmm+yeUPgxq3IWM1pj/XdL1Xhlt38VgkHfO2aV/N1bnqHGpcEkez51pn
fntbr3uXJCq6u/3m4ji2Xc4d+8QJ2adcp1p0rob3dCBBwa9/Jtq8dJqrwTmMi6eUL++frU+tH+r+PJ6e
OYn01nU6vqQYnUfpptodcIdFI/YaCGVoLrAT2cu0vj/eyoi7lw01Zo4en1qwdvrQ36wPeQGhXb7R08Xn
t9Vvmo+7+kGlktWVXjW8tVkL61flbKTBolvrczJN5rrRStRquF84ZGcqCR+uumr9y0nS2WE/4uVBpIu7
+qvdc/PHXFF7p0Csizp2/XNcH/c0p/1VO3rvv7bHK2fyPkpmvbAM4Ot1xhZTd/K+7zlLv1baOb75XTbt
DiQifCb8f8b/0fp7gy1dM8x/J7srTq2jea/2pHK377h1LUtlmAjSr6rThI8P24Ct1HowkP4LAAD//76C
8Qv6BAAA
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
