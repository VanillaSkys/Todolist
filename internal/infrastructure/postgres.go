package infrastructure

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitPostgres() {
	host := viper.GetString("database.postgres.host")
	user := viper.GetString("database.postgres.user")
	password := viper.GetString("database.postgres.password")
	name := viper.GetString("database.postgres.name")
	port := viper.GetInt("database.postgres.port")
	sslmode := viper.GetString("database.postgres.sslmode")
	timezone := viper.GetString("database.postgres.timezone")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, name, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to Postgres : ", err)
	}

	Db = db
}
