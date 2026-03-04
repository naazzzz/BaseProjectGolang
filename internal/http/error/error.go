package error

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"BaseProjectGolang/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// InternalError defines the abstraction for internal error handling.
type InternalError interface {
	GetError() *HTTPError
}

// HTTPError represents the structure of an HTTP error response.
type HTTPError struct {
	Title      string                 `json:"title,omitempty"`
	StatusCode int                    `json:"statusCode"`
	Message    interface{}            `json:"message,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("title: %s, statusCode: %d, message: %v, data: %v, metadata: %v", e.Title, e.StatusCode, e.Message, e.Data, e.Metadata)
}

func (e *HTTPError) GetError() *HTTPError {
	return e
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(title string, statusCode int, message interface{}, data, metadata map[string]interface{}) *HTTPError {
	return &HTTPError{
		Title:      title,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Metadata:   metadata,
	}
}

// NewErrorHandler creates a new Fiber error handler.
func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx fiber.Ctx, receivedError error) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic occurred: %v, stacktrace: %v", r, string(debug.Stack()))
			}
		}()

		httpErr := getErrorDetails(ctx, receivedError)

		data, err := json.Marshal(httpErr)
		if err != nil {
			return err
		}

		log.Println(string(data))

		return ctx.Status(httpErr.StatusCode).JSON(httpErr)
	}
}

// getErrorDetails extracts error details and returns an InternalError.
func getErrorDetails(ctx fiber.Ctx, receivedError error) *HTTPError {
	if receivedError == nil {
		return &HTTPError{}
	}

	externalError := eris.Unpack(receivedError).ErrExternal
	rootError := eris.Unpack(receivedError).ErrRoot
	httpErr := &HTTPError{StatusCode: fiber.StatusInternalServerError}

	switch {
	case errors.Is(externalError, gorm.ErrRecordNotFound):
		httpErr.StatusCode = http.StatusNotFound
		httpErr.Message = nilIfEmpty(rootError.Msg)

	case errors.As(externalError, new(*fiber.Error)):
		var fiberErr *fiber.Error

		errors.As(externalError, &fiberErr)

		httpErr.StatusCode = fiberErr.Code
		httpErr.Message = nilIfEmpty(rootError.Msg)

	case errors.As(externalError, new(validator.ValidationErrors)):
		var validationErrs validator.ValidationErrors

		errors.As(externalError, &validationErrs)

		var validationMap []string

		for _, item := range validationErrs {
			validationMap = append(validationMap, item.Error())
		}

		httpErr.StatusCode = fiber.StatusUnprocessableEntity
		httpErr.Message = validationMap

	default:
		if _, ok := externalError.(InternalError); ok {
			httpErr = externalError.(InternalError).GetError()
		}

		if receivedError.Error() == "missing or malformed JWT" {
			httpErr.StatusCode = http.StatusUnauthorized
		}

		httpErr.Metadata = map[string]interface{}{"description": receivedError.Error()}
	}

	updateMetadata(ctx, httpErr, receivedError)

	httpErr.Title = http.StatusText(httpErr.StatusCode)

	return httpErr
}

// updateMetadata updates the error metadata based on the environment.
func updateMetadata(ctx fiber.Ctx, httpErr *HTTPError, receivedError error) {
	if ctx.App().Config().AppName == config.ProdEnv {
		return
	}

	metadata := map[string]interface{}{
		"trace": eris.Unpack(receivedError).ErrRoot.Stack,
	}

	if externalErr := eris.Unpack(receivedError).ErrExternal; externalErr != nil && externalErr.Error() != "" {
		metadata["description"] = externalErr.Error()
	}

	if httpErr.Metadata != nil {
		httpErr.Metadata["trace"] = eris.Unpack(receivedError).ErrRoot.Stack
	} else {
		httpErr.Metadata = metadata
	}
}

// nilIfEmpty returns nil if the message is empty, otherwise a pointer to the message.
func nilIfEmpty(msg string) interface{} {
	if msg == "" {
		return nil
	}

	return &msg
}
