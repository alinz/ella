package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// constants
const (
	Version = "1.0.0"
)

// enums
type Enum int

const (
	Enum1 Enum = 1
)

func (e *Enum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "enum1":
		*e = Enum1
	default:
		return fmt.Errorf("invalid enum value: %s", string(text))
	}
	return nil
}

func (e Enum) MarshalText() ([]byte, error) {
	var name string
	switch e {
	case Enum1:
		name = "Enum1"
	default:
		return nil, fmt.Errorf("invalid enum value: %d", e)
	}
	return []byte(name), nil
}

// messages

type Message struct {
	ID   int64 `json:"id"`
	Enum Enum  `json:"enum"`
}

// services

type GreetingService interface {
	Echo(ctx context.Context) (*Message, error)
	Ping(ctx context.Context, userID string) error
}

type greetingServiceServer struct {
	service GreetingService
	routes  map[string]serviceMethodHandler
}

var _ http.Handler = (*greetingServiceServer)(nil)

func (s *greetingServiceServer) createEchoMethodHandler() serviceMethodHandler {
	return createServiceMethodHandler(func(ctx context.Context, args *struct {
	}) (ret *struct {
		Message *Message `json:"message"`
	}, err error) {
		ret.Message, err = s.service.Echo(ctx)
		return
	})
}

func (s *greetingServiceServer) createPingMethodHandler() serviceMethodHandler {
	return createServiceMethodHandler(func(ctx context.Context, args *struct {
		UserID string `json:"user_id"`
	}) (ret *struct {
	}, err error) {
		err = s.service.Ping(ctx, args.UserID)
		return
	})
}

func (s *greetingServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := s.routes[r.URL.Path]
	if !ok {
		ResponseError(w, Errorf(http.StatusNotFound, "method %q not found", r.URL.Path))
		return
	}

	handler(r.Context(), w, r)
}

func CreateGreetingServiceServer(service GreetingService) http.Handler {
	server := greetingServiceServer{
		service: service,
	}

	server.routes = map[string]serviceMethodHandler{
		"/rpc/GreetingService/Echo": server.createEchoMethodHandler(),
	}

	return &server
}

// http handler

type serviceMethodHandler func(context.Context, http.ResponseWriter, *http.Request)

func createServiceMethodHandler[ReqMsg, RespMsg any](fn func(ctx context.Context, req *ReqMsg) (*RespMsg, error)) serviceMethodHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ResponseError(w, Errorf(http.StatusMethodNotAllowed, "method %q not allowed", r.Method))
			return
		}

		defer func() {
			if err := recover(); err != nil {
				// In case of a panic, serve a 500 error and then panic.
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}
		}()

		if err := checkContentType(r, "application/json"); err != nil {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		var reqMsg ReqMsg

		if err := json.NewDecoder(r.Body).Decode(&reqMsg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		respMsg, err := fn(ctx, &reqMsg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var buffer bytes.Buffer

		if err := json.NewEncoder(&buffer).Encode(respMsg); err != nil {
			ResponseError(w, err)
			return
		}

		io.Copy(w, &buffer)
		w.WriteHeader(http.StatusOK)
	}
}

func checkContentType(r *http.Request, value string) error {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case value:
		return nil
	default:
		return Errorf(http.StatusBadRequest, "unexpected Content-Type: %q", r.Header.Get("Content-Type"))
	}
}

// error

func response(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	switch body := body.(type) {
	case nil:
		return
	case io.Reader:
		io.Copy(w, body)
	default:
		json.NewEncoder(w).Encode(body)
	}
}

type Error struct {
	code  int
	cause error
	msg   string
}

func (e Error) MarshalText() ([]byte, error) {
	var buffer bytes.Buffer

	json.NewEncoder(&buffer).Encode(struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  e.code,
		Error: e.Error(),
	})

	return buffer.Bytes(), nil
}

func (e *Error) Error() string {
	var sb strings.Builder

	sb.WriteString(e.msg)
	if e.msg != "" && e.cause != nil {
		sb.WriteString(": ")
		sb.WriteString(e.cause.Error())
	} else if e.cause != nil {
		sb.WriteString(e.cause.Error())
	}

	return sb.String()
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Unwrap() error {
	return e.cause
}

func Errorf(code int, format string, args ...interface{}) error {
	return &Error{code: code, msg: fmt.Sprintf(format, args...)}
}

func WrapErr(code int, cause error, msg string) error {
	return &Error{code: code, cause: cause, msg: msg}
}

func ResponseError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *Error:
		response(w, err.code, err)
	default:
		response(w, http.StatusInternalServerError, err)
	}
}
