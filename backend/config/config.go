package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Server     server
	Database   Database
	GoogleOidc GoogleOidc
}

type server struct { // TODO: private type
	Host      string
	Port      string `env:"SERVER_PORT,required"`
	JwtSecret string `env:"JWT_SECRET"`
}

type Database struct { // TODO: private type
	Host     string `env:"MONGODB_URI,required"` // TODO: rename to URI
	Name     string `env:"MONGODB_NAME,required"`
	Username string `env:"MONGODB_USERNAME,required"`
	Password string `env:"MONGODB_PASSWORD,required"`
}

type GoogleOidc struct { // TODO: private type
	ClientId     string `env:"GOOGLE_OIDC_CLIENT_ID,required"`
	ClientSecret string `env:"GOOGLE_OIDC_CLIENT_SECRET,required"`
	RedirectUri  string `env:"GOOGLE_OIDC_REDIRECT_URI,required"`
	IsDevMode    bool   `env:"GOOGLE_OIDC_IS_DEV_MODE"`
}

func Env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("Environment variable " + key + " is not set")
	}

	return value
}

func ToBoolean(strVal string) bool {
	boolVal, _ := strconv.ParseBool(strVal)
	return boolVal
}

var once sync.Once
var config Config

func C(envPrefix ...string) Config {
	if len(envPrefix) > 1 {
		log.Fatal("can pass only one prefix for env but your prefix:", envPrefix)
	}

	once.Do(func() {
		var prefix string
		if len(envPrefix) == 1 {
			prefix = fmt.Sprintf("%s_", envPrefix[0])
		}

		opts := env.Options{
			Prefix: prefix,
		}

		// dbconf := &Database{}
		// if err := env.ParseWithOptions(dbconf, opts); err != nil {
		// 	log.Fatal(err)
		// }

		// googleOidc := &GoogleOidc{}
		// if err := env.ParseWithOptions(googleOidc, opts); err != nil {
		// 	log.Fatal(err)
		// }

		srvConf := &server{}
		if err := env.ParseWithOptions(srvConf, opts); err != nil {
			log.Println("env.ParseWithOptions(srvConf, opts) error", err)
		}

		h, _ := os.Hostname()
		port := srvConf.Port
		if port == "" {
			port = os.Getenv("PORT")
		}

		config = Config{
			Server: server{
				Host:      h,
				Port:      port,
				JwtSecret: srvConf.JwtSecret,
			},
			// Database: Database{
			// 	Host:     dbconf.Host,
			// 	Name:     dbconf.Name,
			// 	Username: dbconf.Username,
			// 	Password: dbconf.Password,
			// },
			// GoogleOidc: GoogleOidc{
			// 	ClientId:     googleOidc.ClientId,
			// 	ClientSecret: googleOidc.ClientSecret,
			// 	RedirectUri:  googleOidc.RedirectUri,
			// 	IsDevMode:    googleOidc.IsDevMode,
			// },
		}
	})

	return config
}
