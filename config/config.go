package config

import "github.com/spf13/viper"

// Config struct
type Config struct {
	Port       string `mapstructure:"PORT"`
	DbUsername string `mapstructure:"DB_USERNAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbURL      string `mapstructure:"DB_URL"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbDatabase string `mapstructure:"DB_DATABASE"`
	DbDebug    bool   `mapstructure:"DB_DEBUG"`
}

// ENV const
var ENV *Config

// LoadConfig load config based on env file
func LoadConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic(err)
	}
}

// GetString get config string
func GetString(key string, def ...string) string {
	value := viper.GetString(key)

	if value == "" && len(def) > 0 {
		return def[0]
	}

	return value
}

// GetInt get config int
func GetInt(key string, def ...int) int {
	value := viper.GetInt(key)

	if value == 0 && len(def) > 0 {
		return def[0]
	}

	return value
}

// GetInt64 get config int64
func GetInt64(key string, def ...int64) int64 {
	value := viper.GetInt64(key)

	if value == 0 && len(def) > 0 {
		return def[0]
	}

	return value
}

// GetFloat64 get config float64
func GetFloat64(key string, def ...float64) float64 {
	value := viper.GetFloat64(key)

	if value == 0 && len(def) > 0 {
		return def[0]
	}

	return value
}

// GetBool get config bool
func GetBool(key string, def ...bool) bool {
	value := viper.GetBool(key)

	if len(def) > 0 {
		return def[0]
	}

	return value
}
