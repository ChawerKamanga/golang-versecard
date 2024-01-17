package main

import "github.com/joho/godotenv"

func loadConfig() error {
	return godotenv.Load()
}
