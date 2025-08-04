package config

type Configuration struct {
	// database fields
	DB_USER   string
	DB_PASS   string
	DB_DBNAME string
	DB_DBHOST string
	DB_PORT   string
}

func GetConfig() Configuration {
	configuration := Configuration{}
	configuration.DB_USER = "dev"
	configuration.DB_PASS = "dev"
	configuration.DB_DBNAME = "dev"
	configuration.DB_DBHOST = "localhost"
	configuration.DB_PORT = "5432"
	return configuration
}
