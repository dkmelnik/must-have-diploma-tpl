package configs

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
)

type Server struct {
	ServerAddr  string `envconfig:"RUN_ADDRESS"`
	PGUri       string `envconfig:"DATABASE_URI"`
	AccrualAddr string `envconfig:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
	JWTSecret   string `envconfig:"JWT_SECRET" default:"some_secret"`
	JWTExp      int    `envconfig:"JWT_EXP" default:"1"`
}

func NewServer() (Server, error) {
	cb := Server{}

	flag.StringVar(&cb.ServerAddr, "a", "8080", "in the form 'port'. If empty, 8080 is used")
	flag.StringVar(&cb.PGUri, "d", "", "string for db connect")
	flag.StringVar(&cb.AccrualAddr, "r", "", "string of accrual system address")
	flag.Parse()

	err := envconfig.Process("", &cb)

	return cb, err
}
