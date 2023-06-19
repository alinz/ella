```
ella = 1.0.0

const1 = example
const2 = this is an example
const3 = `this is an example
a multiline string
`

enum Enum1 uint8 {
  _ = 0
  ONE
  _
  THREE
}

message Message1 {
  firstname: string {
  }
}

service Service1 {
  http Ping() {
   http.method = GET
  }

  rpc Activity() => (result: Message1)
}

rpc Service2 {
}

```
