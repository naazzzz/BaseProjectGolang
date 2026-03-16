package examplectr

import (
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/internal/usecase/exampleusc"
	"BaseProjectGolang/internal/validation"

	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/soner3/flora"
)

type ExampleController struct {
	flora.Component
	*controller.BaseController
	validator validation.IValidator
	handler   *exampleusc.ExampleHandler
}

func NewExampleController(
	base *controller.BaseController,
	handler *exampleusc.ExampleHandler,
	validator validation.IValidator,
) *ExampleController {
	return &ExampleController{
		BaseController: base,
		validator:      validator,
		handler:        handler,
	}
}

// Example godoc
// @Summary Example endpoint
// @tags Example
// @description Example endpoint
// @Accept json
// @Produce json
// @Param dto.ExampleRequest body dto.ExampleRequest true "Dto для логина"
// @Success      200  {object}  dto.ExampleResponse
// @Failure      default  {object}  error.HTTPError "The body of any response with an error"
// @Router /api/example [post]
func (exampleController *ExampleController) Example(ctx fiber.Ctx) error {
	var exampleRequest *dto.ExampleRequest

	if err := ctx.Bind().Body(&exampleRequest); err != nil {
		return eris.New(err.Error())
	}

	if err := exampleController.validator.Validate(exampleRequest); err != nil {
		return err
	}

	cmd := &exampleusc.ExampleCommand{}

	// Copy the request data to the command object by json tags
	if err := copier.Copy(&cmd, &exampleRequest); err != nil {
		return err
	}

	exampleResult, err := exampleController.handler.Execute(ctx, cmd)
	if err != nil {
		return err
	}

	response := &dto.ExampleResponse{}

	// Copy the response data to the response object by json tags
	if err = copier.Copy(&response, &exampleResult); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
