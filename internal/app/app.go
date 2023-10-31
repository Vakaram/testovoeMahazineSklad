package app

import (
	"github.com/Vakaram/testovoeMahazineSklad/internal/handler"
	"github.com/Vakaram/testovoeMahazineSklad/internal/service"
	"github.com/Vakaram/testovoeMahazineSklad/internal/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

// любой пакет должен начинаться со структуры

type app struct {
	address string
	handler http.Handler
	logger  *logrus.Logger //добавляем логгер

}
type Config struct {
	Address          string
	ConnectionString string // сюда должны передать строку // todo
}

func New(cfg Config) *app {
	store := storage.New(storage.Config{
		cfg.ConnectionString,
	})

	//перед этим проинициализурем сервис для юзера
	// ниже не срабатывало тк утиная типизация не срабатывала
	UserService := service.New(service.Config{
		// передадим в store store который инициализурем выше
		Store: store,
	})
	//надо добавить handler  в проект вызовим new у папки handler
	h := handler.New(handler.Config{UserService})

	vozvratim := &app{
		address: cfg.Address,
		handler: h.Handler(),  // вызвали функцию которая станет доаступна в h
		logger:  logrus.New(), // logrus.NewStore() - это встроенно в логрус а не нами написано
	}
	return vozvratim
}

func (a *app) Start() {
	http.ListenAndServe(a.address, a.handler)
}
