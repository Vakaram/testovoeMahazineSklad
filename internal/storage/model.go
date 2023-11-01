package storage

// Goods товары
type Goods struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Orders заказы
type Orders struct {
	ID  int `db:"id"`
	Num int `db:"num"`
}

// OrdersGoods связь таблиц ордер и товары
type OrdersGoods struct {
	OrdersID int `db:"orders_id"`
	GoodsID  int `db:"goods_id"`
	Sum      int `db:"sum"`
}

// Rack стелажи есть главные bool
type Rack struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
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
