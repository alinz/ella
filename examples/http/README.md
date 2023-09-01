# Setup

compile the Ella code to golang

```
ella gen main ./examples/http/*.ella ./examples/http/api.gen.go
```

then run it using the following command

```
go run examples/http/*.go
```

And call the following Curl to get your

```
curl -X POST http://localhost:8080/ella/http/Greeting/Hello -H 'Content-Type: application/json' -d '{"name": "Ella"}'
```
