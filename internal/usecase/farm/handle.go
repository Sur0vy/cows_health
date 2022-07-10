package farm

import (
	go_err "errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/config"
	"github.com/Sur0vy/cows_health.git/errors"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/logger"
)

type Handle interface {
	Get(c echo.Context) error
	Add(c echo.Context) error
	Delete(c echo.Context) error
}

type Handler struct {
	log         *logger.Logger
	farmStorage storages.FarmStorage
}

func NewFarmHandler(fs storages.FarmStorage, log *logger.Logger) Handle {
	return &Handler{
		log:         log,
		farmStorage: fs,
	}
}

func (h *Handler) Get(c echo.Context) error {
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get farms for user: %v", cookie)
	uID := c.Get("UserID")
	if uID == nil {
		h.log.Warn().Msg("Error reading user from storage")
		return c.NoContent(http.StatusInternalServerError)
	}
	farms, err := h.farmStorage.Get(c.Request().Context(), uID.(int))

	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Farms for user getting success")
	return c.JSON(http.StatusOK, farms)
}

func (h *Handler) Add(c echo.Context) error {
	farm := new(models.Farm)
	if err := c.Bind(farm); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	uID := c.Get("UserID")
	if uID == nil {
		h.log.Warn().Msg("Error reading user from storage")
		return c.NoContent(http.StatusInternalServerError)
	}
	farm.UserID = uID.(int)
	err := h.farmStorage.Add(c.Request().Context(), *farm)

	if err != nil {
		if go_err.Is(err, errors.ErrExist) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Farms for user added success")
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Delete(c echo.Context) error {
	farmID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	h.log.Info().Msgf("Delete farm with index: %v", farmID)
	err = h.farmStorage.Delete(c.Request().Context(), farmID)
	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	return c.NoContent(http.StatusOK)
}
