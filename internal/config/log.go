package config

type (
	Logs struct {
		LifeTime   string `envconfig:"LOGS_LIFETIME" default:"2"`
		DirSave    string `envconfig:"LOGS_DIR_SAVE" default:"log"`
		FileName   string `envconfig:"LOGS_FILE_NAME" default:"logs.log"`
		MaxBackups string `envconfig:"LOGS_MAX_BACKUPS" default:"4"`
		MaxSize    string `envconfig:"LOGS_MAX_SIZE" default:"5"`
	}
)

func (logs *Logs) GetMaxSize() string {
	return logs.MaxSize
}

func (logs *Logs) GetMaxBackups() string {
	return logs.MaxBackups
}

func (logs *Logs) GetLifeTime() string {
	return logs.LifeTime
}

func (logs *Logs) GetDirSave() string {
	return logs.DirSave
}

func (logs *Logs) GetFilename() string {
	return logs.FileName
}
