package example

import (
	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/test"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleSuccess(t *testing.T) {
	t.Parallel()
	app, _ := test.GetDefaultAppTest(t, nil)

	jsonLoginRequest, err := json.Marshal(dto.ExampleRequest{Data: "test"})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/example", strings.NewReader(string(jsonLoginRequest)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.FiberInstance.Test(req)
	if err != nil {
		t.Error(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	bodyMap := make(map[string]interface{})

	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
