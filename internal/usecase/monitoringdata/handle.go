package monitoringdata

import (
	"encoding/json"
	"github.com/Sur0vy/cows_health.git/internal/usecase/dataprocessor"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
)

type Handle interface {
	Add(c echo.Context) error
}

type Handler struct {
	log       *logger.Logger
	mdStorage models.MonitoringDataStorage
	processor dataprocessor.Processor
}

func NewMonitoringDataHandler(ms models.MonitoringDataStorage, dp dataprocessor.Processor, log *logger.Logger) Handle {
	return &Handler{
		log:       log,
		mdStorage: ms,
		processor: dp,
	}
}

func (h *Handler) Add(c echo.Context) error {
	defer c.Request().Body.Close()
	input, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	var data []models.MonitoringData
	if err := json.Unmarshal(input, &data); err != nil {
		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	h.log.Info().Msg("Monitoring data will be added")
	var wg sync.WaitGroup
	for _, md := range data {
		wg.Add(1)
		h.processor.Run(c.Request().Context(), md, &wg)
	}
	wg.Wait()
	h.log.Info().Msg("All monitoring data was added")
	return c.NoContent(http.StatusAccepted)
}
