package config

type DBConfig struct {
	DSN string
}

type BlogConfig struct {
	DB DBConfig
}
