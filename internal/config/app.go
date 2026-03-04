package config

type App struct {
	Host       string `envconfig:"APP_HOST"  default:"localhost"`
	URL        string `envconfig:"APP_URL"`
	MainDomain string
}
