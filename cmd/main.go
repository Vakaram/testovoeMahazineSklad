package main

import (
	"fmt"
	"github.com/Vakaram/testovoeMahazineSklad/internal/app"
)

func main() {
	//// Получаем конфиг
	//configStore := storage.ParseConfigDB()
	//// Создаем pool для бд
	//myStore := storage.New(configStore)
	// Создание приложения нового
	newApp := app.NewApp()
	////Создание таблиц в базе данных и схем
	//err := myStore.InitTable()
	//if err != nil {
	//	fmt.Printf("Ошибка при создание таблиц: %s ", err.Error())
	//	return
	//}
	app.Start(newApp)
	fmt.Println("Старт программы ")
}
