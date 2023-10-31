package main

import (
	"github.com/Vakaram/testovoeMahazineSklad/internal/app"
	"github.com/Vakaram/testovoeMahazineSklad/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	configMain := config.New() // здесь должна строка появиться
	myApp := app.New(app.Config{
		Address:          configMain.Address,
		ConnectionString: configMain.DatabaseURL,
	})
	logrus.Info("Программа запущена")
	myApp.Start()

}
