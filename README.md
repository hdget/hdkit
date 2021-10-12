# hdkit 
The tool is help rapidly develop microservice with [hdsdk](github.com/hdget/hdsdk). This tool will help automatically create a microservice based project boilerplate and generating template codes of microservice.
Now it supports two kinds codes for different framework:
- [dapr](https://github.com/dapr/dapr)
- [gokit](https://github.com/go-kit/kit)

## Table of Contents
- [Requirements](#requirements)
- [Usage](#usage)
  - [1. create microservice project](#1-create-microserivce-project)
  - [2. compile protobuf file](#2-compile-protobuf-file)
  - [3. generate microservice code templates](#3-generate-microservice-code-templates)
- [FAQ](#faq)
- [Reference](#reference)

## Requirements

Following library or utitiliy are required:

- [hdsdk](https://github.com/hdget/hdsdk)
- [protoc](https://github.com/google/protobuf/releases) binary used compile protobuf files
- [gogo protobuf](https://github.com/gogo/protobuf) 3rd libray to compile more fast grpc stub files

The`protoc` can be downloaded from above url and extrat binary file into `<GOPATH>/bin` directory.

The`gogo protobuf` tools can be installed with below commands:
```
go get github.com/gogo/protobuf/proto
go get github.com/gogo/protobuf/protoc-gen-gogofaster
go get github.com/gogo/protobuf/gogoproto
```
Or if you already create a project as described in `Usage` part, you can execute `<project>/bin/install_compiler.bat` to install those tools.

## Usage

### 1. create microservice project
```
hdkit new  <project_name> -p <protobuf_filepath> -t <framework>

e,g:
hdkit new example -p example.proto -t dapr
hdkit new example -p example.proto -t gokit
```

> The framework is `dapr` by default, you can specify to new `gokit` project by `-t` option

- Above command will do following:
 - create project boilerplate
  ```
  - <project>
      - autogen
      - service
      - proto
      - bin
      - go.mod
  ```
 - create `<project>/bin` directory and output script files to this dir
  - gen_grpc.bat:  used to compile protobuf files in windows
  - install_compiler.bat: install related grpc tools in windows
 - if command `-p` option specified valid protobuf filepath
  - if succeed find protobuf file, copy the protobuf file to `<project>/proto` directory
  - if `protoc` and related tools does exist, it will try to compile protobuf files under `<project>/proto` directory and output to `<project>/autogen/pb` directory

### 2. compile protobuf file
- If specified protobuf file and compiling it successfully in step1, please ignore this part
- If not specified protobuf in step1, it can use `hdkit gen pb` command to process protobuf
```
hdkit gen pb <project_name> -p <protobuf_filepath>
```
- if specified `-p` option with valid protobuf filepath, it will copy the file to `<project>/proto` directory, and change into directory`<project>/bin` then invoke` gen_grpc.*` script to compile protobuf files. As usual, the compiled pb files will save into `<project>/autogen/pb` directory
- if not specified `-p` option, it will try to find all `*.proto` files under `<project>/proto` directory, then it will compile them and save result to `<project>/autogen/pb` directory

### 3. generate microservice code templates
```
hdkit gen <project> -t <framework>

e,g:
hdkit gen example
hdkit gen example
```

> The framework is `dapr` by default, you can specify to generate `gokit` codes by `-t` option,
> Also please make sure you new `gokit` project before.

- Firstly, it will try to parse `*.pb.go` files under `<project>/autogen/pb` directory one by one until it find a protobuf `service`
- Based on found `service`, it will generate following:
  - `<project>/main.go`
    
     main entry  

  - `<project>/service`
  
     The `service` implementation, please put all business logic here, all files under this directory will not be overwriten

  - if it is gokit based project, it will generate: `<project>/autogen/grpc`
      - `endpoint_<method>.go`
      
          each method in `service` will have a corresponding endpoint, each endpoint implements `hdsdk` `GrpcEndpoint` interface
      
      - `handlers.go`: 
      
         handlers is grpc server handler collection, which will have `New` function and `ServeGrpc` function

## FAQ

to be added

## Reference

- [kit](https://github.com/GrantZheng/kit)
- [truss](https://github.com/metaverse/truss)