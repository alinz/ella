package ast

import (
	"strconv"
	"strings"

	"compiler.ella.to/internal/token"
)

type CustomError struct {
	Token      *token.Token
	Name       *Identifier
	Code       int64
	HttpStatus int
	Msg        *ValueString
}

var _ Statement = (*Enum)(nil)

func (c *CustomError) statementLiteral() {}

func (c *CustomError) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CustomError) String() string {
	var sb strings.Builder

	sb.WriteString("error ")
	sb.WriteString(c.Name.String())
	sb.WriteString(" {")

	sb.WriteString(" Code = ")
	sb.WriteString(strconv.FormatInt(c.Code, 10))

	sb.WriteString(" HttpStatus = ")
	sb.WriteString(HttpStatusCode2String[c.HttpStatus])

	sb.WriteString(" Msg = ")
	sb.WriteString(c.Msg.String())

	sb.WriteString(" }")

	return sb.String()
}

var HttpStatusCode2String = map[int]string{
	100: "Continue",
	101: "SwitchingProtocols",
	102: "Processing",
	103: "EarlyHints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "NonAuthoritativeInfo",
	204: "NoContent",
	205: "ResetContent",
	206: "PartialContent",
	207: "MultiStatus",
	208: "AlreadyReported",
	226: "IMUsed",
	300: "MultipleChoices",
	301: "MovedPermanently",
	302: "Found",
	303: "SeeOther",
	304: "NotModified",
	305: "UseProxy",
	307: "TemporaryRedirect",
	308: "PermanentRedirect",
	400: "BadRequest",
	401: "Unauthorized",
	402: "PaymentRequired",
	403: "Forbidden",
	404: "NotFound",
	405: "MethodNotAllowed",
	406: "NotAcceptable",
	407: "ProxyAuthRequired",
	408: "RequestTimeout",
	409: "Conflict",
	410: "Gone",
	411: "LengthRequired",
	412: "PreconditionFailed",
	413: "RequestEntityTooLarge",
	414: "RequestURITooLong",
	415: "UnsupportedMediaType",
	416: "RequestedRangeNotSatisfiable",
	417: "ExpectationFailed",
	418: "Teapot",
	421: "MisdirectedRequest",
	422: "UnprocessableEntity",
	423: "Locked",
	424: "FailedDependency",
	425: "TooEarly",
	426: "UpgradeRequired",
	428: "PreconditionRequired",
	429: "TooManyRequests",
	431: "RequestHeaderFieldsTooLarge",
	451: "UnavailableForLegalReasons",
	500: "InternalServerError",
	501: "NotImplemented",
	502: "BadGateway",
	503: "ServiceUnavailable",
	504: "GatewayTimeout",
	505: "HTTPVersionNotSupported",
	506: "VariantAlsoNegotiates",
	507: "InsufficientStorage",
	508: "LoopDetected",
	510: "NotExtended",
	511: "NetworkAuthenticationRequired",
}

var HttpStatusString2Code = map[string]int{
	"Continue":                      100,
	"SwitchingProtocols":            101,
	"Processing":                    102,
	"EarlyHints":                    103,
	"OK":                            200,
	"Created":                       201,
	"Accepted":                      202,
	"NonAuthoritativeInfo":          203,
	"NoContent":                     204,
	"ResetContent":                  205,
	"PartialContent":                206,
	"MultiStatus":                   207,
	"AlreadyReported":               208,
	"IMUsed":                        226,
	"MultipleChoices":               300,
	"MovedPermanently":              301,
	"Found":                         302,
	"SeeOther":                      303,
	"NotModified":                   304,
	"UseProxy":                      305,
	"TemporaryRedirect":             307,
	"PermanentRedirect":             308,
	"BadRequest":                    400,
	"Unauthorized":                  401,
	"PaymentRequired":               402,
	"Forbidden":                     403,
	"NotFound":                      404,
	"MethodNotAllowed":              405,
	"NotAcceptable":                 406,
	"ProxyAuthRequired":             407,
	"RequestTimeout":                408,
	"Conflict":                      409,
	"Gone":                          410,
	"LengthRequired":                411,
	"PreconditionFailed":            412,
	"RequestEntityTooLarge":         413,
	"RequestURITooLong":             414,
	"UnsupportedMediaType":          415,
	"RequestedRangeNotSatisfiable":  416,
	"ExpectationFailed":             417,
	"Teapot":                        418,
	"MisdirectedRequest":            421,
	"UnprocessableEntity":           422,
	"Locked":                        423,
	"FailedDependency":              424,
	"TooEarly":                      425,
	"UpgradeRequired":               426,
	"PreconditionRequired":          428,
	"TooManyRequests":               429,
	"RequestHeaderFieldsTooLarge":   431,
	"UnavailableForLegalReasons":    451,
	"InternalServerError":           500,
	"NotImplemented":                501,
	"BadGateway":                    502,
	"ServiceUnavailable":            503,
	"GatewayTimeout":                504,
	"HTTPVersionNotSupported":       505,
	"VariantAlsoNegotiates":         506,
	"InsufficientStorage":           507,
	"LoopDetected":                  508,
	"NotExtended":                   510,
	"NetworkAuthenticationRequired": 511,
}
