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

### Install Golang Tools (Optional)

>> It's ok if some of these do not install on your system

```
go get github.com/klauspost/asmfmt/cmd/asmfmt && go install github.com/klauspost/asmfmt/cmd/asmfmt
go get github.com/derekparker/delve/cmd/dlv && go install github.com/derekparker/delve/cmd/dlv
go get github.com/kisielk/errcheck && go install github.com/kisielk/errcheck
go get github.com/davidrjenni/reftools/cmd/fillstruct && go install github.com/davidrjenni/reftools/cmd/fillstruct
go get github.com/mdempsky/gocode && go install github.com/mdempsky/gocode
go get github.com/rogpeppe/godef && go install github.com/rogpeppe/godef
go get github.com/zmb3/gogetdoc && go install github.com/zmb3/gogetdoc
go get golang.org/x/tools/cmd/goimports && go install golang.org/x/tools/cmd/goimports
go get github.com/golang/lint/golint && go install github.com/golang/lint/golint
go get github.com/alecthomas/gometalinter && go install github.com/alecthomas/gometalinter
go get github.com/fatih/gomodifytags && go install github.com/fatih/gomodifytags
go get golang.org/x/tools/cmd/gorename && go install golang.org/x/tools/cmd/gorename
go get github.com/jstemmer/gotags && go install github.com/jstemmer/gotags
go get golang.org/x/tools/cmd/guru && go install golang.org/x/tools/cmd/guru
go get github.com/josharian/impl && go install github.com/josharian/impl
go get honnef.co/go/tools/cmd/keyify && go install honnef.co/go/tools/cmd/keyify
go get github.com/fatih/motion && go install github.com/fatih/motion
go get github.com/koron/iferr && go install github.com/koron/iferr
go get github.com/maruel/panicparse/cmd/pp && go install github.com/maruel/panicparse/cmd/pp
go get github.com/maruel/panicparse && go install go get github.com/maruel/panicparse
go get github.com/DATA-DOG/goup && go install github.com/DATA-DOG/goup
```

### Install Dependency Manager

```
go get github.com/Masterminds/glide && go install github.com/Masterminds/glide
go get github.com/golang/dep && go install github.com/golang/dep
```

### Install linter

```
go get github.com/alecthomas/gometalinter
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

## GO Resources


* [50 Shades of Go](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/) - Traps, Gotchas, and Common Mistakes for New Golang Devs.
* [A Tour of Go](http://tour.golang.org/) - Interactive tour of Go.
* [Build web application with Golang](https://github.com/astaxie/build-web-application-with-golang) - Golang ebook intro how to build a web app with golang.
* [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin) - Get familiar with Gin and find out how it can help you reduce boilerplate code and build a request handling pipeline.
* [Games With Go](http://gameswithgo.org/) - A video series teaching programming and game development.
* [Go By Example](https://gobyexample.com/) - Hands-on introduction to Go using annotated example programs.
* [Go Cheat Sheet](https://github.com/a8m/go-lang-cheat-sheet) - Go's reference card.
* [Go database/sql tutorial](http://go-database-sql.org/) - Introduction to database/sql.
* [Golangbot](https://golangbot.com/learn-golang-series/) - Tutorials to get started with programming in Go.
* [Hackr.io](https://hackr.io/tutorials/learn-golang) - Learn Go from the best online golang tutorials submitted & voted by the golang programming community.
* [How to Use Godog for Behavior-driven Development in Go](https://semaphoreci.com/community/tutorials/how-to-use-godog-for-behavior-driven-development-in-go) - Get started with Godog — a Behavior-driven development framework for building and testing Go applications.
* [Learn Go with TDD](https://github.com/quii/learn-go-with-tests) - Learn Go with test-driven development.
* [Working with Go](https://github.com/mkaz/working-with-go) - Intro to go for experienced programmers.
* [Your basic Go](http://yourbasic.org/golang) - Huge collection of tutorials and how to's
* [Awesome Go](https://github.com/avelino/awesome-go)
