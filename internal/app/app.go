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

func Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		// Обработка введенного текста
		//todo получить запросы по сплиту , запятой дальше придумаем че делать
		//тут создадим и отправим в  фукнцию для дальнейшей работы и логики
		_, err := SplitRequest(text)
		if err != nil {
			logrus.Panic(err)
		}
		//fmt.Printf("Вот введеный вами текст : %s", text)
	}
}

func SplitRequest(text string) ([]storage.RequestedOrders, error) {
	nums := strings.Split(text, ",")
	var requestNums []storage.RequestedOrders
	for _, v := range nums {
		numInt, _ := strconv.Atoi(strings.TrimSpace(v))
		requestNums = append(requestNums, storage.RequestedOrders{Num: numInt})
	}
	for _, orderInt := range requestNums {
		fmt.Printf("Заказ такой : %d \n ", orderInt.Num)
	}
	return requestNums, nil
}
