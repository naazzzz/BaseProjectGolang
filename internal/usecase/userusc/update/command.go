package update

import domainUser "BaseProjectGolang/internal/domain/userdmn"

type Command struct {
	Username string
	NewUser  *domainUser.User
}
