package contextkeys

type contextKey string

const (
	AuthorizedUser   contextKey = "authorized_user"
	TransactionKey   contextKey = "transaction"
	RequestCtxKey    contextKey = "request-ctx"
	FactoryKey       contextKey = "factory"
	QbCtxKey         contextKey = "queryBuilder"
	PaginationCtxKey contextKey = "pagination"
)

// константы из кода(querybuilder, header)
const (
	ServiceConfigsCtxKey        string = "service_configs"
	UsernameParamsKey           string = "username"
	PasswordParamsKey           string = "password"
	SignatureParamsKey          string = "signature"
	SignatureTimestampParamsKey string = "signature_timestamp"
)
