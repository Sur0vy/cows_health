package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

type Handle interface {
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Register(c echo.Context) error

	GetFarms(c echo.Context) error
	AddFarm(c echo.Context) error
	DelFarm(c echo.Context) error

	GetCowBreeds(c echo.Context) error
	GetCows(c echo.Context) error

	AddCow(c echo.Context) error
	DelCows(c echo.Context) error
	GetCowInfo(c echo.Context) error

	GetBolusesTypes(c echo.Context) error
	AddMonitoringData(c echo.Context) error

	ResponseBadRequest(c echo.Context) error
}

type BaseHandler struct {
	log         *logger.Logger
	userStorage storage.UserStorage
	farmStorage storage.FarmStorage
}

func NewBaseHandler(us storage.UserStorage, ds storage.FarmStorage, log *logger.Logger) Handle {
	return &BaseHandler{
		log:         log,
		userStorage: us,
		farmStorage: ds,
	}
}

func (h *BaseHandler) Login(c echo.Context) error {
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
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	h.log.Info().Msgf("login success (cookie: %v)", cookie)
	return c.NoContent(http.StatusOK)
}

func (h *BaseHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = config.Cookie
	cookie.Value = ""
	cookie.Expires = time.Time{}
	c.SetCookie(cookie)
	h.log.Info().Msg("logout success")
	return c.NoContent(http.StatusOK)
}

func (h *BaseHandler) Register(c echo.Context) error {
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
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func (h *BaseHandler) GetFarms(c echo.Context) error {
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get farms for user: %v", cookie)
	u := h.userStorage.GetUser(c.Request().Context(), cookie.Value)
	farms, err := h.farmStorage.GetFarms(c.Request().Context(), u.ID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Farms for user getting success")
	return c.JSON(http.StatusOK, farms)
}

func (h *BaseHandler) AddFarm(c echo.Context) error {
	defer c.Request().Body.Close()
	input, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var farm storage.Farm
	if err := json.Unmarshal(input, &farm); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	cookie, _ := c.Cookie(config.Cookie)
	u := h.userStorage.GetUser(c.Request().Context(), cookie.Value)

	farm.UserID = u.ID
	err = h.farmStorage.AddFarm(c.Request().Context(), farm)

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
	h.log.Info().Msg("Farms for user added success")
	return c.NoContent(http.StatusCreated)
}

func (h *BaseHandler) DelFarm(c echo.Context) error {
	farmID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	h.log.Info().Msgf("Delete farm with index: %v", farmID)
	err = h.farmStorage.DelFarm(c.Request().Context(), farmID)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *BaseHandler) GetCows(c echo.Context) error {
	farmIDStr := c.Param("id")
	h.log.Info().Msgf("farm ID: %s", farmIDStr)
	farmID, err := strconv.Atoi(farmIDStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var cows string
	cows, err = h.farmStorage.GetCows(c.Request().Context(), farmID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Cows for user getting success")
	return c.JSON(http.StatusOK, cows)
}

func (h *BaseHandler) GetCowInfo(c echo.Context) error {
	cowIDStr := c.Param("id")
	h.log.Info().Msgf("cow ID: %s", cowIDStr)
	cowID, err := strconv.Atoi(cowIDStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var cow string
	cow, err = h.farmStorage.GetCowInfo(c.Request().Context(), cowID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Cow info for user getting success")
	return c.JSON(http.StatusOK, cow)
}

func (h *BaseHandler) AddMonitoringData(c echo.Context) error {
	defer c.Request().Body.Close()
	input, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var data []storage.MonitoringData
	if err := json.Unmarshal(input, &data); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	h.log.Info().Msg("Monitoring data will be added")
	var wg sync.WaitGroup
	for _, md := range data {
		wg.Add(1)
		dp := storage.NewDataProcessor(h.farmStorage, h.log)
		dp.Run(c.Request().Context(), md, &wg)
	}
	wg.Wait()
	h.log.Info().Msg("All monitoring data was added")
	return c.NoContent(http.StatusAccepted)
}

func (h *BaseHandler) GetBolusesTypes(c echo.Context) error {
	types, err := h.farmStorage.GetBolusesTypes(c.Request().Context())
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("boluses types getting success")
	return c.String(http.StatusOK, types)
}

func (h *BaseHandler) GetCowBreeds(c echo.Context) error {
	breeds, err := h.farmStorage.GetCowBreeds(c.Request().Context())
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			return c.NoContent(http.StatusNoContent)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	h.log.Info().Msg("Breeds getting success")
	return c.JSON(http.StatusOK, breeds)
}

func (h *BaseHandler) AddCow(c echo.Context) error {
	defer c.Request().Body.Close()
	input, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var cow storage.Cow
	if err := json.Unmarshal(input, &cow); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	cow.AddedAt = time.Now()
	//добавим в таблицу коров
	err = h.farmStorage.AddCow(c.Request().Context(), cow)
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
	h.log.Info().Msg("Cow for user added success")
	return c.NoContent(http.StatusCreated)
}

func (h *BaseHandler) DelCows(c echo.Context) error {
	defer c.Request().Body.Close()
	IDs, err := getIDFromJSON(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}
	err = h.farmStorage.DeleteCows(c.Request().Context(), IDs)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return c.NoContent(http.StatusConflict)
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *BaseHandler) ResponseBadRequest(c echo.Context) error {
	h.log.Info().Msgf("bad request. Error code: %v", http.StatusBadRequest)
	return c.NoContent(http.StatusBadRequest)
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
