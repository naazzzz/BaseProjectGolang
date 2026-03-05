package update

import (
	domainUser "BaseProjectGolang/internal/domain/user"
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
) error {

	user, err := h.userRepository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return err
	}

	user.Username = cmd.NewUser.Username
	user.Password = cmd.NewUser.Password

	return h.userRepository.Update(ctx, user)
}
