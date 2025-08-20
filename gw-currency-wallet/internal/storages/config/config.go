package config

type ServerConfig struct {
	Port int // порт для запуска сервера

}

type DBConfig struct {
	Port         int
	Host         string
	User         string
	Password     string
	DatabaseName string
}
