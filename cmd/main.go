package main

import (
	"fmt"
	"github.com/Vakaram/testovoeMahazineSklad/internal/app"
	"github.com/Vakaram/testovoeMahazineSklad/internal/storage"
	"testing"
)

func main() {
	// Получаем конфиг
	configStore := storage.ParseConfigDB()
	// Создаем pool для бд
	myStore := storage.New(configStore)
	//Создание таблиц в базе данных и схем
	err := myStore.InitTable()
	testing.Init()
	if err != nil {
		fmt.Printf("Ошибка при создание таблиц: %s ", err.Error())
		return
	}

	app.Start(myStore)
	fmt.Println("Старт программы ")
}
