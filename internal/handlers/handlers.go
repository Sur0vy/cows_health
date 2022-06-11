package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

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

	//GetFarmInfo(c *gin.Context)
	//GetCows(c *gin.Context)

	//AddCow(c *gin.Context)
	//
	//GetBolusesTypes(c *gin.Context)
	//AddBolusData(c *gin.Context)
	//
	//GetCowInfo(c *gin.Context)
	//GetCowBreeds(c *gin.Context)

	ResponseBadRequest(c *gin.Context)
}

type BaseHandler struct {
	storage storage.Storage
}

func NewBaseHandler(s *storage.Storage) Handle {
	return &BaseHandler{
		storage: *s,
	}
}

func (h *BaseHandler) Login(c *gin.Context) {
	logger.Wr.Info().Msgf("Handler IN: %v", c)
	defer logger.Wr.Info().Msgf("Handler OUT: %v", c)
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Wr.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	logger.Wr.Info().Msg(string(body))
	var user storage.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Wr.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var cookie string
	cookie, err = h.storage.GetUserHash(c, user)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusUnauthorized)
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusUnauthorized)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return
		}
	}

	c.SetCookie(config.Cookie, cookie, 36000, "/", "", false, false)
	c.Writer.WriteHeader(http.StatusOK)
	logger.Wr.Info().Msgf("login success (cookie: %v)", cookie)
}

func (h *BaseHandler) Logout(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.SetCookie(config.Cookie, "", 0, "/", "", false, false)
	c.Writer.WriteHeader(http.StatusOK)
	logger.Wr.Info().Msg("logout success")
}

func (h *BaseHandler) Register(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Wr.Warn().Msgf("Register failed. Bad request.")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var user storage.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Wr.Warn().Msgf("Register failed. Bad request.")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var cookie string
	cookie, err = h.storage.AddUser(c, user)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			c.Writer.WriteHeader(http.StatusConflict)
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusConflict)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			return
		}
	}
	c.SetCookie(config.Cookie, cookie, 3600, "/", "", false, true)
	c.Writer.WriteHeader(http.StatusOK)
}

func (h *BaseHandler) GetFarms(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	logger.Wr.Info().Msgf("Get farms for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)
	farms, err := h.storage.GetFarms(c, u.ID)

	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusNoContent)
			c.AbortWithStatus(http.StatusNoContent)
		default:
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	logger.Wr.Info().Msg("Farms for user getting success")
	c.String(http.StatusOK, farms)
}

func (h *BaseHandler) AddFarm(c *gin.Context) {
	cookie, _ := c.Cookie(config.Cookie)
	logger.Wr.Info().Msgf("Add farm for user: %v", cookie)
	u := h.storage.GetUser(c, cookie)

	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Wr.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var farm storage.Farm
	if err := json.Unmarshal(input, &farm); err != nil {
		logger.Wr.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	farm.UserID = u.ID
	err = h.storage.AddFarm(c, farm)

	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusConflict)
			c.AbortWithStatus(http.StatusConflict)
		default:
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	logger.Wr.Info().Msg("Farms for user added success")
	c.Status(http.StatusCreated)
}

func (h *BaseHandler) DelFarm(c *gin.Context) {
	farmID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Wr.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	logger.Wr.Info().Msgf("Delete farm with index: %v", farmID)
	err = h.storage.DelFarm(c, farmID)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusConflict)
			c.AbortWithStatus(http.StatusConflict)
		default:
			logger.Wr.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	c.String(http.StatusOK, "")
}

//func (h *BaseHandler) GetFarmInfo(c *gin.Context) {
//	c.Writer.Header().Set("Content-Type", "application/json")
//	farmIdStr := c.Param("id")
//	fmt.Printf("farm id = %s\n", farmIdStr)
//	farmID, err := strconv.Atoi(farmIdStr)
//	if err != nil {
//		c.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	var farmInfo string
//	farmInfo, err = h.storage.GetFarmInfo(c, farmID)
//	if err != nil {
//		switch err.(type) {
//		case *storage.EmptyError:
//			c.Writer.WriteHeader(http.StatusNoContent)
//			return
//		default:
//			c.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	} else {
//		c.String(http.StatusOK, farmInfo)
//	}
//}
//
//func (h *BaseHandler) GetCows(c *gin.Context) {
//	c.Writer.Header().Set("Content-Type", "application/json")
//	farmIdStr := c.Param("id")
//	fmt.Printf("farm id = %s\n", farmIdStr)
//	farmID, err := strconv.Atoi(farmIdStr)
//	if err != nil {
//		c.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	var cows string
//	cows, err = h.storage.GetCows(c, farmID)
//	if err != nil {
//		switch err.(type) {
//		case *storage.EmptyError:
//			c.Writer.WriteHeader(http.StatusNoContent)
//			return
//		default:
//			c.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	} else {
//		c.String(http.StatusOK, cows)
//	}
//}
//
//func (h *BaseHandler) GetCowInfo(c *gin.Context) {
//	c.Writer.Header().Set("Content-Type", "application/json")
//	cowIdStr := c.Param("id")
//	fmt.Printf("cow id = %s\n", cowIdStr)
//	cowID, err := strconv.Atoi(cowIdStr)
//	if err != nil {
//		c.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	var cowInfo string
//	cowInfo, err = h.storage.GetCowInfo(c, cowID)
//	if err != nil {
//		switch err.(type) {
//		case *storage.EmptyError:
//			c.Writer.WriteHeader(http.StatusNoContent)
//			return
//		default:
//			c.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	} else {
//		c.String(http.StatusOK, cowInfo)
//	}
//}
//
//func (h *BaseHandler) DelCows(c *gin.Context) {
//	c.Writer.Header().Set("Content-Type", "application/json")
//	IDs, err := getIDFromJSON(c.Request.Body)
//
//	if err != nil {
//		c.String(http.StatusBadRequest, "")
//		return
//	}
//
//	h.storage.DeleteCows(c, IDs)
//	c.String(http.StatusAccepted, "")
//}
//
//func (h *BaseHandler) GetBolusesTypes(c *gin.Context) {
//	c.Writer.Header().Set("Content-Type", "application/json")
//	types, err := h.storage.GetBolusesTypes(c)
//	if err != nil {
//		switch err.(type) {
//		case *storage.EmptyError:
//			c.Writer.WriteHeader(http.StatusNoContent)
//			return
//		default:
//			c.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	} else {
//		c.String(http.StatusOK, types)
//	}
//}
//
//func (h *BaseHandler) GetCowBreeds(c *gin.Context) {
//	c.Writer.Header().Set("Content-Type", "application/json")
//	types, err := h.storage.GetCowBreeds(c)
//	if err != nil {
//		switch err.(type) {
//		case *storage.EmptyError:
//			c.Writer.WriteHeader(http.StatusNoContent)
//			return
//		default:
//			c.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	} else {
//		c.String(http.StatusOK, types)
//	}
//}
//
//func (h *BaseHandler) AddCow(c *gin.Context) {
//	input, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		c.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	fmt.Println(string(input))
//	var cow storage.Cow
//	if err := json.Unmarshal(input, &cow); err != nil {
//		c.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = h.storage.AddCow(c, cow)
//	if err != nil {
//		c.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	c.Status(http.StatusCreated)
//
//}
//
//func (h *BaseHandler) AddBolusData(c *gin.Context) {
//	input, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		c.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	fmt.Println(string(input))
//	var monitoringData storage.MonitoringData
//	if err := json.Unmarshal(input, &monitoringData); err != nil {
//		c.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = h.storage.AddMonitoringData(c, monitoringData)
//	if err != nil {
//		switch err.(type) {
//		case *storage.EmptyError:
//			c.Writer.WriteHeader(http.StatusConflict)
//			return
//		default:
//			c.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	}
//	c.Status(http.StatusCreated)
//}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
	logger.Wr.Info().Msgf("bad request. Error code: %v", http.StatusBadRequest)
}

func getIDFromJSON(reader io.ReadCloser) ([]int, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	fmt.Printf("parse JSON array: body = %s\n", input)

	var IDs []int
	if err := json.Unmarshal(input, &IDs); err != nil {
		fmt.Printf("\tID unmarshal error: %s\n", err)
		return nil, err
	}
	fmt.Printf("\tID unmarshal success")
	return IDs, nil
}
