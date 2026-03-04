package config

type (
	Secure struct {
		AuthPrivateKey string `envconfig:"AUTH_PRIVATE_KEY" default:"secret"`
	}
)
