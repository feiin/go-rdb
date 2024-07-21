package setting

import (
	"os"
)

type ENV string

const (
	DEV  ENV = "development"
	PROD ENV = "production"
)

func (e ENV) IsProduction() bool {
	return e == PROD
}

func (e ENV) IsDev() bool {
	return e == DEV
}

func GetEnv() ENV {
	appenv := os.Getenv("APP_ENV")
	if appenv != "" {
		allEnv := map[string]struct{}{
			"development": {},
			"production":  {},
		}

		if _, ok := allEnv[appenv]; ok {
			return ENV(appenv)
		}
	}

	return DEV
}
