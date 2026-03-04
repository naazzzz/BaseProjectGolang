package trait

import (
	"BaseProjectGolang/pkg/fasthttpmock"

	"github.com/valyala/fasthttp"
)

func MockHTTPClient(pairs *fasthttpmock.RequestResponsePairs) *fasthttpmock.WrapClient {
	fastClient := &fasthttp.Client{}

	client := fasthttpmock.NewWrapClient(fastClient)

	mockClient := fasthttpmock.NewClient(pairs, fasthttpmock.Equal, fasthttpmock.Copy)

	// // switch to mock usage
	client.SetMockClient(mockClient)

	// switch to normal usage
	// client.SetMockClient(nil)

	return client
}

func CreatePairsForMockClient(
	requestURL string,
	requestBody *string,
	requestMethod string,
	responseBody string,
	responseStatus int,
) *fasthttpmock.RequestResponsePairs {
	pairs := fasthttpmock.NewRequestResponsePairs()
	{
		request := &fasthttp.Request{}
		request.Header.SetMethod(requestMethod)
		request.SetRequestURI(requestURL)

		if requestBody != nil {
			request.SetBodyString(*requestBody)
		}

		response := &fasthttp.Response{}
		response.SetStatusCode(responseStatus)
		response.SetBodyString(responseBody)

		pairs.Add(request, response)
	}

	return pairs
}
