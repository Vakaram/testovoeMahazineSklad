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

// RackGoods стелажи на которых лежит товар
type RackGoods struct {
	RackID  int  `db:"rack_id"`
	GoodsID int  `db:"goods_id"`
	Is_main bool `db:"is_main"`
}
