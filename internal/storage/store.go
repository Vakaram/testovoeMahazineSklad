package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	//"strconv"
)

type Store struct {
	pool *pgxpool.Pool
}

type Config struct {
	DatabaseURL string `yaml:"database_url"`
}

func ParseConfigDB() Config {
	configPath := "./internal/config/local.yaml" // os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}

func New(cfg Config) *Store {
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	s := &Store{
		pool: dbPool,
	}
	logrus.Info("Подключились к БД POOL ")
	return s
}

// InitTable наполнение тестовыми данными бд
func (st *Store) InitTable() error {
	// Создание схемы
	_, err := st.pool.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS store")
	if err != nil {
		fmt.Println("Ошибка при создании схемы:", err)
		return err
	}
	// Создание таблицы goods
	_, err = st.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS store.goods (
 id SERIAL PRIMARY KEY,
 name VARCHAR
 )`)
	if err != nil {
		log.Fatalf("Failed to create goods table: %v", err)
	}

	// Создание таблицы orders
	_, err = st.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS store.orders (
 id SERIAL PRIMARY KEY,
 num INT
)`)
	if err != nil {
		log.Fatalf("Failed to create orders table: %v", err)
	}

	// Создание таблицы rack
	_, err = st.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS store.rack (
 	id SERIAL PRIMARY KEY,
	name VARCHAR,
	is_main BOOLEAN
)`)
	if err != nil {
		log.Fatalf("Failed to create rack table: %v", err)
	}
	// Создание таблицы orders_goods
	_, err = st.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS store.orders_goods (
    id SERIAL,
    orders_id INT REFERENCES store.orders(id),
    goods_id INT REFERENCES store.goods(id),
    sum INT,
    PRIMARY KEY (orders_id, goods_id)
)`)
	if err != nil {
		log.Fatalf("Failed to create orders_goods table: %v", err)
	}

	// Создание таблицы extra_rack
	_, err = st.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS store.extra_rack  (
    id SERIAL,
  	goods_id INT REFERENCES store.goods (id),
  	rack_id INT REFERENCES store.rack (id),
  	PRIMARY KEY (rack_id, goods_id)

)`)
	if err != nil {
		log.Fatalf("Failed to create extra_rack table: %v", err)
	}

	logrus.Info("Создали таблицы")
	return nil
}

// ChecOrderInDB проверяет есть ли такой заказ вообще и возвращает 2 списка в котором есть и в котором нет данных
func (st *Store) ChecOrderInDB(id int) (bool, error) {
	var exists bool // существует?
	err := st.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM store.orders WHERE num = $1)", id).Scan(&exists)
	if err != nil {
		// обработка ошибки
		logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
		return false, err
	}
	if exists {
		//logrus.Info("Запись существует по ID верну тру")
		return true, nil
	} else {
		//logrus.Info("Запись не существует по ID")
		return false, nil
	}
}

// FullInfoPage вернет нам id оредров и начнет строить струтктуру FullInfoPage
func (st *Store) FullInfoPage(id []int) (readyAnswer []FullInfoPage, ordersDontExist []int, err error) {
	//то что будем возвращать
	var idOrdersAndGoodsReturn []FullInfoPage
	var noDataOrderReturn []int // если не будет при проверке то добавим сюда

	for _, v := range id {
		// вывозвим проверку есть ли такой заказ в базе
		check, err := st.ChecOrderInDB(v)
		if err != nil {
			return nil, nil, err
		}
		if check == true {
			//fmt.Print("Зашли в Тру ")
			// получим id этого заказа для таблицы orders_goods
			idInOrder, err := st.ZaprosOrderID(v)
			if err != nil {
				return nil, nil, err
			}
			//fmt.Printf("Получаю вот такой список ID в ZaprosOrderID %d ", idInOrder)
			idOrdersAndGoodsReturn = append(idOrdersAndGoodsReturn, FullInfoPage{IdOrderDB: idInOrder, NumOrder: v})

		} else {
			fmt.Print("Зашли в Фолс ")
			//если такого заказа нет то добавим этот id в ordersDontExist
			noDataOrderReturn = append(noDataOrderReturn, v)
		}
	}
	//fmt.Printf("Структуру если Заказ был такая до  %v", idOrdersAndGoodsReturn)
	//fmt.Printf("Спискок num если заказа нет в бд %v", noDataOrderReturn)

	//Далее предлагаю обогатить нашу структуру доп данными =)
	addOrdresdForFULL, err := st.AddOrdersInFullPage(idOrdersAndGoodsReturn)
	//fmt.Printf("Структуру если Заказ был такая После11111  %v", addOrdresdForFULL)
	//теперь обоготим данными о стелажах rack
	idOrdersAndGoodsReturn, err = st.AddRackInFullPage(addOrdresdForFULL)
	//fmt.Printf("Структуру если Заказ был такая После4444  %+v \n", idOrdersAndGoodsReturn)

	return idOrdersAndGoodsReturn, noDataOrderReturn, nil

}

// ZaprosOrderID возвращает ID наших ордеров обратившись в ORDERS
func (st *Store) ZaprosOrderID(idNum int) (idOrders int, err error) {
	var zeroIdErr int
	//fmt.Printf("\n Вызвали функцию Запрос Ордера ИД и получил на вход такой инт %d ", idNum)
	row, err := st.pool.Query(context.Background(), "SELECT id FROM store.orders WHERE num = $1", idNum)
	if err != nil {
		// обработка ошибки
		logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
		return zeroIdErr, err

	}

	var id int
	if row.Next() {
		err = row.Scan(&id)
		if err != nil {
			fmt.Printf("\n Ошибка при скане ,%v,%v", err)
			return zeroIdErr, err
		}
	} else {
		fmt.Println("Отсутствуют результаты запроса в сканирование ")
		return zeroIdErr, errors.New("Отсутствуют результаты запроса")
	}

	//fmt.Printf("Вижу вот такой id %d для num %d", id, idNum)
	return id, err
}

func (st *Store) AddOrdersInFullPage(fullPage []FullInfoPage) (fullPageAddGoods []FullInfoPage, err error) {
	var fullPageReturn []FullInfoPage //
	// сделаем запрос в orders_goods по id из fullInfoPage и получим там наши товары
	//elem == 0
	for _, v := range fullPage {
		//fmt.Print("искать буду для id ", v.IdOrderDB)
		rows, err := st.pool.Query(context.Background(), ""+
			//переименовать store.orders_goods_rask на store.orders_goods
			"SELECT goods_id,sum FROM store.orders_goods_rask WHERE orders_id = $1",
			v.IdOrderDB)

		//fmt.Printf("Запрос выглядит так SELECT goods_id,sum FROM store.orders_goods WHERE orders_id = $1,%d или %d", v.IdOrderDB)

		if err != nil {
			// обработка ошибки
			logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
			return nil, err

		}
		//создаем новый фулл падже и туда прокидываем наши данные которые знаем
		// обявляем goods который засуним в fullPage
		var GoodsAdd []Goods
		//что заберем из запроса
		for rows.Next() {
			var goods Goods
			err = rows.Scan(&goods.ID, &goods.Sum)
			if err != nil {
				fmt.Printf("\n Ошибка при скане ,%v,%v", err)
				return nil, err
			}
			//в цикле обогащаем данные через next собираем все товары подходящшие по id
			// Дополнительный запрос на получение имени товара
			err = st.pool.QueryRow(context.Background(), "SELECT name FROM store.goods WHERE id = $1", goods.ID).Scan(&goods.Name)
			if err != nil {
				fmt.Printf("\n Ошибка при получении имени товара ,%v,%v", err)
				return nil, err
			}
			GoodsAdd = append(GoodsAdd, goods)
		}
		//fmt.Printf("После обогащения циклом :v ", GoodsAdd)
		//после сканирования и обогащения Orders то мы обогатим FullPageReturn для [0] нулевого жлемента вооот
		fullPageReturn = append(fullPageReturn, FullInfoPage{IdOrderDB: v.IdOrderDB, NumOrder: v.NumOrder, Goods: GoodsAdd})

	}

	return fullPageReturn, nil

}

// AddRackInFullPage Добавляет стелажи
func (st *Store) AddRackInFullPage(fullPage []FullInfoPage) (fullPageAddGoods []FullInfoPage, err error) {
	//Пришла вот такая стрктуруа
	//fmt.Printf("Структуру если Заказ был такая После3333  %v", fullPage)

	var fullPageReturn []FullInfoPage //
	// сделаем запрос в extra_rack получим все связи с таблицей -- вынесим в отдельную функцию
	// Далее по этим связям наполним данные rack для FullInfoPage
	// Хитро ренджим по элементам как по дереву
	for _, v := range fullPage {
		var goodsAdd []Goods
		//fmt.Printf("ПРоход для элемента Fill Page : %s ", strconv.Itoa(v.Goods[i].ID))
		for _, vGods := range v.Goods {
			idRack, err := st.GiveExtraRackRack_id(vGods.ID)
			if err != nil {
				return nil, err
			}
			//получили id_rack для товара и теперь получим его стелажи Rack
			//fmt.Printf("Id для поиска в rack  :%v ", idRack)
			rackForGods, err := st.GiveRackByIdRackID(idRack)
			if err != nil {
				return []FullInfoPage{}, nil
			}
			goodsAdd = append(goodsAdd, Goods{ID: vGods.ID, Name: vGods.Name, Sum: vGods.Sum, Rack: rackForGods})
			//fmt.Printf("Goods получился вот такой : %v", goodsAdd)
			// Осталось только добавить это к общему завершаещему массиву
		}
		fullPageReturn = append(fullPageReturn, FullInfoPage{v.IdOrderDB, v.NumOrder, goodsAdd})
	}
	//fmt.Printf("Вижу 123 fullPageRetur %v", fullPageReturn)

	return fullPageReturn, nil
}

// GiveExtraRackRack_id вернет связанные с id товара все его места в на стелажах
func (st *Store) GiveExtraRackRack_id(goodsId int) (extraRackId []int, err error) {
	var extraRackIdReturn []int
	// делать запрос на получение массива int из extraRack
	rows, err := st.pool.Query(context.Background(), ""+
		//переименовать store.orders_goods_rask на store.orders_goods
		"SELECT rack_id FROM store.extra_rack WHERE goods_id = $1",
		goodsId)

	if err != nil {
		// обработка ошибки
		logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
		return nil, err

	}
	//создаем новый фулл падже и туда прокидываем наши данные которые знаем
	// обявляем goods который засуним в fullPage
	//что заберем из запроса
	for rows.Next() {
		var goods_id int
		err = rows.Scan(&goods_id)
		if err != nil {
			fmt.Printf("\n Ошибка при скане ,%v,%v", err)
			return nil, err
		}
		//в цикле обогащаем данные через next собираем все товары подходящшие по id
		extraRackIdReturn = append(extraRackIdReturn, goods_id)
	}
	//fmt.Printf("Id для поиска в rack  :%v ", extraRackIdReturn)
	//после сканирования и обогащения Orders то мы обогатим FullPageReturn для [0] нулевого жлемента вооот

	return extraRackIdReturn, nil

}

func (st *Store) GiveRackByIdRackID(rackID []int) (rackForGoods []Rack, err error) {
	// объявим структруу новую и туда добавим данные
	var rackForGoodsReturn []Rack
	for _, v := range rackID {
		rows, err := st.pool.Query(context.Background(), ""+
			//переименовать store.orders_goods_rask на store.orders_goods
			"SELECT id,name,is_main FROM store.rack WHERE id = $1",
			v)
		if err != nil {
			// обработка ошибки
			logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
			return []Rack{}, err
		}

		var rackValue Rack
		for rows.Next() {
			err = rows.Scan(&rackValue.ID, &rackValue.Name, &rackValue.IsMain)
			if err != nil {
				fmt.Printf("\n Ошибка при скане ,%v,%v", err)
				return []Rack{}, err
			}
		}
		rackForGoodsReturn = append(rackForGoodsReturn, rackValue)

	}
	//fmt.Printf("Вижу!!!!!!!! :", rackForGoodsReturn)
	return rackForGoodsReturn, err
}
