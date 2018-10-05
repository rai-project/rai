# RAI Client [![Travis Build Status](https://travis-ci.org/rai-project/rai.svg?branch=master)](https://travis-ci.org/rai-project/rai) [![AppVeyor Build status](https://ci.appveyor.com/api/projects/status/h6cmmpaejpo4at8m?svg=true)](https://ci.appveyor.com/project/abduld/rai)

> The RAI client is an executable downloaded by the students and runs on the students' machines.
> The executable requires no library dependencies and works on all the main operating systems and CPU architectures.
> Both features reduce the likelihood that students will have technical difficulties running the client.
> Students use the RAI client to interact with a distributed elastic system to submit jobs.

## Demo

[![asciicast](https://asciinema.org/a/6k5e96itnqu6ekbji60c3kgy4.png)](https://asciinema.org/a/6k5e96itnqu6ekbji60c3kgy4)

## Download Binaries

The code is continuously built and published. The client can be downloaded from the following URLs (depending on your OS and Architecture):

| Operating System | Architecture | Stable Version Link                                                             |
| ---------------- | ------------ | ------------------------------------------------------------------------------- |
| Linux            | amd64        | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-amd64.tar.gz)   |
| Linux            | armv5        | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-armv5.tar.gz)   |
| Linux            | armv6        | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-armv6.tar.gz)   |
| Linux            | armv7        | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-armv7.tar.gz)   |
| Linux            | arm64        | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-arm64.tar.gz)   |
| Linux            | ppc64        | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-ppc64.tar.gz)   |
| Linux            | ppc64le      | [URL](http://files.rai-project.com/dist/rai/stable/latest/linux-ppc64le.tar.gz) |
| OSX/Darwin       | amd64        | [URL](http://files.rai-project.com/dist/rai/stable/latest/darwin-amd64.tar.gz)  |
| Windows          | amd64        | [URL](http://files.rai-project.com/dist/rai/stable/latest/windows-amd64.tar.gz) |

## Building From Source

This is not recommended unless you are interested in developing and/or deploying `rai` on your personal cluster. To build from source simple run

```bash
go get -u github.com/rai-project/rai
```

You will need an extra secret key if you build from source.

-   Create a `.rai_config.yml` in the `rai` directory. You can copy the existing `rai_config.yml` as a starting point.
-   Run rai with `go run -tags develop main.go -d -v -s <app-secret> -p <project-folder>`

Alternatively, you can place the app secret in `~/.rai_secret` and just do

    go run -tags develop main.go -d -v -p <project-folder>

The `-tags develop` casuses rai to read the local configuration instead of using an embedded one.

## Usage

To run the client, use

```bash
rai -p <project folder>
```

From a user's point a view when the client runs, the local directory specified by `-p` gets uploaded to the server and extracted into the `/src` directory on the server. The server then executes the build commands from the `rai_build.yml` specification within the `/build` directory. Once the commands have been run, or there is an error, a zipped version of that `/build` directory is available from the server for download.

The server limits the task time to be an hour with a maximum of 8GB of memory being used within a session. The output `/build` directory is only available to be downloaded from the server for a short amount of time. Networking is also disabled on the execution server. Contact the teaching assistants if this is an issue.

#### Other Options

      -c, --color         Toggle color output.
      -d, --debug         Toggle debug mode.
      -p, --path string   Path to the directory you wish to submit. Defaults to the current working directory. (default "current working directory")
      -v, --verbose       Toggle verbose mode.

On Windows, it might be useful to disable the colored output. You can do that by using the `-c=false` option

## Setting your Profile

Each student will be contacted by a TA and given a secret key to use this service. Do not share your key with other users. The secret key is used to authenticate you with the server.

The `RAI_SECRET_KEY`, `RAI_TEAM_NAME`, and `RAI_ACCESS_KEY` should be specified in your `~/.rai_profile` (Linux/OSX) or `%HOME%/.rai_profile` (Windows -- for me this is `C:\Users\abduld\.rai_profile`) in the following way.

```yaml
profile:
  firstname: Abdul
  lastname: Dakkak
  username: abduld
  email: dakkak@illinois.edu
  access_key: XXXXXXXXXXXXXXXXXXX
  secret_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

## Project Build Specification

The `rai_build.yml` must exist in your project directory. In some cases, you may not be able to execute certain builtin bash commands, in this scenario the current workaround is to create a bash file and insert the commands you need to run. You can then execute the bash script within `rai_build.yml`.

The `rai_build.yml` is written as a [Yaml](http://yaml.org/) ([Spec](http://www.yaml.org/spec/1.2/spec.html)) file and has the following structure.

```yaml
rai:
  version: 0.2 # this is required
  image: nimbix/ubuntu-cuda-ppc64le:latest # nimbix/ubuntu-cuda-ppc64le:latest is a docker image
                                           # You can specify any image found on dockerhub
resources:
  cpu:
    architecture: ppc64le
  gpu:
    architecture: pascal
    count: 1 # tell the system that you're using a gpu
  network: false
commands:
  build:
    - echo "Building project"
    # Use CMake to generate the build files. Remember that your directory gets uploaded to /src
    - cmake /src
    # Run the make file to compile the project.
    - make
    # here we break the long command into multiple lines. The Yaml
    # format supports this using a block-strip command. See
    # http://stackoverflow.com/a/21699210/3543720 for info
    - >-
      ./mybinary -i input1,input2 -o output
```

Syntax errors will be reported, and the job will not be executed. You can check if your file is in a valid yaml format by using tools such as [Yaml Validator](http://codebeautify.org/yaml-validator).

## Building Docker Images

Most of the images on [Docker Hub](http://hub.docker.com) are compiled for X86 architectures. If you are using PPC64le, Power 8 architecture, e.g. Minsky, then you will have to build your Docker image from scratch. RAI has support for building Docker images on the host system.

1.  Create your Dockerfile we have created some example files that can be used as base and/or inspiration:
    [CUDNN](https://github.com/rai-project/Dockerfiles/tree/master/caffe2),
    [OpenCV](https://github.com/rai-project/Dockerfiles/tree/master/opencv),
    [CUMF](https://github.com/rai-project/Dockerfiles/tree/master/cumf),
    [NCCL](https://github.com/rai-project/Dockerfiles/tree/master/nccl),
    ...
    Refer to [Docker Syntax Reference](https://docs.docker.com/engine/reference/builder/) to understand the Docker commands.

2.  Tell the RAI client that you want to build a Dockerfile. This can be done by modifying the `.rai-build.yml` file to include the following:

```yaml
commands:
  build_image:
  image_name: your_user_name/your_image_name:your_image_version # example dakkak/cudnn:6.0
    dockerfile: "./Dockerfile" # the location of the Dockerfile on your local file system
  build: ...
```

3.  Run `rai` as if you are submitting the project. RAI will build and use the image you've specified.

A repository containing prebuilt Dockerfiles for PPC64le is available [Here](https://github.com/rai-project/Dockerfiles) and [Here](https://github.com/rai-project/Dockerfiles-ppc64le) and we accept contributions and/or fixes.

### Disabling Caching

By default, `rai` will not rebuild a docker image if it has the same name as a preexisting image on the system.
You can disable that by changing the `nocache` option to `true` in the `rai_build.yml` file.

```yaml
rai:
  version: 0.2
resources:
  cpu:
    architecture: ppc64le
  network: false
commands:
  build_image:
    image_name: rai/cumf:8.0
    dockerfile: "./Dockerfile"
    no_cache: true
```

### Publishing Docker Images

Docker images built using `rai` can be published on DockerHub.
You will have to explicitly tell `rai` to push the image in the `rai_build` file

```yaml
rai:
  version: 0.2
commands:
  build_image:
    image_name: c3sr/celery:4.0.2
    dockerfile: "./Dockerfile"
    push:
      push: true
```

#### Specifying DockerHub Credentials

There are two ways of specifying the DockerHub credentials. Through the `~/.rai_profile` file (prefered) by adding a `dockerhub` section e.g.

```yaml
profile:
  firstname: Abdul
  lastname: Dakkak
  ...
  dockerhub:
    username: dakkak
    password: ==AES32==PASS
```

or by placing it in the `rai_build.yml` file

```yaml
rai:
  version: 0.2
commands:
  build_image:
    image_name: c3sr/celery:4.0.2
    dockerfile: "./Dockerfile"
    push:
      push: true
      credentials:
        username: dakkak
        password: ==AES32==PASS
```

The password can be encrypted using the `rai encrypt` command.

## CUDA Profiling

Profiling can be performed using `nvprof`. Place the following build commands in your `rai_build.yml` file

```yaml
    - >-
      nvprof --cpu-profiling on --export-profile timeline.nvprof --
      ./mybinary -i input1,input2 -o output
    - >-
      nvprof --cpu-profiling on --export-profile analysis.nvprof --analysis-metrics --
      ./mybinary -i input1,input2 -o output
```

You could change the input and test datasets. This will output two files `timeline.nvprof` and `analysis.nvprof` which can be viewed using the `nvvp` tool (by performing a `file>import`). You will have to install the nvvp viewer on your machine to view these files.

_NOTE:_ `nvvp` will only show performance metrics for GPU invocations, so it may not show any analysis when you only have serial code.

## Benchmark

```
export GOTRACEBACK=all
go build -tags=bench
./rai bench --concurrency_count=10 --iteration_count=100 -s <<SECRET>> -p ./_fixtures/cuda_runtime |& panicparse
```


## Reporting Issues

Please use the [Github issue manager] to report any issues or suggestions.

Include the outputs of

```bash
rai version
```

as well as the output of

```bash
rai buildtime
```

In your bug report. You can also invoke the `rai` command with verbose and debug outputs using

```bash
rai --verbose --debug
```

## Usage

-   [PUMPS 2018 Summer School](https://github.com/illinois-impact/pumps-ai)
-   [ECE408 Fall2017 Project](https://github.com/webgpu/2017fa_ece408_project)
-   [ECE508 Spring2017 MP](https://github.com/webgpu/ece508-convlayer)
-   ECE508 Spring2017 Project
-   [ECE408 Fall2016 Project](https://github.com/webgpu/2016fa_ece408project)

## Publications

### Papers

-   [Dakkak, Abdul et al. “RAI: A Scalable Project Submission System for Parallel Programming Courses.” (2017).](docs/assets/rai-edupar17.pdf)

### Posters

-   [![GTC 2017](assets/gtc-rai-poster.png)](docs/assets/gtc-rai-poster.png)

### Presentations

#### GTC 2017

{::nomarkdown}

<iframe src="https://docs.google.com/presentation/d/1dnn0yJ_VJdJ108gv8TDMcuIoP6eUtoQMkGJ-zRvGZl0/embed?start=true&loop=false&delayms=3000" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>

{:/}

## License

NCSA/UIUC © [Abdul Dakkak](http://impact.crhc.illinois.edu/Content_Page.aspx?student_pg=Default-dakkak)

[github issue manager]: https://github.com/rai-project/rai/issues
