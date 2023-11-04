package storage

// Orders заказы
type Orders struct {
	ID  int `db:"id"`
	Num int `db:"num"`
}

// OrdersGoods связь таблиц ордер и товары
type OrdersGoods struct {
	OrdersID int  `db:"orders_id"`
	Rack_id  int  `db:"rack_id"`
	Is_main  bool `db:"is_main"`
	GoodsID  int  `db:"goods_id"`
	Sum      int  `db:"sum"`
}

// GoodsRacks стелажи на которых лежит товар
type GoodsRacks struct {
	GoodsID int  `db:"goods_id"`
	RackID  int  `db:"rack_id"`
	IsMain  bool `db:"is_main"`
}

// ордеры которые приходят чтобы получить инфу
type RequestedOrders struct {
	Num int
}

//// Страница сборки заказов ответ
//type OrderAssemblyPage struct {
//	RackPage  Rack
//	GoodsPage Goods
//	OrdersNum Orders
//	OrderSum  OrdersGoods
//}

// Страница сборки заказов ответ
type OrderAssemblyPage struct {
	RackPage    Rack        //стелаж
	Description Description // описание что лежит в стелаже для удобного парса

}

type Description struct {
	GoodsPage Goods       //id товара и название
	OrdersNum Orders      //номер заказа
	OrderSum  OrdersGoods //кол-во в заказе штук
}

// Итоговая модель при парсе которой получится показать страницу заказа
type FullInfoPage struct {
	IdOrderDB int
	NumOrder  int
	Goods     []Goods // товаров может быть несколько
}

// Goods товары
type Goods struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Sum  int    `db:"sum"`
	Rack []Rack // содержание для товара его стелажа главного и допов тоже может быть несколько
}

// Rack стелажи есть главные bool
// Будут содержаться несколько вариантов []Rack
type Rack struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	IsMain bool   `db:"is_main"`
}
