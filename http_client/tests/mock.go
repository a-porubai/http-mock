package httpclientmock

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

type ResponseFactory func() *http.Response

type Client struct {
	requests  map[string][]*http.Request
	responses map[string][]ResponseFactory
	mu        sync.Mutex
}

func New() *Client {
	return &Client{
		requests:  make(map[string][]*http.Request),
		responses: make(map[string][]ResponseFactory),
	}
}

func (c *Client) Reset() {
	c.requests = make(map[string][]*http.Request)
	c.responses = make(map[string][]ResponseFactory)
}

func (c *Client) MockResponse(reqURL string, factory ResponseFactory) {
	c.responses[reqURL] = []ResponseFactory{factory}
}

func (c *Client) MockNextResponse(reqURL string, factory ResponseFactory) {
	c.responses[reqURL] = append(c.responses[reqURL], factory)
}

func (c *Client) GetRequests(reqURL string) ([]*http.Request, error) {
	req, ok := c.requests[reqURL]

	if ok {
		return req, nil
	}

	return nil, errors.Errorf("%q URL was not requested", reqURL)
}

func (c *Client) RemoveProcessedRequest(reqURL string, reqIndex int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.requests[reqURL] == nil {
		return errors.New("processed request not found")
	}

	if len(c.requests[reqURL]) <= reqIndex {
		return errors.New("processed request index is invalid")
	}

	c.requests[reqURL] = append(c.requests[reqURL][:reqIndex], c.requests[reqURL][reqIndex+1:]...)

	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	reqURL := req.URL.String()

	//c.mu.Lock()
	//c.requests[reqURL] = append(c.requests[reqURL], req)
	//c.mu.Unlock()

	respFactories, ok := c.responses[reqURL]
	if !ok {
		return nil, errors.Errorf("response for %q URL not mocked", reqURL)
	}

	respFactory := respFactories[0]

	if len(respFactories) > 1 {
		c.responses[reqURL] = respFactories[1:]
	}

	fmt.Println(respFactory())

	return respFactory(), nil
}
