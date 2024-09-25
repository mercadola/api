package config

import (
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type Configuration struct {
	TokenAuth              *jwtauth.JWTAuth
	URI                    string `mapstructure:"MONGODB_URI"`
	CustomerCollection     string `mapstructure:"CUSTOMER_COLLECTION"`
	ProductCollection      string `mapstructure:"PRODUCT_COLLECTION"`
	ShoppingListCollection string `mapstructure:"SHOPPING_LIST_COLLECTION"`
	DB                     string `mapstructure:"DATABASE"`
	Port                   string `mapstructure:"PORT"`
	AppName                string `mapstructure:"APP_NAME"`
	ENV                    string `mapstructure:"ENV"`
	JWTSecret              string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn           int    `mapstructure:"JWT_EXPIRES_IN"`
}

var cfg *Configuration

func LoadConfiguration() error {
	cfg = new(Configuration)
	rootPath := os.Getenv("ROOT_PATH")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath(rootPath)

	err := viper.ReadInConfig()
	if err != nil {
		logger := slog.Default()
		logger.Warn("Arquivo .env não encontrado. Usando variáveis de ambiente.")
		viper.AutomaticEnv()
		c := Configuration{}
		BindEnvs(c)
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

func BindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
