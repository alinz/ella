```
┌─┐┬  ┬  ┌─┐
├┤ │  │  ├─┤
└─┘┴─┘┴─┘┴ ┴ v0.0.2
```

Ella, is yet another compiler to produce Go and Typescript code based on simple and easy-to-read schema IDL. There are many tools like gRPC, Twirp or event WebRPC to generate codes, but this little compiler is designed based on my views of 10+ years of developing backend and APIs. I wanted to simplify the tooling and produce almost perfect optimized, handcrafted code that can be read and understood, not like the output of gRPC.

Ella's schema went through several iterations to make it easier for extension and backward compatibility in future releases.

> **NOTE:**
>
> Ella's code generated has been used in couple of production projects, it has some designs that might not fit your needs, but I think it might solve a large number of projects. Also Ella's only emit `Go`, as a server and client, and `Typescript` as only client. This is intentinal as it servers my needs. However it can be extetended to produce other languages code by traversing the AST code.

# Installation

to install Ella's compiler, simply use the go install command

```bash
go install compiler.ella.to@latest
```

# Usage

Simplicity applies to the CLI command as well, it looks for all files that need to be compiled and outputs the result to the designated file. The extension of the output file tells the compiler whether you want to produce the typescript or golang code. That's pretty much of it.

For example, the following command, will generate `api.gen.go` in `/api` folder with the package name `api` and will read all the ella files inside `./schema` folder.

```bash
ella gen api /api/api.gen.go ./schema/*.ella
```

Also, we can format the schema as well to have a consistent look by running the following command

```bash
ella fmt ./schema/*.ella
```

The full CLI documentation can be accessed by running Ella command without any arguments

```
┌─┐┬  ┬  ┌─┐
├┤ │  │  ├─┤
└─┘┴─┘┴─┘┴ ┴ v0.0.2

Usage: ella [command]

Commands:
  - fmt Format one or many files in place using glob pattern
        ella fmt <glob path>

  - gen Generate code from a folder to a file and currently
        supports .go and .ts extensions
        ella gen <pkg> <output path to file> <search glob paths...>

  - ver Print the version of ella

example:
  ella fmt ./path/to/*.ella
  ella gen rpc ./path/to/output.go ./path/to/*.ella
  ella gen rpc ./path/to/output.ts ./path/to/*.ella ./path/to/other/*.ella
```

# Schema

There is a simple role when writing a schema in Ella's IDL, all identifier needs to be `Pascal Case`. This eliminates colliding with reserved keywords.

## const

A constant is a value that is immutable and can be used in field and method options.

### int

```
const BuildNumber = 10
```

### float

```
const Pi = 3.14
```

### string

```
const DoubleQuote = "1.0.0"

const SingleQuote = 'Hello World'

const MultiLine = `Hello
  This is really cool

  bye
`
```

### bool

```
const Debug = true
```

### byte size

This is a very helpful time as it makes it easier to write values in bytes type. The following postfixes can be used: `b`, `kb`, `mb`, `gb`, `tb`, `pb`, `eb`

> Note: the number has to be integer. No floating is premitted. In order to represent `1.1kb`, use the lower postfix `b` and represent the number as follows: `1100b`.

```
const MaxFileUploadSize = 100mb
```

### duration

> Note: `duration` type is also the same as `byte size` type. Only interger value is permitted.

This is a very helpful time as it makes it easier to write values in duration type. The following postfixes can be used: `ns`, `us`, `ms`, `s`, `m`, `h`

```
const MaxWaitTime = 5h
```

## enum

`enum` is a way to define a series of const values under the same category. In Golang, there is no such thing as `enum` and usually, people use a custom type and assign values to it. Ella's compiler does the heavy lifting of that and generates the most optimized version of the Go representation that supports both `yaml` and `json` marshal and unmarshal operations. It also supports ignoring value using `_` keyword.

> Note: Ella's enum type doesn't have a type, because behind the scene it will generate approporate type which is most memory efficient and optimize.

```
enum UserType {
  _
  Normal
  Guest
  Root
}

enum UserStatus {
  _
  Active = 10
  Deactive
  Deleted = 65
}
```

## model

model is a way to define a series of variables under the same category, similar to `struct` in Go.

```
model User {
  Firstname: string
  Lastname: string
  Age: int8
  LocationMap: []string
  Parents: []User
  ComplexMap: map<string, []User>
  CreatedAt: timestamp
}
```

> Note: Model's field type can be any of the default types such as `byte`, `bool`, `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `timestamp`, `string`, and complex types such as `map<key, value>` and array `[]type`, or it can be any other model's or enum's name type.

### field options

field options is a way to customize and assign values to each field of the model. Currently, there are the following predefined field options available.

- JsonOmitEmpty

is a boolean value that adds `omitempty` inside `Go` generated struct tag.

```
model User {
  Username: string
  Password: string {
    JsonOmitEmpty = true
  }
}
```

So basically the above model generates a `Go` struct as follows:

```golang
type User struct {
  Username string `json:"username"`
  Password string `json:"password,omitempty"`
}
```

- Json

Json option fields are a way to either not marshal and unmarshal the value or renaming the field during marshal and unmarshal operations.

```
model User {
  Firstname: string {
    Json = "firstName"
  }
}

model AnotherUser {
  Username: string
  Password: string {
    Json = false
  }
}
```

The above model will be converted to the following `Go` code

```Golang
type User struct {
  Firstname string `json:"firstName"`
}

type AnotherUser struct {
  Username string `json:"username"`
  Password string `json:"-"`
}
```

- Yaml

is the same as `Json`. Currently, this is the only option available for `Yaml`.

```
model User {
  Firstname: string {
    Yaml = "firstName"
  }
}

model User {
  Username: string
  Password: string {
    Yaml = false
  }
}
```

## service

### http

#### stream

#### file upload

### rpc

### method options

# References

- Here is the list of reserved keywords:

  - const
  - enum
  - model
  - http
  - rpc
  - service
  - byte
  - bool
  - int8
  - int16
  - int32
  - int64
  - uint8
  - uint16
  - uint32
  - uint64
  - float32
  - float64
  - timestamp
  - string
  - map
  - any
  - file
  - stream

- The logo was generated [here](https://patorjk.com/software/taag/#p=display&f=Calvin%20S&t=ella)

```

```
