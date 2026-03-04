package crypto

import (
	"strconv"
	"strings"

	"github.com/serkanalgur/phpfuncs"
	"github.com/xyproto/randomstring"
)

const (
	HashVersionMd5    = 0
	HashVersionSha256 = 1
	HashVersionLatest = 1
	DefaultSaltLength = 32
	PasswordHash      = 0
	PasswordSalt      = 1
	PasswordVersion   = 2
	CipherLatest      = 2
	Delimiter         = ":"
)

type MGPassword struct {
	hashVersionMap  map[int]interface{}
	passwordHashMap map[int]interface{}
	cipher          int
}

func NewMGPassword() *MGPassword {
	hashVersionMap := map[int]interface{}{
		HashVersionMd5:    "md5",
		HashVersionSha256: "sha256",
	}

	passwordHashMap := map[int]interface{}{
		PasswordHash:    "",
		PasswordSalt:    "",
		PasswordVersion: HashVersionLatest,
	}

	cipher := CipherLatest

	return &MGPassword{
		hashVersionMap,
		passwordHashMap,
		cipher,
	}
}

func (mg *MGPassword) IsValidHash(password string, hash string) bool {
	mg.explodePasswordHash(hash)

	for _, hashVersion := range mg.getPasswordVersion() {
		password = mg.hash(mg.getPasswordSalt()+password, hashVersion)
	}

	return mg.compareStrings(password, mg.getPasswordHash())
}

func (mg *MGPassword) explodePasswordHash(hash string) map[int]interface{} {
	explodedPassword := strings.Split(hash, Delimiter)
	explodedPasswordMap := make(map[int]interface{})

	for i := 0; i < len(explodedPassword); i++ {
		explodedPasswordMap[i] = explodedPassword[i]
	}

	for key, defaultValue := range mg.passwordHashMap {
		pass, ok := explodedPasswordMap[key]
		if ok {
			mg.passwordHashMap[key] = pass
		} else {
			switch defaultValue := defaultValue.(type) {
			case int:
				value := strconv.Itoa(defaultValue)
				mg.passwordHashMap[key] = value

			default:
				mg.passwordHashMap[key] = defaultValue
			}
		}
	}

	return mg.passwordHashMap
}

func (mg *MGPassword) getPasswordVersion() []int {
	strSlice := strings.Split(mg.passwordHashMap[PasswordVersion].(string), Delimiter)
	intSlice := make([]int, len(strSlice))

	for i := 0; i < len(strSlice); i++ {
		intSlice[i], _ = strconv.Atoi(strSlice[i])
	}

	return intSlice
}

func (mg *MGPassword) hash(data string, version int) string {
	return phpfuncs.Hash(mg.hashVersionMap[version].(string), data)
}

func (mg *MGPassword) getPasswordSalt() string {
	return mg.passwordHashMap[PasswordSalt].(string)
}

func (mg *MGPassword) getPasswordHash() string {
	return mg.passwordHashMap[PasswordHash].(string)
}

func (mg *MGPassword) compareStrings(expected string, actual string) bool {
	return expected == actual
}

func (mg *MGPassword) GetHash(password string, salt bool, version int) string {
	if !salt {
		return mg.hash(password, version)
	}

	saltStr := randomstring.String(DefaultSaltLength)

	return strings.Join([]string{mg.hash(saltStr+password, version), saltStr, strconv.Itoa(version)}, Delimiter)
}
