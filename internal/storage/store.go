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

	// Создание таблицы orders_goods
	_, err = st.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS store.orders_goods (
 orders_id INT REFERENCES store.orders(id),
 goods_id INT REFERENCES store.goods(id),
 sum INT,
 PRIMARY KEY (orders_id, goods_id)
)`)
	if err != nil {
		log.Fatalf("Failed to create orders_goods table: %v", err)
	}

	// Создание таблицы rack
	_, err = st.pool.Exec(context.Background(), `
	CREATE TABLE 
	IF NOT EXISTS store.rack (
 	id SERIAL PRIMARY KEY,
	name VARCHAR
)`)
	if err != nil {
		log.Fatalf("Failed to create rack table: %v", err)
	}

	// Создание таблицы rack_goods
	_, err = st.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS store.rack_goods (
 rack_id INT REFERENCES store.rack(id),
 goods_id INT REFERENCES store.goods(id),
 is_main BOOLEAN,
 PRIMARY KEY (rack_id, goods_id)
)`)
	if err != nil {
		log.Fatalf("Failed to create rack_goods table: %v", err)
	}

	logrus.Info("Создали таблицы")
	// наполнение данными таблиц
	st.addDate()

	return nil
}

func (st *Store) addDate() error {
	// Наполнение таблицы goods
	goods := []Goods{
		{Name: "Ноутбук"},
		{Name: "Телевизор"},
		{Name: "Телефон"},
		{Name: "Компьютер"},
		{Name: "Часы"},
		{Name: "Микрофон"},
	}

	for _, g := range goods {
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

	logrus.Info("Заполнили всеми нужными данными")
	return nil
}
