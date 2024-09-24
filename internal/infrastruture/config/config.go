package config

import (
	"os"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type Database struct {
	URI                    string `mapstructure:"MONGODB_URI"`
	CustomerCollection     string `mapstructure:"CUSTOMER_COLLECTION"`
	ProductCollection      string `mapstructure:"PRODUCT_COLLECTION"`
	ShoppingListCollection string `mapstructure:"SHOPPING_LIST_COLLECTION"`
	DB                     string `mapstructure:"DATABASE"`
}

type Configuration struct {
	TokenAuth    *jwtauth.JWTAuth
	Database     `mapstructure:",squash"`
	Port         string `mapstructure:"PORT"`
	AppName      string `mapstructure:"APP_NAME"`
	ENV          string `mapstructure:"ENV"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn int    `mapstructure:"JWT_EXPIRES_IN"`
	CosmoToken   string `mapstructure:"COSMO_TOKEN"`
	CosmoUrl     string `mapstructure:"COSMO_URL"`
}

var cfg *Configuration

func LoadConfiguration() error {
	cfg = new(Configuration)
	rootPath := os.Getenv("ROOT_PATH")

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath(rootPath)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return err
	}

	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return nil
}

func GetConfig() (*Configuration, error) {
	if cfg == nil {
		err := LoadConfiguration()
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
