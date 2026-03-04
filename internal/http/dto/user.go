package dto

import (
	userModel "BaseProjectGolang/internal/infrastructure/database/orm/model/user"
	"BaseProjectGolang/pkg/converter"
)

type UserRequest struct {
	Username string `validate:"required,max=191" json:"username"`
	Password string `validate:"required,max=128" json:"password"`
	Active   bool   `validate:"required" json:"active"`
}

func (userRequest *UserRequest) ToUserModel() (user *userModel.User, err error) {
	user, err = converter.TypeConverter[*userModel.User](userRequest)
	user.Password = userRequest.Password

	return
}
