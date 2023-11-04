package storage

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"os"
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
	// наполнение данными таблиц отключать после первого прогона если не нужны данные иначе будет ошибка ну значит данные уже есть так что не страшно но могут быть дубляжи
	//err = st.addDate()
	//if err != nil {
	//	fmt.Printf("Ошибка при заполнение данных : %s ", err)
	//}

	return nil
}

// addDate для наполнения данными
func (st *Store) addDate() error {
	// Наполнение таблицы goodsData
	goodsData := []Goods{
		{Name: "Ноутбук"},
		{Name: "Телевизор"},
		{Name: "Телефон"},
		{Name: "Компьютер"},
		{Name: "Часы"},
		{Name: "Микрофон"},
	}

	for _, g := range goodsData {
		_, err := st.pool.Exec(context.Background(), `
  INSERT INTO store.goods (name)
  VALUES ($1)
 `, g.Name)
		if err != nil {
			return fmt.Errorf("Failed to insert data into Goods table: %v", err)
		}
	}

	// добавление стелажей
	rackData := []Rack{
		{Name: "А"},
		{Name: "Б"},
		{Name: "В"},
		{Name: "Ж"},
		{Name: "З"},
	}

	for _, g := range rackData {
		_, err := st.pool.Exec(context.Background(), `
  INSERT INTO store.rack (name)
  VALUES ($1)
 `, g.Name)
		if err != nil {
			return fmt.Errorf("Failed to insert data into Rack table: %v", err)
		}
	}

	//Добавление заказов
	ordersData := []Orders{
		{Num: 10},
		{Num: 11},
		{Num: 14},
		{Num: 15},
	}

	for _, g := range ordersData {
		_, err := st.pool.Exec(context.Background(), `
  INSERT INTO store.orders (num)
  VALUES ($1)
 `, g.Num)
		if err != nil {
			return fmt.Errorf("Failed to insert data into Orders table: %v", err)
		}
	}

	//Связь заказы товар и кол-во
	ordersGoodsData := []OrdersGoods{
		//{OrdersID: 1,Rack_id:,Is_main: GoodsID: 1, Sum: 2}, //1 id  - это 10 заказ
		{OrdersID: 1, GoodsID: 3, Sum: 1},
		{OrdersID: 1, GoodsID: 6, Sum: 1},
		{OrdersID: 2, GoodsID: 2, Sum: 3}, //2 id - это 11 заказ и тд
		{OrdersID: 3, GoodsID: 1, Sum: 3},
		{OrdersID: 3, GoodsID: 4, Sum: 4},
		{OrdersID: 4, GoodsID: 5, Sum: 1},
	}

	for _, g := range ordersGoodsData {
		_, err := st.pool.Exec(context.Background(), `
  		INSERT INTO store.orders_goods
    	( orders_id,goods_id,sum)
 		VALUES ($1,$2,$3)`,
			g.OrdersID, g.GoodsID, g.Sum)
		if err != nil {
			return fmt.Errorf("Failed to insert data into Orders table: %v", err)
		}
	}

	goodsRackData := []GoodsRacks{
		{1, 1, true},
		{2, 1, true},
		{3, 2, true},
		{3, 5, false},
		{3, 3, false},
		{4, 4, true},
		{5, 4, true},
		{5, 1, false},
		{6, 4, true},
	}

	for _, g := range goodsRackData {
		_, err := st.pool.Exec(context.Background(), `
  		INSERT INTO store.goods_racks
    	(goods_id ,rack_id ,is_main)
 		VALUES ($1,$2,$3)`,
			g.GoodsID, g.RackID, g.IsMain)
		if err != nil {
			return fmt.Errorf("Failed to insert data into Orders table: %v", err)
		}
	}

	logrus.Info("Заполнили всеми нужными данными")
	return nil
}

//todo сделать так чтобы эта функция дополнительно тех заказов которые есть в одну группу собирала а те которые есть в другую

// ChecOrderInDB проверяет есть ли такой заказ вообще
func (st *Store) ChecOrderInDB(id int) (bool, error) {
	var exists bool // существует?
	err := st.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM store.orders WHERE num = $1)", id).Scan(&exists)
	if err != nil {
		// обработка ошибки
		logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
		return false, err
	}
	if exists {
		logrus.Info("Запись существует по ID")
		return true, nil
	} else {
		logrus.Info("Запись не существует по ID")
		return false, nil
	}
}

// Запрос по id Order наши товары и их кол-во в заказе
func (st *Store) PoluchaemSpisokSIdOrderAndRack(id []int) (idGoods []Orders, err error) {
	var goodsIDsum []Orders

	for _, v := range id {
		check, _ := st.ChecOrderInDB(v)
		if check == true {
			//fmt.Printf("True вижу для id %d", v)
			//Получаем соответствующие id для наших заказов
			idOrders, err := st.GiveIdInOrders(id)
			if err != nil {
				fmt.Printf("Ошибочка =_)1 ", err)
				return nil, err
			}
			//получим все товары связанные с этим id
			_, err = st.GiveAllGoodsAndSumInOrders(idOrders)
			if err != nil {
				fmt.Printf("Ошибочка =_)2 ", err)
				return nil, err
			}
			//получаем из orders_goods_rask id товаров и их sum в стркутуре

		} else {
			return nil, err
		}
	}
	fmt.Printf("Структуру такая %v", goodsIDsum)

	return goodsIDsum, nil
}

// GiveIdInOrders возвращает ID наших ордеров обратившись в ORDERS
func (st *Store) GiveIdInOrders(id []int) (idOrders []int, err error) {
	var ids []int
	for _, v := range id {
		rows, err := st.pool.Query(context.Background(), "SELECT id FROM store.orders WHERE num = $1", v)
		if err != nil {
			// обработка ошибки
			logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
			return nil, err

		}

		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				// обработка ошибки
				logrus.Error("Ошибка при сканировании строки: ", err)
				return nil, err
			}
			ids = append(ids, id)
		}
	}
	fmt.Printf("\nВот получил id по Orders их id: %d \n", ids[0])
	return ids, err
}
func (st *Store) GiveAllGoodsAndSumInOrders(idOrders []int) (goodsAmdSumInOrders []OrdersGoods, err error) {

	var idGoodsAndSuminOrders []OrdersGoods
	for _, v := range idOrders {
		rows, err := st.pool.Query(context.Background(), "SELECT orders_id,goods_id,sum FROM store.orders_goods WHERE orders_id = $1", v)
		fmt.Printf("Вот запрос твой %s ", rows)
		if err != nil {
			// обработка ошибки
			logrus.Error("Ошибка при выполнении запроса к базе данных: ", err)
			return nil, err

		}

		for rows.Next() {
			var orders_id int
			var goods_id int
			var sum int
			err := rows.Scan(&goods_id, &sum)
			if err != nil {
				// обработка ошибки
				logrus.Error("Ошибка при сканировании строки: ", err)
				return nil, err
			}
			idGoodsAndSuminOrders =
				append(idGoodsAndSuminOrders, OrdersGoods{OrdersID: orders_id, GoodsID: goods_id, Sum: sum})
		}
	}
	fmt.Printf("\n Получил вот такие товары для заказа их id = %v \n", idGoodsAndSuminOrders)
	fmt.Printf("\n zzzzz= %s ,%s \n", idGoodsAndSuminOrders[0].GoodsID, idGoodsAndSuminOrders[0].Sum)

	return nil, err
}
