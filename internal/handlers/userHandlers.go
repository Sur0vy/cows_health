package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

type UserHandle interface {
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Register(c echo.Context) error

	ResponseBadRequest(c echo.Context) error
}

type UserHandler struct {
	log         *logger.Logger
	userStorage storage.UserStorage
}

func NewUserHandler(us storage.UserStorage, log *logger.Logger) UserHandle {
	return &UserHandler{
		log:         log,
		userStorage: us,
	}
}

func (h *UserHandler) Login(c echo.Context) error {
	h.log.Info().Msgf("Handler IN: %v", c)
	defer h.log.Info().Msgf("Handler OUT: %v", c)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	defer c.Request().Body.Close()
	h.log.Info().Msg(string(body))
	var user storage.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var hash string
	hash, err = h.userStorage.GetUserHash(c.Request().Context(), user)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusUnauthorized)
			return c.NoContent(http.StatusUnauthorized)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	cookie := new(http.Cookie)
	cookie.Name = config.Cookie
	cookie.Value = hash
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	h.log.Info().Msgf("login success (cookie: %v)", cookie)
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = config.Cookie
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Time{}
	c.SetCookie(cookie)
	h.log.Info().Msg("logout success")
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) Register(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Register failed. Bad request.")
		return c.NoContent(http.StatusBadRequest)
	}
	var user storage.User
	err = json.Unmarshal(body, &user)
	if (err != nil) || (user.Login == "") || (user.Password == "") {
		h.log.Warn().Msgf("Register failed. Bad request.")
		return c.NoContent(http.StatusBadRequest)
	}
	err = h.userStorage.AddUser(c.Request().Context(), user)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	var hash string
	hash, err = h.userStorage.GetUserHash(c.Request().Context(), user)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusUnauthorized)
			return c.NoContent(http.StatusUnauthorized)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	cookie := new(http.Cookie)
	cookie.Name = config.Cookie
	cookie.Value = hash
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) ResponseBadRequest(c echo.Context) error {
	h.log.Info().Msgf("bad request. Error code: %v", http.StatusBadRequest)
	return c.NoContent(http.StatusBadRequest)
}
