run the following command to generate code all the examples

first install the `ella-gen`

```bash
go install ella.to/cmd/ella-gen
```

then the following will generate code for all examples

```bash
ella-gen -t rpc,workload -i ./examples/constants -o ./examples/constants
ella-gen -t rpc,workload -i ./examples/enums -o ./examples/enums
ella-gen -t rpc,workload -i ./examples/messages -o ./examples/messages
ella-gen -t rpc,workload -i ./examples/services -o ./examples/services
```
