package app

import (
	"bufio"
	"fmt"
	"github.com/Vakaram/testovoeMahazineSklad/internal/storage"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type app struct {
	Store *storage.Store
}

// NewApp создает новое приложение с конфигом для БД
func NewApp() *app {
	// Получаем конфиг
	configStore := storage.ParseConfigDB()
	// Создаем pool для бд
	myStore := storage.New(configStore)
	//Создание таблиц в базе данных и схем
	err := myStore.InitTable()
	if err != nil {
		fmt.Printf("Ошибка при создание таблиц: %s ", err.Error())

	}
	app := &app{
		Store: myStore,
	}
	return app
}

// Start запуск приложение и ожидание ввода
func Start(a *app) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		//разобрали ответ на id заказов которые нужно чекнуть в бд
		_, err := SplitRequest(text)
		if err != nil {
			logrus.Panic(err)
		}
		//получаем инфо из orders_goods

		////todo сделать массив элементов которые есть в бд и которых нет в бд и их выносить в отдельный овтет после списка ниже будет красиво

	}
}

// SplitRequest сплитует заказы по запятой формирует стрктуру
func SplitRequest(text string) ([]storage.RequestedOrders, error) {
	//todo сделать проверку а цифра ли это?
	nums := strings.Split(text, ",")
	var requestNums []storage.RequestedOrders
	for _, v := range nums {
		numInt, _ := strconv.Atoi(strings.TrimSpace(v))
		requestNums = append(requestNums, storage.RequestedOrders{Num: numInt})
	}
	for _, orderInt := range requestNums {
		fmt.Printf("Заказ такой :%d \n ", orderInt.Num)
	}

	return requestNums, nil
}
