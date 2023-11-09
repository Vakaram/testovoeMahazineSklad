package app

// todo добавить запрос для того чтобы еще имя писалось в сткрутуре
import (
	"bufio"
	"fmt"
	"github.com/Vakaram/testovoeMahazineSklad/internal/storage"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
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
		intOrdersSplit, err := SplitRequest(text)
		if err != nil {
			logrus.Panic(err)
		}
		//получаем инфо из orders_goods стркутру Goods с заказами и данные по которым ответа нет
		readyAnswer, ordersDontExist, err := a.Store.FullInfoPage(intOrdersSplit)
		if err != nil {
			fmt.Printf("Ошибка в app " + err.Error())
			return
		}
		//fmt.Printf("Пришло 999 %+v", readyAnswer)

		//Тут будет формирование новой стуктуры
		itogStructForText, err := SortInRack(readyAnswer)
		//fmt.Printf("Пришло 123 %v", itogStructForText)
		//fmt.Printf("Заказы которые не нашли  %d", ordersDontExist)

		// Тут уже полученную структуру преобразуем в читаемый вид
		//передали структуру, данные которых нет по запросам, и те на которые готов ответ intOrdersSplit

		itog, err := BeautifulText(itogStructForText, ordersDontExist, intOrdersSplit)
		fmt.Printf(itog)
		////todo сделать массив элементов которые есть в бд и которых нет в бд и их выносить в отдельный овтет после списка ниже будет красиво
		//logrus.Info(readableAnswer)
	}
}

// SplitRequest сплитует заказы по запятой формирует стрктуру
func SplitRequest(text string) (orderNum []int, err error) {
	//todo сделать проверку а цифра ли это?
	nums := strings.Split(text, ",")
	var requestNums []int
	for _, v := range nums {
		numInt, _ := strconv.Atoi(strings.TrimSpace(v))
		requestNums = append(requestNums, numInt)
	}
	return requestNums, nil
}

func SortInRack(oldStruct []storage.FullInfoPage) (newStruct []storage.RackItog, err error) {
	// Создаем карту для группировки по главным стелажам
	racksMap := make(map[string][]storage.GoodsItog)

	// Обходим все элементы структуры Atest
	for _, fullInfo := range oldStruct {
		// Обходим все товары в заказе
		var falseRack string
		for _, goods := range fullInfo.Goods {
			for _, rack1 := range goods.Rack {
				if rack1.IsMain != true {
					falseRack += rack1.Name + ","
				}
			}
			// Проверяем, является ли стелаж главным
			for _, rack := range goods.Rack {

				if rack.IsMain {
					// Создаем объект GoodsItog
					goodsItog := storage.GoodsItog{
						IdGoods:   goods.ID,
						Name:      goods.Name,
						Order:     fullInfo.NumOrder,
						Sum:       goods.Sum,
						ExtraRack: falseRack,
					}

					// Добавляем объект GoodsItog в карту по имени главного стелажа
					racksMap[rack.Name] = append(racksMap[rack.Name], goodsItog)
				}
			}
			//очищаем список стелажей допов
			falseRack = ""
		}
	}

	// Создаем срез RackItog
	var rackItog []storage.RackItog

	// Создаем объекты RackItog на основе карты главных стелажей
	for rackName, goodsItog := range racksMap {
		rack := storage.RackItog{
			RackName:  rackName,
			GoodsItog: goodsItog,
		}
		rackItog = append(rackItog, rack)
	}
	// Сортируем структуру RackItog по полю RackName
	sort.Slice(rackItog, func(i, j int) bool {
		return rackItog[i].RackName < rackItog[j].RackName
	})
	//fmt.Printf("Вот структура итог %v", rackItog)

	return rackItog, nil
}

// преобразует структуру приходяющую в красивый текст
//можно возвращать если товар не найден тоже строку

func BeautifulText(itogForText []storage.RackItog, ordersDontExist []int, numRequest []int) (finish string, err error) {
	//===Стеллаж
	var numReq string
	for _, order := range numRequest {
		numReq += strconv.Itoa(order) + ","
	}
	// удаляет заданный символ справа
	numReq = strings.TrimRight(numReq, ",")

	var stringReturn string
	stringReturn = "=+=+=+=\nСтраница сборки заказов " + numReq + "\n"
	for _, v := range itogForText {
		var goods string
		for _, v2 := range v.GoodsItog {
			if v2.ExtraRack == "" {
				goods += v2.Name + "(id=" + strconv.Itoa(v2.IdGoods) + ")" + "\n" + "заказ " + strconv.Itoa(v2.Order) + ", " + strconv.Itoa(v2.Sum) + " шт" + "\n" + "\n"
			} else {
				goods += v2.Name + "(id=" + strconv.Itoa(v2.IdGoods) + ")" + "\n" + "заказ " + strconv.Itoa(v2.Order) + ", " + strconv.Itoa(v2.Sum) + " шт" + "\n" + "доп стеллаж: " + v2.ExtraRack + "\n" + "\n"

			}
		}
		stringReturn += "===Стелаж " + v.RackName + "\n" + goods
	}

	return stringReturn, nil
}
