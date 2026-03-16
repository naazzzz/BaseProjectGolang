package create

import (
	domainUser "BaseProjectGolang/internal/domain/userdmn"
	"context"
)

type Handler struct {
	userRepository domainUser.IUserRepository
}

func NewHandler(
	userRepository domainUser.IUserRepository,
) *Handler {
	return &Handler{
		userRepository: userRepository,
	}
}

func (h *Handler) Execute(
	ctx context.Context,
	cmd *Command,
) (*Result, error) {

	user := &domainUser.User{
		Username: cmd.Username,
		Password: cmd.Password,
		Active:   true,
	}

	id, err := h.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &Result{
		ID: id,
	}, nil
}
