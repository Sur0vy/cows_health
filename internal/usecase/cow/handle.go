package cow

import (
	"encoding/json"
	go_err "errors"
	"github.com/Sur0vy/cows_health.git/errors"
	"github.com/Sur0vy/cows_health.git/logger"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type Handle interface {
	Add(c echo.Context) error
	Get(c echo.Context) error
	Delete(c echo.Context) error
	GetBreeds(c echo.Context) error
	GetInfo(c echo.Context) error
}

type Handler struct {
	log *logger.Logger
	cs  models.CowStorage
}

func NewCowHandler(cs models.CowStorage, log *logger.Logger) Handle {
	return &Handler{
		log: log,
		cs:  cs,
	}
}

func (h *Handler) Add(c echo.Context) error {
	defer c.Request().Body.Close()
	input, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var cow models.Cow
	if err := json.Unmarshal(input, &cow); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	if cow.AddedAt.IsZero() {
		cow.AddedAt = time.Now()
	}
	//добавим в таблицу коров
	err = h.cs.Add(c.Request().Context(), cow)
	if err != nil {
		if go_err.Is(err, errors.ErrExist) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Cow for user added success")
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Get(c echo.Context) error {
	farmIDStr := c.Param("id")
	h.log.Info().Msgf("farm ID: %s", farmIDStr)
	farmID, err := strconv.Atoi(farmIDStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var cows []models.Cow
	cows, err = h.cs.Get(c.Request().Context(), farmID)

	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Cows for user getting success")
	return c.JSON(http.StatusOK, cows)
}

func (h *Handler) Delete(c echo.Context) error {
	defer c.Request().Body.Close()
	IDs, err := getIDFromJSON(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	err = h.cs.Delete(c.Request().Context(), IDs)
	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) GetInfo(c echo.Context) error {
	cowIDStr := c.Param("id")
	h.log.Info().Msgf("cow ID: %s", cowIDStr)
	cowID, err := strconv.Atoi(cowIDStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var cow models.CowInfo
	cow, err = h.cs.GetInfo(c.Request().Context(), cowID)

	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Cow info for user getting success")
	return c.JSON(http.StatusOK, cow)
}

func (h *Handler) GetBreeds(c echo.Context) error {
	breeds, err := h.cs.GetBreeds(c.Request().Context())
	if err != nil {
		if go_err.Is(err, errors.ErrEmpty) {
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		} else {
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Breeds getting success")
	return c.JSON(http.StatusOK, breeds)
}

func getIDFromJSON(reader io.ReadCloser) ([]int, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var IDs []int
	if err := json.Unmarshal(input, &IDs); err != nil {
		return nil, err
	}
	return IDs, nil
}
