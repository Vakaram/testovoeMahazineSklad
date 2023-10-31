package service

// делается чтобы user напрямую не мог вызывать базу данных
import (
	"context"
	"github.com/Vakaram/testovoeMahazineSklad/internal/models"
)

type userStorage interface {
	Create(ctx context.Context, user models.CreateUser) (models.User, error)
	Update(ctx context.Context, id int, user models.UpdateUser) (models.User, error)
	FidByID(ctx context.Context, id int) (models.User, error)
	//типо добавили новый метод
	FidByEmail(ctx context.Context, email string) (models.User, error)

	Delete(ctx context.Context, id int) error
}
