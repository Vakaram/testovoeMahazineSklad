package handler

import (
	"context"
	"github.com/Vakaram/testovoeMahazineSklad/internal/models" // все будут зависить от моделей но от друг друга зависеть не будут
	// суть сделать так чтобы мы обрабатывали и меняли юзера внутри ничего не меняем а вот самого юзера меняем и зависим
)

// UserService делаем интерфейс которому потом сделаем стркутруы с методами
// предпологаем что появится сервис у кторого будут вот такие 4 метода
type UserService interface {
	Create(ctx context.Context, user models.CreateUser) (models.User, error)
	Update(ctx context.Context, id int, user models.UpdateUser) (models.User, error)
	FidByID(ctx context.Context, id int) (models.User, error)
	Delete(ctx context.Context, id int) error
}
