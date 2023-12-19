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

For example the following command, will generate `api.gen.go` in `/api` folder with the package name `api` and will read all the ella files inside `./schema` folder.

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

# References

- Logo was generated [here](https://patorjk.com/software/taag/#p=display&f=Calvin%20S&t=ella)
