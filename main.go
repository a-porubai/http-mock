package main

import (
	pingv1 "awesomeProject/gen/connect/ping/v1"
	"awesomeProject/gen/connect/ping/v1/pingv1connect"
	httpclientmock "awesomeProject/http_client/tests"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"connectrpc.com/connect"
)

func main() {
	url := "http://localhost:8089"

	httpClient := httpclientmock.New()

	mockedBody := `{"number": 42,"text": "test"}`
	httpClient.MockResponse(fmt.Sprintf("%v/connect.ping.v1.PingService/Ping", url), func() *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(mockedBody)),
		}
	})

	client := pingv1connect.NewPingServiceClient(
		httpClient,
		url,
		connect.WithGRPC(),
	)
	req := connect.NewRequest(&pingv1.PingRequest{
		Number: 42,
	})
	req.Header().Set("Some-Header", "hello from connect")
	res, err := client.Ping(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	log.Println(res.Header().Get("Some-Other-Header"))
}
