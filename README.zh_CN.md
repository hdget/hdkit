# hdkit 
写这个工具的初衷是为了方便使用[hdget sdk](github.com/hdget/sdk)来快速开发微服务应用，该工具可以快速的创建项目蓝图，并根据protobuf文件定义快速生成Go程序代码。

## 1. 创建微服务项目
```
hdkit new <project_name> -p <protobuf_filepath>
```
- 以上命令将会创建微服务项目目录，该命令做了以下动作：
 - 创建项目目录结构
  ```
  - <project>
      - autogen
      - service
      - proto
      - bin
      - go.mod
  ```
 - 创建`<project>/bin`目录并将相关脚本文件输出保存到该目录
  - gen_grpc.bat:  windows下通过`protoc`编译grpc的stub
  - install_compiler.bat: windows下安装grpc编译工具的脚本
 - 如果`-p`参数指定了合法的protobuf文件的路径：
  - 如果能成功找到该路径，将该protobuf文件拷贝到`<project>/proto`目录中去
  - 如果`protoc`和相关的编译工具都存在，会尝试编译`<project>/proto`目录下的protobuf文件并输出保存到`<project>/autogen/pb`目录下

## 二、编译protobuf
- 如果第一步创建微服务项目时指定了protobuf文件并编译成功，请忽略该步骤
- 如果第一步创建微服务项目时未指定protobuf文件，可以通过`gen pb`命令来处理protobuf文件如下：
```
hdkit gen pb <project_name> -p <protobuf_filepath>
```
- 如果`-p`参数指定的protobuf文件路径正确，该命令会将其拷贝到`<project>/proto`目录，并进入`<project>/bin`目录去执行`gen_grpc`脚本去编译protobuf文件，编译后的pb文件输出保存在`<project>/autogen/pb`目录
- 如果未指定`-p`参数，则尝试在`<project>/proto`目录下去查找所有`*.proto`文件，扎到后会尝试编译并保存到`<project>/autogen/pb`目录

## 三、生成微服务代码
```
hdkit gen service <project>
```
- 首先其会在`<project>/autogen/pb`目录下查找编译了的`*.pb.go`文件，尝试从里面找到第一个`service`接口
- 然后根据找到的`service`接口，生成如下内容：
- `<project>/service`: 该目录时grpc接口的实现目录，工具自动生成了服务接口实现结构和对应的方法模板，后续需要在该目录下实现所有业务逻辑
- `<project>/autogen/grpc`: 该目录下保存了所有实现`hdget/sdk`下的`microservice`接口的文件，包括：
- `handlers.go`: 所有服务endpoint的集合struct,以及对应的New函数
- `endpoint_method.go`: 所有服务方法的对应method的实现

## 参考

- [kit](https://github.com/GrantZheng/kit)
- [truss](https://github.com/metaverse/truss)