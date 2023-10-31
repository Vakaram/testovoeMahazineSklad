package service

import (
	"context"
	"errors"
	"github.com/Vakaram/sterAuto/internal/models"
)

// тут описываем методы для утинной типиазции и тогда она сработает
// теперь сюда передадим интерфейс
type service struct {
	store userStorage
}

type Config struct {
	Store userStorage
}

func New(cfg Config) *service {
	s := &service{
		store: cfg.Store,
	}
	return s
}

func (s *service) Create(ctx context.Context, newUser models.CreateUser) (models.User, error) {
	//такие раз и в легку проверку по email сделали
	user, _ := s.store.FidByEmail(ctx, newUser.Email)
	if user.Id == 0 {
		return models.User{}, errors.New("email zaniat , email already exist")

	}
	return s.store.Create(ctx, newUser)
}

func (s *service) FidByID(ctx context.Context, id int) (models.User, error) {
	//а не рекурсия ли это ?ааа тут мы типо передаем в store а у него там свои такие же методы будут уухххх
	return s.store.FidByID(ctx, id)
}
func (s *service) Delete(ctx context.Context, id int) error {
	//перед удалением надо узнать а есть ли вообще такой пользователь

	user, _ := s.store.FidByID(ctx, id)
	if user.Id == 0 {
		//а вот и вернули ошибку из models тк мы её создали и чтобы строки не писать по сути в ней и есть строка
		return models.ErrUserNotFound
	}
	return s.store.Delete(ctx, id)
}

func (s *service) Update(ctx context.Context, id int, updUser models.UpdateUser) (models.User, error) {
	user, _ := s.store.FidByID(ctx, id)
	if user.Id == 0 {
		//а вот и вернули ошибку из models тк мы её создали и чтобы строки не писать по сути в ней и есть строка
		return models.User{}, models.ErrUserNotFound
	}

	//а вот уже возвращаем то что вернет нам бд
	return s.store.Update(ctx, id, updUser)
}
