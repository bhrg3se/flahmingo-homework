package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

// ParseConfig uses viper to read and parse config file.
func ParseConfig(path string) Config {
	var config Config
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic("config file not found in " + filepath.Join(path))
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(absPath)

	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "flahmingo")

	viper.SetDefault("server.listen", "127.0.0.1:9090")
	viper.SetDefault("database.user", "flahmingo")
	viper.SetDefault("database.user", "flahmingo")
	viper.SetDefault("database.user", "flahmingo")
	viper.SetDefault("database.user", "flahmingo")

	if err = viper.ReadInConfig(); err != nil {
		logrus.Fatalf("could not read config file: %v", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		logrus.Fatalf("config file invalid: %v", err)
	}

	return config
}

type Config struct {
	Database struct {
		User         string `toml:"user"`
		Password     string `toml:"password"`
		Host         string `toml:"host"`
		Port         string `toml:"port"`
		Name         string `toml:"name"`
		SSL          bool   `toml:"ssl"`
		CaCertPath   string `json:"caCertPath"`
		UserCertPath string `json:"userCertPath"`
		UserKeyPath  string `json:"userKeyPath"`
	} `toml:"database"`

	Logging struct {
		Level string `toml:"logging"`
	} `toml:"logging"`

	Server struct {
		Listen string `toml:"listen"`
	} `toml:"server"`

	GoogleCloud struct {
		ProjectID string `toml:"projectID"`
	} `toml:"googleCloud"`

	Twilio struct {
		AccountSID  string `toml:"accountSid"`
		AuthToken   string `toml:"authToken"`
		PhoneNumber string `toml:"phoneNumber"`
	} `toml:"twilio"`
}
