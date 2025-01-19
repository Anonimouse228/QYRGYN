package config

import "os"

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL") // Example: "host=localhost port=5432 user=postgres dbname=social password=yourpassword sslmode=disable"
}

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func GetPort() string {
	return os.Getenv("PORT")
}
