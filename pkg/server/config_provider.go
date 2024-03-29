package server

import (
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type applicationConfig struct {
	DBConnectionString string `mapstructure:"DB_CONNECTION_STRING" validate:"required"`
}

type ConfigurationProvider struct {
	config *applicationConfig
}

func (cfg *ConfigurationProvider) loadConfig() {

	envFilePath := os.Getenv("ENV_FILE")
	if len(envFilePath) > 0 {
		viper.SetConfigFile(envFilePath)
		viper.SetConfigType("env")
	}

	_ = viper.ReadInConfig()

	var config applicationConfig
	_ = viper.Unmarshal(&config)

	st := reflect.ValueOf(&config).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		if envKey, ok := field.Tag.Lookup("mapstructure"); ok {
			value := os.Getenv(envKey)
			if len(value) > 0 {
				st.Field(i).SetString(value)
			}

		}
	}

	v := validator.New()
	if err := v.Struct(config); err != nil {
		panic(err)
	}

	cfg.config = &config
}

func (cfg *ConfigurationProvider) GetDBConnectionString() *string {
	return &cfg.config.DBConnectionString

}

func NewConfigurationProvider() *ConfigurationProvider {
	cfgProvider := &ConfigurationProvider{}
	cfgProvider.loadConfig()
	return cfgProvider
}
