package rpc

import (
	"context"
	"fmt"
	"net/http"
	"path"
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
	StatusStream(ctx context.Context, userID string) (<-chan *Message, error)
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

func (s *greetingServiceServer) createStatusStreamMethodHandler() serviceMethodHandler {
	return createStreamServiceMethod(func(ctx context.Context, args *struct {
		UserID string `json:"user_id"`
	}) (<-chan *stremEvent, error) {
		stream1, err := s.service.StatusStream(ctx, args.UserID)
		if err != nil {
			return nil, err
		}

		out := mergeChannels(
			ctx,
			createEventStream(ctx, "stream", stream1),
		)

		return out, nil
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
		"/rpc/GreetingService/Echo":         server.createEchoMethodHandler(),
		"/rpc/GreetingService/Ping":         server.createPingMethodHandler(),
		"/rpc/GreetingService/StatusStream": server.createStatusStreamMethodHandler(),
	}

	return &server
}

type greetingServiceClient struct {
	client *http.Client
	host   string
}

var _ GreetingService = (*greetingServiceClient)(nil)

func (s *greetingServiceClient) Echo(ctx context.Context) (_ *Message, err error) {
	url, err := urlPathJoin(s.host, "/rpc/GreetingService/Echo")
	if err != nil {
		return
	}

	in := emptyStruct{}

	out := struct {
		Message *Message `json:"message"`
	}{}

	err = callServiceMethod(ctx, s.client, url, &in, &out)
	if err != nil {
		return nil, err
	}

	return out.Message, nil
}

func (s *greetingServiceClient) Ping(ctx context.Context, userID string) error {
	url := path.Join(s.host, "/rpc/GreetingService/Ping")

	in := struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	}

	out := emptyStruct{}

	err := callServiceMethod(ctx, s.client, url, &in, &out)
	if err != nil {
		return err
	}

	return nil
}

func (s *greetingServiceClient) StatusStream(ctx context.Context, userID string) (<-chan *Message, error) {
	url := path.Join(s.host, "/rpc/GreetingService/StatusStream")

	in := struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	}

	out1 := make(chan *Message)

	streamMapper := streamMapper{
		"stream": parseStreamData(out1),
	}

	err := callServiceStreamMethod(ctx, s.client, url, &in, streamMapper)
	if err != nil {
		return nil, err
	}

	return out1, nil
}

func CreateGreetingServiceClient(host string, client *http.Client) GreetingService {
	return &greetingServiceClient{
		host:   host,
		client: client,
	}
}
