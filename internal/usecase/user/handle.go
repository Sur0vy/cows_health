package user

import (
	"encoding/json"
	go_err "errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/errors"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
)

type Handle interface {
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Register(c echo.Context) error
}

type Handler struct {
	log *logger.Logger
	us  models.UserStorage
}

func NewUserHandler(us models.UserStorage, log *logger.Logger) Handle {
	return &Handler{
		log: log,
		us:  us,
	}
}

func (h *Handler) Login(c echo.Context) error {
	h.log.Info().Msgf("Handler IN: %v", c)
	defer h.log.Info().Msgf("Handler OUT: %v", c)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	defer c.Request().Body.Close()
	h.log.Info().Msg(string(body))
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var hash string
	hash, err = h.us.GetHash(c.Request().Context(), user)
	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusUnauthorized)
			return c.NoContent(http.StatusUnauthorized)
		} else {
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

func (h *Handler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = config.Cookie
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Time{}
	c.SetCookie(cookie)
	h.log.Info().Msg("logout success")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Register(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Register failed. Bad request.")
		return c.NoContent(http.StatusBadRequest)
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if (err != nil) || (user.Login == "") || (user.Password == "") {
		h.log.Warn().Msgf("Register failed. Bad request.")
		return c.NoContent(http.StatusBadRequest)
	}
	err = h.us.Add(c.Request().Context(), user)
	if err != nil {
		if go_err.Is(err, errors.ErrExist) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	var hash string
	hash, err = h.us.GetHash(c.Request().Context(), user)
	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusUnauthorized)
			return c.NoContent(http.StatusUnauthorized)
		} else {
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
