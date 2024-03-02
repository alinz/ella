package token

const (
	Error                    Type = -1   // Error token type which indicates error
	EOF                      Type = iota // EOF token type which indicates end of input
	Identifier                           // identifier
	Const                                // const
	Enum                                 // enum
	Model                                // model
	Http                                 // http
	Rpc                                  // rpc
	Service                              // service
	Byte                                 // byte
	Bool                                 // bool
	Int8                                 // int8
	Int16                                // int16
	Int32                                // int32
	Int64                                // int64
	Uint8                                // uint8
	Uint16                               // uint16
	Uint32                               // uint32
	Uint64                               // uint64
	Float32                              // float32
	Float64                              // float64
	Timestamp                            // timestamp
	String                               // string
	Map                                  // map
	Array                                // array []
	Any                                  // any
	File                                 // file
	Stream                               // stream
	ConstDuration                        // 1ns, 1us, 1ms, 1s, 1m, 1h
	ConstBytes                           // 1b, 1kb, 1mb, 1gb, 1tb, 1pb, 1eb
	ConstFloat                           // 1.0
	ConstInt                             // 1
	ConstStringSingleQuote               // 'string'
	ConstStringDoubleQuote               // "string"
	ConstStringBacktickQoute             // `string`
	ConstBool                            // true, false
	ConstNull                            // null
	Return                               // =>
	Assign                               // =
	Colon                                // :
	Comma                                // ,
	Extend                               // ...
	OpenCurly                            // {
	CloseCurly                           // }
	OpenParen                            // (
	CloseParen                           // )
	OpenAngle                            // <
	CloseAngle                           // >
	RightComment                         // #
	TopComment                           // #
	CustomError                          // error
)

func (t Type) String() string {
	switch t {
	case Error:
		return "Error"
	case EOF:
		return "EOF"
	case Identifier:
		return "Identifier"
	case Const:
		return "Const"
	case Enum:
		return "Enum"
	case Model:
		return "Model"
	case Http:
		return "Http"
	case Rpc:
		return "Rpc"
	case Service:
		return "Service"
	case Byte:
		return "Byte"
	case Bool:
		return "Bool"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
	case Uint32:
		return "Uint32"
	case Uint64:
		return "Uint64"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Timestamp:
		return "Timestamp"
	case String:
		return "String"
	case Map:
		return "Map"
	case Array:
		return "Array"
	case Any:
		return "Any"
	case File:
		return "File"
	case Stream:
		return "Stream"
	case ConstDuration:
		return "ConstDuration"
	case ConstBytes:
		return "ConstBytes"
	case ConstFloat:
		return "ConstFloat"
	case ConstInt:
		return "ConstInt"
	case ConstStringSingleQuote:
		return "ConstStringSingleQuote"
	case ConstStringDoubleQuote:
		return "ConstStringDoubleQuote"
	case ConstStringBacktickQoute:
		return "ConstStringBacktickQoute"
	case ConstBool:
		return "ConstBool"
	case ConstNull:
		return "ConstNull"
	case Return:
		return "Return"
	case Assign:
		return "Assign"
	case Colon:
		return "Colon"
	case Comma:
		return "Comma"
	case Extend:
		return "Extend"
	case OpenCurly:
		return "OpenCurly"
	case CloseCurly:
		return "CloseCurly"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	case OpenAngle:
		return "OpenAngle"
	case CloseAngle:
		return "CloseAngle"
	case RightComment:
		return "RightComment"
	case TopComment:
		return "TopComment"
	case CustomError:
		return "CustomError"
	default:
		return "Unknown"
	}
}
