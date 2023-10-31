package handler

import (
	"encoding/json"
	"errors"
	"github.com/Vakaram/testovoeMahazineSklad/internal/models"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type handler struct {
	UserService UserService //тот сервис для usera сюда теперь запишем
	router      *mux.Router
}

type Config struct {
	UserService UserService
}

func New(cfg Config) *handler {
	h := &handler{
		router:      mux.NewRouter(),
		UserService: cfg.UserService,
	}
	return h

}

// Handler Описываем тут ручки наши
func (h *handler) Handler() http.Handler {
	h.router.HandleFunc("/user", h.createUser).Methods(http.MethodPost)
	h.router.HandleFunc("/user/{id}", h.getUser).Methods(http.MethodGet)
	h.router.HandleFunc("/user/{id}", h.updateUser).Methods(http.MethodPatch)
	h.router.HandleFunc("/user/{id}", h.deleteUser).Methods(http.MethodDelete)
	return h.router
}

// todo добавить логирование во все места
func (h *handler) createUser(res http.ResponseWriter, req *http.Request) {
	// для создания пользователя нужен контекст мы его берем из запроса
	ctx := req.Context()
	body, err := io.ReadAll(req.Body) // читает все тело
	// обязательно закрываем соедениние
	defer req.Body.Close()
	if err != nil {
		// вернем ошибку
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var newUser models.CreateUser
	if err := json.Unmarshal(body, &newUser); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.UserService.Create(ctx, newUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// вызвали функцию для возвращааения пользователя в json
	returnUser(res, user)

}

func (h *handler) getUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, err := getID(req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := h.UserService.FidByID(ctx, id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// вызвали функцию для возвращааения пользователя в json
	returnUser(res, user)

}

func (h *handler) updateUser(res http.ResponseWriter, req *http.Request) {

	// для создания пользователя нужен контекст мы его берем из запроса
	ctx := req.Context()

	id, err := getID(req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(req.Body) // читает все тело
	// обязательно закрываем соедениние
	defer req.Body.Close()
	if err != nil {
		// вернем ошибку
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var updUser models.UpdateUser
	if err := json.Unmarshal(body, &updUser); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := h.UserService.Update(ctx, id, updUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// вызвали функцию для возвращааения пользователя в json
	returnUser(res, user)

}

func (h *handler) deleteUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, err := getID(req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.UserService.Delete(ctx, id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	//возвращаем код удаления
	res.WriteHeader(http.StatusNoContent)

}

// чтобы сократить код
func returnUser(res http.ResponseWriter, user models.User) {
	respounse, err := json.Marshal(user)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.Write(respounse)
	res.WriteHeader(http.StatusCreated)

}

func getID(req *http.Request) (int, error) {
	idstring, ok := mux.Vars(req)["id"]
	if !ok {
		return 0, errors.New("id is required")
	}

	// atoi возвращает сразу ошибку и инт так что передаем только его
	return strconv.Atoi(idstring)

}
