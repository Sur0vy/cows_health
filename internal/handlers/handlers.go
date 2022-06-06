package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

type Handle interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)

	GetFarms(c *gin.Context)
	AddFarm(c *gin.Context)
	GetFarmInfo(c *gin.Context)
	GetCows(c *gin.Context)
	DelCows(c *gin.Context)

	GetBoluses(c *gin.Context)
	GetBolusesTypes(c *gin.Context)
	AddBolus(c *gin.Context)
	DelBoluses(c *gin.Context)
	AddBolusData(c *gin.Context)

	GetCowInfo(c *gin.Context)
	GetCowBreeds(c *gin.Context)
	AddCow(c *gin.Context)

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
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))
	var user storage.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var cookie string
	cookie, err = h.storage.CheckUser(c, user)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	c.SetCookie(config.Cookie, cookie, 36000, "/", "", false, false)
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Println(c.Cookie(config.Cookie))
}

func (h *BaseHandler) Logout(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.SetCookie(config.Cookie, "", 0, "/", "", false, false)
	c.Writer.WriteHeader(http.StatusOK)
}

func (h *BaseHandler) Register(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var user storage.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var cookie string
	cookie, err = h.storage.AddUser(c, user)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			c.Writer.WriteHeader(http.StatusConflict)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	c.SetCookie(config.Cookie, cookie, 3600, "/", "", false, true)
	c.Writer.WriteHeader(http.StatusOK)
}

func (h *BaseHandler) GetFarms(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cookie, _ := c.Cookie(config.Cookie)
	userID, _ := h.storage.GetUser(c, cookie)
	farms, err := h.storage.GetFarms(c, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusNoContent)
	} else {
		c.String(http.StatusOK, farms)
	}
}

func (h *BaseHandler) AddFarm(c *gin.Context) {
	cookie, _ := c.Cookie(config.Cookie)
	userID, _ := h.storage.GetUser(c, cookie)

	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var farm storage.Farm
	if err := json.Unmarshal(input, &farm); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//проверка наличия такой фермы
	farm.UserID = userID
	err = h.storage.AddFarm(c, farm)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			c.Writer.WriteHeader(http.StatusConflict)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	c.Status(http.StatusCreated)
}

func (h *BaseHandler) GetFarmInfo(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	farmIdStr := c.Param("id")
	fmt.Printf("farm id = %s\n", farmIdStr)
	farmID, err := strconv.Atoi(farmIdStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	var farmInfo string
	farmInfo, err = h.storage.GetFarmInfo(c, farmID)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.String(http.StatusOK, farmInfo)
	}
}

func (h *BaseHandler) GetCows(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	farmIdStr := c.Param("id")
	fmt.Printf("farm id = %s\n", farmIdStr)
	farmID, err := strconv.Atoi(farmIdStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	var cows string
	cows, err = h.storage.GetCows(c, farmID)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.String(http.StatusOK, cows)
	}
}

func (h *BaseHandler) GetCowInfo(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	cowIdStr := c.Param("id")
	fmt.Printf("cow id = %s\n", cowIdStr)
	cowID, err := strconv.Atoi(cowIdStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	var cowInfo string
	cowInfo, err = h.storage.GetCowInfo(c, cowID)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.String(http.StatusOK, cowInfo)
	}
}

func (h *BaseHandler) DelCows(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	IDs, err := getIDFromJSON(c.Request.Body)

	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	h.storage.DeleteCows(c, IDs)
	c.String(http.StatusAccepted, "")
}

func (h *BaseHandler) GetBoluses(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	farmIdStr := c.Param("id")
	fmt.Printf("farm id = %s\n", farmIdStr)
	farmId, err := strconv.Atoi(farmIdStr)
	if err != nil {
		fmt.Printf("\tError: no boluses")
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	var boluses string
	boluses, err = h.storage.GetBoluses(c, farmId)
	if err != nil {
		fmt.Printf("\tError: no cows")
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	} else {
		c.String(http.StatusOK, boluses)
	}
}

func (h *BaseHandler) GetBolusesTypes(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	types, err := h.storage.GetBolusesTypes(c)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.String(http.StatusOK, types)
	}
}

func (h *BaseHandler) GetCowBreeds(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	types, err := h.storage.GetCowBreeds(c)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.String(http.StatusOK, types)
	}
}

func (h *BaseHandler) AddCow(c *gin.Context) {
	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(string(input))
	var cow storage.Cow
	if err := json.Unmarshal(input, &cow); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.AddCow(c, cow)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)

}

func (h *BaseHandler) AddBolus(c *gin.Context) {
	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(string(input))
	var bolus storage.Bolus
	if err := json.Unmarshal(input, &bolus); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.AddBolus(c, bolus)
	if err != nil {
		switch err.(type) {
		case *storage.ExistError:
			c.Writer.WriteHeader(http.StatusConflict)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	c.Status(http.StatusCreated)
}

func (h *BaseHandler) DelBoluses(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	IDs, err := getIDFromJSON(c.Request.Body)

	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	h.storage.DeleteBoluses(c, IDs)
	c.String(http.StatusAccepted, "")
}

func (h *BaseHandler) AddBolusData(c *gin.Context) {
	input, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(string(input))
	var monitoringData storage.MonitoringData
	if err := json.Unmarshal(input, &monitoringData); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.AddMonitoringData(c, monitoringData)
	if err != nil {
		switch err.(type) {
		case *storage.EmptyError:
			c.Writer.WriteHeader(http.StatusConflict)
			return
		default:
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	c.Status(http.StatusCreated)
}

func (h *BaseHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
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
