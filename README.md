```
┌─┐┬  ┬  ┌─┐
├┤ │  │  ├─┤
└─┘┴─┘┴─┘┴ ┴ v0.0.2
```

Ella, is a yet another compiler to produce Go and Typescript code based on simple and easy to read schema or IDL. There are many tools like gRPC, Twirp or event Webrpc to generate codes, but this little compiler is designed based on my views of 10+ years of developing backend and apis. I wanted to simplify the tooling and produce most optimized, handcrafted code that can be read and understand not like the output of gRPC.

Ella's schema went through number of iterations to make it easier for extention and backward compatibility in the future releases.

> **NOTE:**
>
> Ella's code generated has been used in couple of production projects, it has some designs that might not fit your needs, but I think it might solve a large number of projects. Also Ella's only emit `Go`, as a server and client, and `Typescript` as only client. This is intentinal as it servers my needs. However it can be extetended to produce other languages code by traversing the AST code.

# Installation

in order to install ella's compiler, simply use go install command

```bash
go install ella.to/cmd/ella
```

# Usage

Simplicity applies to the cli command as well, it basically looks for all files that needs to be compile and output the result to the designated file. The extension of the output file tells the compiler weather you want to produce the typescript or golang code. That's pretty much of it.

```bash
ella gen api /api/api.gen.go ./schema/*.ella
```

Also we can format the schema as well to have a consitent look by running the following command

```bash
ella fmt ./schema/*.ella
```

The full cli documention can be access by running ella command without any arguments

```
Usage: ella [command]

Commands:
  - fmt Format one or many files in place using glob pattern
        ella fmt <glob path>

  - gen Generate code from a folder to a file and
        currently supports .go and .ts
        ella gen <pkg> <search glob path> <output path to file>

  - ver Print the version of ella

example:
  ella fmt ./path/to/*.ella
  ella gen rpc ./path/to/*.ella ./path/to/output.go
  ella gen rpc ./path/to/*.ella ./path/to/output.ts
```

# Features

# Design

# Schema

## const

creates a const values with varity of types. Constants can be used to store global variables for using inside model's and service method's options as a way of sharing common values.

supports the following types

- numbers, which includes int, float, duration like 1s or 5ns
- duration: 1s, 5ns
- bytesize, 10mb, 5gb
- bool, true and false
- string, single quote, double quotes and multiline using "`"

## enum

## model

## service

defines a skelton for each service

### methods

There are two types of methods available in each service.

#### http

#### rpc

# References

- Logo was generated [here](https://patorjk.com/software/taag/#p=display&f=Calvin%20S&t=ella)
