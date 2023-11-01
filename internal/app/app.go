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

// Start запуск приложение и ожидание ввода
func Start(myStore *storage.Store) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		//разобрали ответ на id заказов которые нужно чекнуть в бд
		requestOrdersInput, err := SplitRequest(text)
		if err != nil {
			logrus.Panic(err)
		}
		//теперь получим наши ID товаров связанные с заказом
		massivOrders := []int{}
		for _, v := range requestOrdersInput {
			massivOrders = append(massivOrders, v.Num)
		}

		order, err := myStore.ChecOrder(massivOrders)
		//logrus.Info("Получили заказы функцию вызвали")
		if err != nil {
			logrus.Debug(err)
			return
		}
		for _, g := range order {
			fmt.Printf("Получил вот такой товар по поиску  : %s ", g)
		}
	}
}

// SplitRequest сплитует заказы по запятой формирует стрктуру
func SplitRequest(text string) ([]storage.RequestedOrders, error) {
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
