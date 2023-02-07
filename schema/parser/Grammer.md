```
	// Byte
	// Bool
	// Any
	// Null
	// Uint8
	// Uint16
	// Uint32
	// Uint64
	// Int8
	// Int16
	// Int32
	// Int64
	// Float32
	// Float64
	// String
	// Timestamp

a = 10

enum A int64 {
	A = 1
}

message B {
	A: String {
		json.omitempty

	}
}

service C {

}


```

```
E -> T
E ->

T -> CONSTANT E
T -> ENUM E
T -> MESSAGE E
T -> SERVICE E

CONSTANT -> ID = VALUE

ENUM -> enum ID TYPE { VALUES }
VALUES -> ID VALUE
VALUE -> ID = VALUE
VALUE -> ID
VALUE ->

MESSAGE -> message ID { FIELDS }
FIELDS -> ID : TYPE { OPTIONS }
FIELDS -> ID : TYPE
FIELDS ->

SERVICE -> service ID { METHODS }
METHODS -> METHOD METHODS
METHODS -> METHOD

METHOD -> ID ( ARGS )
METHOD -> ID ( ARGS ) => ( RETURNS )

ARGS -> ARGS ARG
ARGS -> ARG

ARG -> ID : TYPE



E -> T && T
E -> T || T
E -> T

T -> ! id
T -> ! ( E )
T -> id {== != < > >= <=} V
T -> id

V -> C , V
V -> C

C -> string
C -> number
C -> byte_size
```
