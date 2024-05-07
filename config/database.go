package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type databaseConfiguration struct {
	Host     string `env:"POSTGRES_HOST,required" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT,required" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	Database string `env:"POSTGRES_DATABASE" envDefault:"postgres"`
}

func (config *databaseConfiguration) ToDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.Host, config.Port, config.User, config.Password, config.Database)
}

func DatabaseConnection() *gorm.DB {
	config := databaseConfiguration{}
	if err := env.Parse(&config); err != nil {
		fmt.Printf("%+v\n", err)
	}
	connection, err := gorm.Open(postgres.Open(config.ToDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return connection
}
