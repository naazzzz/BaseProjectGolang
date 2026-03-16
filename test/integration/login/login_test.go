package login

import (
	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/test"
	"BaseProjectGolang/test/trait"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	t.Parallel()
	app, cfg := test.GetDefaultAppTest(t, nil)

	user := trait.CreateUserWithServiceInfo(app.Database, cfg, nil)

	jsonLoginRequest, err := json.Marshal(dto.LoginRequest{Username: user.Username, Password: "secret"})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(jsonLoginRequest)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.FiberInstance.Test(req, fiber.TestConfig{Timeout: 10000000})
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
	assert.Contains(t, bodyMap, "tokendmn")
	assert.Contains(t, bodyMap, "token_type")
	assert.Contains(t, bodyMap, "expires_at")
}

func TestLoginErrorWrongPassword(t *testing.T) {
	t.Parallel()
	app, cfg := test.GetDefaultAppTest(t, nil)

	user := trait.CreateUserWithServiceInfo(app.Database, cfg, nil)

	jsonLoginRequest, err := json.Marshal(dto.LoginRequest{Username: user.Username, Password: "wrong_secret"})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(jsonLoginRequest)))
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

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Contains(t, bodyMap["message"], "Wrong username or password")
}

func TestLoginErrorWrongUsername(t *testing.T) {
	t.Parallel()
	app, cfg := test.GetDefaultAppTest(t, nil)

	trait.CreateUserWithServiceInfo(app.Database, cfg, nil)

	jsonLoginRequest, err := json.Marshal(dto.LoginRequest{Username: "wrong_username", Password: "secret1"})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(jsonLoginRequest)))
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

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Contains(t, bodyMap["message"], "Wrong username or password")
	assert.Contains(t, bodyMap["metadata"].(map[string]interface{})["description"], "record not found")
}

func TestLoginValidationErrorRequiredUsernameAndPassword(t *testing.T) {
	t.Parallel()
	app, cfg := test.GetDefaultAppTest(t, nil)

	trait.CreateUserWithServiceInfo(app.Database, cfg, nil)

	jsonLoginRequest, err := json.Marshal(dto.LoginRequest{Username: "", Password: ""})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(jsonLoginRequest)))
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

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Contains(t, bodyMap["message"].([]interface{})[0], "Key: 'LoginRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag")
	assert.Contains(t, bodyMap["message"].([]interface{})[1], "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag")
}

func TestLoginValidationErrorMinMaxUsernameAndPassword(t *testing.T) {
	t.Parallel()
	app, cfg := test.GetDefaultAppTest(t, nil)

	trait.CreateUserWithServiceInfo(app.Database, cfg, nil)

	jsonLoginRequest, err := json.Marshal(
		dto.LoginRequest{
			Username: "pZmmynUdpRHnZbyiLqojkhXHyzXVCZoIZUztiyLJsQIcEMOXmnbOIwPQXLthKaUdWFGcOhIxmipXUSTtFsBKAlJvBDkmUZHUwMMPfNtCAtEBtOmihNYuNzCGcfFCCpmwijRARJMOYkJnuNzLIARRLIdwaflsLGOTwAVwaOiNJjkhRQgSDWgjWmzjIdUQuey2",
			Password: "se",
		})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(jsonLoginRequest)))
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

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Contains(t, bodyMap["message"].([]interface{})[0], "Key: 'LoginRequest.Username' Error:Field validation for 'Username' failed on the 'max' tag")
	assert.Contains(t, bodyMap["message"].([]interface{})[1], "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag")
}
