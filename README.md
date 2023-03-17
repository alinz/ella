# ella

`ella` is a command-line tool designed to accelerate the creation of REST endpoints in Golang and Typescript. While it's currently exclusively focused on Golang and Typescript, it has the potential to expand its support to other languages.

currently, it supports the following features

- [x] support comment
- [x] constant values
- [x] enums
  - [x] support only int or uint
- [x] messages
  - [x] support basic types such as int, uint, string and byte
  - [x] support complex types such as maps, array
- [x] services
  - [x] configure http method from GET or POST
  - [x] support multiple stream values using SSE
- [x] compatible with standard http.Handle
- [x] support both generating server and client code
- [x] generated code is readable and handcrafted

# Installation

install it using go install

```bash
go install ella.to/cmd/ella-gen@latest
```

# Schema

ella has its own IDL which is very familiar as it was borrowed from great projects such as gRPC and webrpc. It has some novel ideas as well which are presented in this document.

First, the extension of the IDL files should be `.ella` and run the following command to generate the output in golang:

```bash
ella-gen -i ./schema -o ./rpc/rpc.gen.go
```

or in typescript by changing the file ext from `.go` to `.ts`

```bash
ella-gen -i ./schema -o ./rpc/rpc.gen.ts
```

Basically the above command will generate rpc and read all `*.ella` files and generate a single file call `rpc.gen.go`.

Ella's IDL consist of 4 pars

- Types
- Constants
- Enums
- Message
- Service

## Types

There are a couple of common types supported in Ella IDL, here are the list of the

- int8, int16, int32 and int64
- uint8, uint16, uint32 and uint64
- byte
- any
- string
- map
- bool

```
map<int8, string>
```

- array

```
[]int8
```

- custom type using message keyword

## Constants

Constants are single immutable variable that helps to add additional information to generated code. There is only a mandatory constant named `Ella` which indicates the version of the IDL.

All constants' names should start with capital letters to indicate that they will be exported.

Here is an example:

```
Ella = "0.0.1"

MyStr1 = "Hello world!!"
MyStr2 = 'Hello world2!!'

MyNum1 = 1
MyNum2 = 2.4

MyBool = true
```

this will generate the following Golang code

```golang
const (
	Ella   = "0.0.1"
	MyBool = true
	MyNum1 = 1
	MyNum2 = 2.4
	MyStr1 = "Hello world!!"
	MyStr2 = "Hello world2!!"
)
```

## Enums

Enum doesn't need any explanation, here is an example of that

```
Ella = "0.0.1"

enum MyStatus int8 {
  Sleeping = 1
  Wake
  Working = 5
  Execsise
}
```

## Message

Message is a way to define complex and custom types. Each field in the message has optional feature. by default each field is exported and has a snake_case json tag.

For example

```
Ella = "0.0.1"

message Profile {
  Firstname: string {
    json = first_name
  }

  Lastname: string
}
```

the go representation of `Profile` is as follows:

```golang
type Profile struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"lastname"`
}
```

# Service

service is way to define both server and client stubs. Currently the following feature is supported

- [x] support multiple input and multiple output
- [x] support multiple streams out using multiplex SSE
- [x] support changing the either POST or GET for each method
- [-] support fileupload (pending)

for example:

```
Ella = "0.0.1"

service MyCustomService {
  Status(update: bool) => (lists: stream string)
}
```

the above code is converted to the following interface which can be used for both client and server

```golang
type MyCustomService interface {
	Status(ctx context.Context, update bool) (lists <-chan string, err error)
}
```
