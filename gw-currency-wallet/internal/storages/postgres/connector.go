package postgres

import (
	"flag"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBOptions struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int
}

func NewConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// после инициализации вызвать flag.Parse()
func InitFlags(DBcfg *DBOptions) {
	flag.StringVar(&DBcfg.DBUser, "POSTGRES_USER", os.Getenv("POSTGRES_USER"), "имя пользователя postgres")
	flag.StringVar(&DBcfg.DBPassword, "POSTGRES_PASSWORD", os.Getenv("POSTGRES_PASSWORD"), "пароль пользователя postgres")
	flag.StringVar(&DBcfg.DBName, "POSTGRES_DB", os.Getenv("POSTGRES_DB"), "название базы данных postgres")
	flag.StringVar(&DBcfg.DBHost, "POSTGRES_HOST", "localhost", "хост сервера postgres")
	flag.IntVar(&DBcfg.DBPort, "POSTGRES_PORT", 5432, "порт сервера postgres")
}
