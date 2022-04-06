package config

import "sync"

const (
	AppEnvProd    = "prod"
	AppEnvStaging = "staging"
	AppEnvTest    = "test"
)

var once sync.Once

func InitialiseConfig() {
	once.Do(func() {
		LoadDBConfig()
	})
}
