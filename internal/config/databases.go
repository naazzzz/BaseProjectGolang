package config

type (
	Databases struct {
		ConnectionMaxLifeTime string `envconfig:"CONNECTION_MAX_LIFE_TIME" default:"300"`
		MaxOpenConnections    string `envconfig:"MAX_OPEN_CONNECTIONS" default:"10"`
		MaxIdleConnections    string `envconfig:"MAX_IDLE_CONNECTIONS" default:"10"`
		Pgsql                 *Pgsql
	}

	Pgsql struct {
		Database string `envconfig:"POSTGRES_DB" default:"monitor"`
		User     string `envconfig:"POSTGRES_USER" default:"admin"`
		Password string `envconfig:"POSTGRES_PASSWORD" default:"master123"`
		Host     string `envconfig:"POSTGRES_HOST" default:"db_pgsql"`
		Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
	}
)
