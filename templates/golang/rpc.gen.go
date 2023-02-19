package rpc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

func callServiceMethod[ReqMsg, RespMsg any](ctx context.Context, client *http.Client, url string, in *ReqMsg, out *RespMsg) (err error) {
	var reqBody io.Reader

	if !isStructEmpty(in) {
		reqBody, err = jsonEncoder(in)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil
	}

	if isStructEmpty(out) {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

type streamMapper map[string]func(id int64, data string)

func callServiceStreamMethod[ReqMsg any](ctx context.Context, client *http.Client, url string, in *ReqMsg, out streamMapper) (err error) {
	var reqBody io.Reader

	if !isStructEmpty(in) {
		reqBody, err = jsonEncoder(in)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return nil
	}

	go func() {
		defer resp.Body.Close()

		for event := range parseStreamEvents(resp.Body) {
			if handler, ok := out[event.event]; ok {
				handler(event.id, event.data)
			}
		}
	}()

	return
}

func parseStreamData[T any](out chan<- *T) func(id int64, data string) {
	return func(id int64, data string) {
		var item T
		err := json.Unmarshal([]byte(data), &item)
		if err != nil {
			return
		}
		out <- &item
	}
}

type serviceMethodHandler func(context.Context, http.ResponseWriter, *http.Request)

type stremEvent struct {
	id    int64
	event string
	data  string
}

func parseStreamEvents(r io.Reader) <-chan *stremEvent {
	out := make(chan *stremEvent)

	scanner := bufio.NewScanner(r)

	// Set the scanner's split function to split on "\n\n"
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Return nothing if at end of file and no data passed
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		idx := bytes.Index(data, []byte("\n\n"))
		if idx >= 0 {
			return idx + 2, data[:idx], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		// We need more data
		return 0, nil, nil
	})

	secondPart := func(value string) (string, bool) {
		segments := strings.Split(value, ":")
		if len(segments) != 2 {
			return "", false
		}
		return strings.TrimSpace(segments[1]), true
	}

	go func() {
		defer close(out)

		for scanner.Scan() {
			item := scanner.Text()
			lines := strings.Split(item, "\n")

			if len(lines) != 3 {
				continue
			}

			identifier, ok := secondPart(lines[0])
			if !ok {
				continue
			}

			id, err := strconv.ParseInt(identifier, 10, 64)
			if err != nil {
				continue
			}

			event, ok := secondPart(lines[1])
			if !ok {
				continue
			}

			data, ok := secondPart(lines[2])
			if !ok {
				continue
			}

			out <- &stremEvent{
				id:    id,
				event: event,
				data:  data,
			}
		}
	}()

	return out
}

func createEventStream[T any](ctx context.Context, name string, stream <-chan T) <-chan *stremEvent {
	out := make(chan *stremEvent)

	go func() {
		defer close(out)
		var id int64 = 1
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-stream:
				if !ok {
					return
				}
				b, _ := json.Marshal(item)
				out <- &stremEvent{
					id:    id,
					event: name,
					data:  string(b),
				}
				id++
			}
		}
	}()

	return out
}

func createStreamServiceMethod[ReqMsg any](fn func(ctx context.Context, req *ReqMsg) (<-chan *stremEvent, error)) serviceMethodHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ResponseError(w, Errorf(http.StatusMethodNotAllowed, "method %q not allowed", r.Method))
			return
		}

		fluser, ok := w.(http.Flusher)
		if !ok {
			ResponseError(w, Errorf(http.StatusInternalServerError, "response writer does not support flushing"))
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

		events, err := fn(ctx, &reqMsg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		fluser.Flush()

		var lastID int64
		for item := range events {
			fmt.Fprintf(w, "id: %d\nevent: %s\ndata: %s\n\n", item.id, item.event, item.data)
			fluser.Flush()
			lastID = item.id
		}

		lastID++
		fmt.Fprintf(w, "id: %d\nevent: done\ndata: {}\n\n", lastID)
	}
}

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

// Utility functions

func mergeChannels[Resp any](ctx context.Context, chans ...<-chan Resp) <-chan Resp {
	out := make(chan Resp)
	done := make(chan struct{})
	wait := make(chan struct{}, len(chans))

	output := func(c <-chan Resp) {
		defer func() {
			wait <- struct{}{}
		}()

		for n := range c {
			select {
			case <-done:
				return
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	for _, c := range chans {
		go output(c)
	}

	go func() {
		defer close(done)
		defer close(out)

		for i := 0; i < len(chans); i++ {
			select {
			case <-wait:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
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

func urlPathJoin(host string, paths ...string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, path.Join(paths...))
	return u.String(), nil
}

type emptyStruct struct{}

func isStructEmpty(value any) bool {
	_, ok := value.(*emptyStruct)
	return ok
}

func jsonEncoder(value any) (io.Reader, error) {
	r, w := io.Pipe()
	go func() {
		err := json.NewEncoder(w).Encode(value)
		if err != nil {
			w.CloseWithError(err)
			return
		}
		w.Close()
	}()
	return r, nil
}

// Error

func writeResponse(w http.ResponseWriter, code int, body any) {
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

	err := json.NewEncoder(&buffer).Encode(struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  e.code,
		Error: e.Error(),
	})

	return buffer.Bytes(), err
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
		writeResponse(w, err.code, err)
	default:
		writeResponse(w, http.StatusInternalServerError, err)
	}
}
