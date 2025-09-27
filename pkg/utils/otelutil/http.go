package otelutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var wrappedTransport = otelhttp.NewTransport(http.DefaultTransport)
var client = http.Client{Transport: wrappedTransport}

type HttpRequest struct {
	ctx      context.Context
	header   http.Header
	query    url.Values
	formData url.Values
	body     any
}

func R(ctx context.Context) *HttpRequest {
	return &HttpRequest{
		ctx:      ctx,
		header:   http.Header{},
		query:    url.Values{},
		formData: url.Values{},
	}
}

func (request *HttpRequest) request(method, url string) (*http.Request, error) {
	var body io.Reader
	if request.body != nil {
		s, err := jsoniter.MarshalToString(request.body)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(s)
	} else if len(request.formData) > 0 {
		body = strings.NewReader(request.formData.Encode())
	}
	if len(request.query) > 0 {
		divder := "?"
		if len(strings.Split(url, "?")) > 1 {
			divder = "&"
		}
		url = fmt.Sprintf("%s%s%s", url, divder, request.query.Encode())
	}
	req, err := http.NewRequestWithContext(
		request.ctx,
		method,
		url,
		body,
	)
	req.Header = request.header
	return req, err
}

func (request *HttpRequest) SetHeader(key, value string) *HttpRequest {
	request.header.Set(key, value)
	return request
}

func (request *HttpRequest) SetQuery(key, value string) *HttpRequest {
	request.query.Set(key, value)
	return request
}

func (request *HttpRequest) SetForm(key, value string) *HttpRequest {
	request.formData.Set(key, value)
	return request
}

func (request *HttpRequest) SetFormData(data map[string]string) *HttpRequest {
	for k, v := range data {
		request.formData.Set(k, v)
	}
	return request
}

func (request *HttpRequest) SetJSON(data any) *HttpRequest {
	request.SetHeader("Content-Type", "application/json")
	request.body = data
	return request
}

func (request *HttpRequest) SetAuthToken(token string) *HttpRequest {
	request.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return request
}

func (request *HttpRequest) Execute(method string, url string) (*HttpResponse, error) {
	req, err := request.request(method, url)
	if err != nil {
		return nil, err
	}
	return HttpDo(req)
}

func (request *HttpRequest) Get(url string) (*HttpResponse, error) {
	return request.Execute(http.MethodGet, url)
}

func (request *HttpRequest) Post(url string) (*HttpResponse, error) {
	return request.Execute(http.MethodPost, url)
}

func (request *HttpRequest) Put(url string) (*HttpResponse, error) {
	return request.Execute(http.MethodPut, url)
}

type HttpResponse struct {
	statusCode int
	status     string
	header     http.Header
	body       []byte
}

func (response *HttpResponse) StatusCode() int {
	return response.statusCode
}

func (response *HttpResponse) Status() string {
	return response.status
}

func (response *HttpResponse) Header() http.Header {
	return response.header
}

func (response *HttpResponse) Body() []byte {
	return response.body
}

func HttpDo(req *http.Request) (*HttpResponse, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &HttpResponse{
		statusCode: resp.StatusCode,
		status:     resp.Status,
		header:     resp.Header,
		body:       body,
	}, nil
}
