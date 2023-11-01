package storage

import (
	"context"
	"github.com/Vakaram/testovoeMahazineSklad/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
)

type storage struct {
	poll *pgxpool.Pool
}

type Config struct {
	ConnectionString string
}

func New(cfg Config) *storage {
	//проверка пула на ошибку
	poolNew, err := pgxpool.New(context.Background(), cfg.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	//всегда передавай ссылку на storage чтобы изменить его иначен не вернуь будет ошибка
	s := &storage{
		poll: poolNew,
	}

	return s
}

// init создание таблиц
func (s *storage) init() {
	// todo написать таблицы

}

func (s *storage) Create(ctx context.Context, user models.CreateUser) (models.User, error) {
	query := `
    INSERT INTO users (first_name, last_name, age, email)
    VALUES ($1, $2, $3, $4)
    RETURNING id, first_name, last_name, age, email
  `
	var newUser models.User
	err := s.poll.QueryRow(ctx, query, user.FirstName, user.LastName, user.Age, user.Email).Scan(&newUser.ID, &newUser.FirstName, &newUser.LastName, &newUser.Age, &newUser.Email)
	if err != nil {
		return models.User{}, err
	}

	return models.User{}, nil
}
func (s *storage) Update(ctx context.Context, id int, user models.UpdateUser) (models.User, error) {
	return models.User{}, nil
	// todo написать запрос
}
func (s *storage) FidByID(ctx context.Context, id int) (models.User, error) {
	return models.User{}, nil
	// todo написать запрос
}

// типо добавили новый метод
func (s *storage) FidByEmail(ctx context.Context, email string) (models.User, error) {
	return models.User{}, nil
	// todo написать запрос
}

func (s *storage) Delete(ctx context.Context, id int) error {
	return nil
	// todo написать запрос
}
