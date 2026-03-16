package delete

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
) error {

	return h.userRepository.Delete(
		ctx,
		cmd.Username,
		cmd.Force,
	)
}
