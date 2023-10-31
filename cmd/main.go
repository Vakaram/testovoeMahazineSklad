package main

import "testovoezadanie1/internal/app"

func main() {
	config := config.New() // здесь должна строка появиться
	app := app.New(app.Config{
		Address:          "127.0.0.1:8080",
		ConnectionString: config.stringConnectionDB,
	})
	app.Start()
}
