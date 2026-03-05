package update

import domainUser "BaseProjectGolang/internal/domain/user"

type Command struct {
	Username string
	NewUser  *domainUser.User
}
