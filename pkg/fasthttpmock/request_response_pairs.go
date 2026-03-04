package fasthttpmock

import (
	"sync"
	"testing"

	"github.com/valyala/fasthttp"
)

type requestResponsePair struct {
	request  *fasthttp.Request
	response *fasthttp.Response
}

type RequestResponsePairs struct {
	mu    sync.Mutex
	pairs []requestResponsePair
}

func NewRequestResponsePairs() *RequestResponsePairs {
	return &RequestResponsePairs{}
}

func (p *RequestResponsePairs) Add(request *fasthttp.Request, response *fasthttp.Response) {
	// Deep copy to avoid later mutations (Reset/reuse) zeroing stored data
	reqCopy := &fasthttp.Request{}
	request.CopyTo(reqCopy)

	respCopy := &fasthttp.Response{}
	response.CopyTo(respCopy)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.pairs = append(p.pairs, requestResponsePair{
		request:  reqCopy,
		response: respCopy,
	})
}

func (p *RequestResponsePairs) AddWithPairs(pair *RequestResponsePairs) {
	// Deep copy incoming pairs to detach from external mutations
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pr := range pair.pairs {
		reqCopy := &fasthttp.Request{}
		pr.request.CopyTo(reqCopy)

		respCopy := &fasthttp.Response{}
		pr.response.CopyTo(respCopy)
		p.pairs = append(p.pairs, requestResponsePair{request: reqCopy, response: respCopy})
	}
}

func (p *RequestResponsePairs) AssertSend(t *testing.T, callback func(request *fasthttp.Request) bool) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pair := range p.pairs {
		if callback(pair.request) {
			return true
		}
	}

	t.Error("It is impossible to find a request that satisfies the callback conditions")

	return false
}
