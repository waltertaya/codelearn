package initialisers

import "github.com/joho/godotenv"

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}
