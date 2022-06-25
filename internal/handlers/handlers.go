package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

type Handle interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)

	GetFarms(c *gin.Context)
	AddFarm(c *gin.Context)
	DelFarm(c *gin.Context)

	GetCowBreeds(c *gin.Context)
	GetCows(c *gin.Context)

	AddCow(c *gin.Context)
	DelCows(c *gin.Context)
	GetCowInfo(c *gin.Context)

	GetBolusesTypes(c *gin.Context)
	AddMonitoringData(c *gin.Context)

	ResponseBadRequest(c *gin.Context)
}

type BaseHandler struct {
	log     *logger.Logger
	storage storage.Storage
}

func NewBaseHandler(s storage.Storage, log *logger.Logger) Handle {
	return &BaseHandler{
		log:     log,
		storage: s,
	}
}

func (h *BaseHandler) Login(c *gin.Context) {
	h.log.Info().Msgf("Handler IN: %v", c)
	defer h.log.Info().Msgf("Handler OUT: %v", c)
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	h.log.Info().Msg(string(body))
	var user storage.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var cookie string
	cookie, err = h.storage.GetUserHash(c, user)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusUnauthorized)
			h.log.Warn().Msgf("Error with code: %v", http.StatusUnauthorized)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return
		}
	}

	c.SetCookie(config.Cookie, cookie, 36000, "/", "", false, false)
	c.Writer.WriteHeader(http.StatusOK)
	h.log.Info().Msgf("login success (cookie: %v)", cookie)
}

func (h *BaseHandler) Logout(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.SetCookie(config.Cookie, "", 0, "/", "", false, false)
	c.Writer.WriteHeader(http.StatusOK)
	h.log.Info().Msg("logout success")
}

func (h *BaseHandler) Register(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Warn().Msgf("Register failed. Bad request.")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var user storage.User
	err = json.Unmarshal(body, &user)
	if (err != nil) || (user.Login == "") || (user.Password == "") {
		h.log.Warn().Msgf("Register failed. Bad request.")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var cookie string
	cookie, err = h.storage.AddUser(c, user)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			c.Writer.WriteHeader(http.StatusConflict)
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return
		}
	}
	c.SetCookie(config.Cookie, cookie, 3600, "/", "", false, true)
	c.Writer.WriteHeader(http.StatusOK)
}

func (h *BaseHandler) GetFarms(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get farms for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}
	farms, err := h.storage.GetFarms(c, u.ID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			c.AbortWithStatus(http.StatusNoContent)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("Farms for user getting success")
	c.String(http.StatusOK, farms)
}

func (h *BaseHandler) AddFarm(c *gin.Context) {
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Add farm for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}
	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var farm storage.Farm
	if err := json.Unmarshal(input, &farm); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	farm.UserID = u.ID
	err = h.storage.AddFarm(c, farm)

	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			c.AbortWithStatus(http.StatusConflict)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("Farms for user added success")
	c.Status(http.StatusCreated)
}

func (h *BaseHandler) DelFarm(c *gin.Context) {
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Delete farm for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}
	farmID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	h.log.Info().Msgf("Delete farm with index: %v", farmID)
	err = h.storage.DelFarm(c, u.ID, farmID)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			c.AbortWithStatus(http.StatusConflict)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.String(http.StatusOK, "")
}

func (h *BaseHandler) GetCows(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get cows for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	farmIDStr := c.Param("id")
	h.log.Info().Msgf("farm ID: %s", farmIDStr)
	farmID, err := strconv.Atoi(farmIDStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var cows string
	cows, err = h.storage.GetCows(c, farmID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			c.AbortWithStatus(http.StatusNoContent)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("Cows for user getting success")
	c.String(http.StatusOK, cows)
}

func (h *BaseHandler) GetCowInfo(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get cows info user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	cowIDStr := c.Param("id")
	h.log.Info().Msgf("cow ID: %s", cowIDStr)
	cowID, err := strconv.Atoi(cowIDStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var cow string
	cow, err = h.storage.GetCowInfo(c, cowID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			c.AbortWithStatus(http.StatusNoContent)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("Cow info for user getting success")
	c.String(http.StatusOK, cow)
}

func (h *BaseHandler) AddMonitoringData(c *gin.Context) {
	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var data []storage.MonitoringData
	if err := json.Unmarshal(input, &data); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.log.Info().Msg("Monitoring data will be added")
	var wg sync.WaitGroup
	for _, md := range data {
		wg.Add(1)
		storage.ProcessMonitoringData(c, &wg, h.storage, md, h.log)
	}
	wg.Wait()
	h.log.Info().Msg("All monitoring data was added")
	c.Status(http.StatusAccepted)
}

func (h *BaseHandler) GetBolusesTypes(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get bolus types ", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	types, err := h.storage.GetBolusesTypes(c)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			c.AbortWithStatus(http.StatusNoContent)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("boluses types getting success")
	c.String(http.StatusOK, types)
}

func (h *BaseHandler) GetCowBreeds(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Get cows breeds", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}
	breeds, err := h.storage.GetCowBreeds(c)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			c.AbortWithStatus(http.StatusNoContent)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("Breeds getting success")
	c.String(http.StatusOK, breeds)
}

func (h *BaseHandler) AddCow(c *gin.Context) {
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Add cow for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}
	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var cow storage.Cow
	if err := json.Unmarshal(input, &cow); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cow.AddedAt = time.Now()
	//добавим в таблицу коров
	err = h.storage.AddCow(c, cow)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			c.AbortWithStatus(http.StatusConflict)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	h.log.Info().Msg("Cow for user added success")
	c.Status(http.StatusCreated)
}

func (h *BaseHandler) DelCows(c *gin.Context) {
	cookie, _ := c.Cookie(config.Cookie)
	h.log.Info().Msgf("Delete cows for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	if u == nil {
		h.log.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	var IDs, err = getIDFromJSON(c.Request.Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.storage.DeleteCows(c, IDs)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
			c.AbortWithStatus(http.StatusConflict)
			return
		default:
			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	c.String(http.StatusAccepted, "")
}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
	h.log.Info().Msgf("bad request. Error code: %v", http.StatusBadRequest)
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
