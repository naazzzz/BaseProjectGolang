package fasthttpmock

import (
	"testing"

	"github.com/valyala/fasthttp"
)

type (
	IHTTPClient interface {
		Do(request *fasthttp.Request, response *fasthttp.Response) error
	}

	WrapClient struct {
		realClient   *fasthttp.Client
		mockClient   *Client
		mocked       bool
		SuccessSends *RequestResponsePairs
	}
)

func NewWrapClient(realClient *fasthttp.Client) *WrapClient {
	return &WrapClient{
		realClient:   realClient,
		SuccessSends: &RequestResponsePairs{},
	}
}

func (c *WrapClient) Do(request *fasthttp.Request, response *fasthttp.Response) error {
	var err error

	if c.mocked {
		err = c.mockClient.Do(request, response)
	} else {
		err = c.realClient.Do(request, response)
	}

	if err != nil {
		return err
	}

	c.SuccessSends.Add(request, response)

	return nil
}

func (c *WrapClient) AddPairsToMockHTTPClient(pairs *RequestResponsePairs) {
	c.mockClient.pairs.AddWithPairs(pairs)
}

func (c *WrapClient) SetMockClient(mockClient *Client) {
	c.mockClient = mockClient
	c.mocked = mockClient != nil
}

func (c *WrapClient) AssertSend(t *testing.T, callback func(request *fasthttp.Request) bool) bool {
	return c.SuccessSends.AssertSend(t, callback)
}
