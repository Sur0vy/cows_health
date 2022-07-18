package monitoringdata

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/internal/usecase/dataprocessor"
	"github.com/Sur0vy/cows_health.git/logger"
)

type Handle interface {
	Add(c echo.Context) error
}

type Handler struct {
	log       *logger.Logger
	mdStorage storages.MonitoringDataStorage
	processor dataprocessor.DataProcessor
}

func NewMonitoringDataHandler(ms storages.MonitoringDataStorage, dp dataprocessor.DataProcessor, log *logger.Logger) Handle {
	return &Handler{
		log:       log,
		mdStorage: ms,
		processor: dp,
	}
}

func (h *Handler) Add(c echo.Context) error {
	var data []models.MonitoringData
	if err := c.Bind(&data); err != nil {
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
