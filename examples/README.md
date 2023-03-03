run the following command to generate code all the examples

first install the `ella-gen`

```bash
go install github.com/alinz/ella.to/cmd/ella-gen
```

then the following will generate code for all examples

```bash
ella-gen rpc -i ./examples/constants -o ./examples/constants
ella-gen rpc -i ./examples/enums -o ./examples/enums
ella-gen rpc -i ./examples/messages -o ./examples/messages
ella-gen rpc -i ./examples/services -o ./examples/services
```
