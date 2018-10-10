# Developer Machine Setup

This section outlines how to install tools that would aid you in developing `rai`.

## Installing Go Compiler

RAI is developed using [golang](https://golang.org/) which needs to be installed for code to be compiled from source.
You can install Golang either through [Go Version Manager](https://github.com/moovweb/gvm)(recommended) or from the instructions on the [golang site](https://golang.org/). We recommend the Go Version Manager.

The following are instruction on how to install Go 1.11 through Go version manager.
Go version 1.11+ is required to compile RAI.

Download the [GVM](https://github.com/moovweb/gvm) using

```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```

Add the following line to your `.bashrc`(or `.zshrc` if using zsh) to set up the GVM environment.
This is sometimes done for you by default.

```
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
```

You can then install the Go 1.11 binary and set it as the default

```
gvm install go1.11 -B
gvm use go1.11 --default
```

`gvm` will setup both your `$GOPATH` and `$GOROOT` and you can validate that the installation completed by invoking

```sh
$ go env
GOARCH="amd64"
GOBIN=""
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="/home/abduld/.gvm/pkgsets/go1.11/global"
GORACE=""
GOROOT="/home/abduld/.gvm/gos/go1.11"
GOTOOLDIR="/home/abduld/.gvm/gos/go1.11/pkg/tool/linux_amd64"
GCCGO="gccgo"
CC="gcc"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build917072201=/tmp/go-build -gno-record-gcc-switches"
CXX="g++"
CGO_ENABLED="1"
PKG_CONFIG="pkg-config"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
```

## Installing Helper Tools

### Directory Movement

Install [Fasd](https://github.com/clvv/fasd)

### VimGo

Install [VimGo](https://github.com/fatih/vim-go)

### VSCode

Install [VSCode](https://code.visualstudio.com/) with [Golang](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go) extension.

### Install Golang linters

```
go install github.com/klauspost/asmfmt/cmd/asmfmt
go install github.com/derekparker/delve/cmd/dlv
go install github.com/kisielk/errcheck
go install github.com/davidrjenni/reftools/cmd/fillstruct
go install github.com/mdempsky/gocode
go install github.com/rogpeppe/godef
go install github.com/zmb3/gogetdoc
go install golang.org/x/tools/cmd/goimports
go install github.com/golang/lint/golint
go install github.com/alecthomas/gometalinter
go install github.com/fatih/gomodifytags
go install golang.org/x/tools/cmd/gorename
go install github.com/jstemmer/gotags
go install golang.org/x/tools/cmd/guru
go install github.com/josharian/impl
go install honnef.co/go/tools/cmd/keyify
go install github.com/fatih/motion
go install github.com/koron/iferr
```

### Install Dependency Manager

```
go install github.com/Masterminds/glide
go install github.com/golang/dep
```

### Install linter

```
go install github.com/alecthomas/gometalinter
gometalinter --install
```

## Downloading source code using `rai-srcmanager`

First, install the `rai-srcmanager` by

```sh
go get -u -v github.com/rai-project/rai-srcmanager
```

Download the required public repositories by

```sh
rai-srcmanager update --public
```

Now all the relevant repositories should now be in `$GOPATH/src/github.com/rai-project`.
