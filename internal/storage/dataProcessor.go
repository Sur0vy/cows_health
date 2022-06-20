package storage

import (
	"context"
	"sync"
	"time"

	"github.com/Sur0vy/cows_health.git/internal/logger"
)

func ProcessMonitoringData(c context.Context, wg *sync.WaitGroup, s Storage, data MonitoringData) {
	defer wg.Done()
	//запросим, есть ли болюс
	cowID := s.HasBolus(c, data.BolusNum)
	if cowID != -1 {
		data.CowID = cowID
		logger.Wr.Info().Msgf("cow id = %v", data.CowID)
	} else {
		logger.Wr.Warn().Msgf("no cow data for bolus %v", data.BolusNum)
		return
	}
	data.AddedAt = time.Now()
	//запишем данные
	if err := s.AddMonitoringData(c, data); err == nil {
		logger.Wr.Info().Msgf("monitoring data added success for %v", data.CowID)
	} else {
		logger.Wr.Warn().Err(err).Msgf("adding monitoring data error %v", data.CowID)
		return
	}
	//запросим данные за последние 10 минут
	mds, err := s.GetMonitoringData(c, data.CowID, 10)
	if err != nil {
		logger.Wr.Info().Msgf("monitoring data get success for %v", data.CowID)
		return
	} else {
		logger.Wr.Warn().Err(err).Msgf("getting monitoring data error %v", data.CowID)
		for i, md := range mds {
			logger.Wr.Info().Msgf("%d = %v", i, md)
		}
	}

	//запустим алгоримт определения здоровья и половой активности
	var health Health
	health.UpdatedAt = time.Now()

	//запишем в состояние о корове
	if err := s.UpdateHealth(c, health); err == nil {
		logger.Wr.Info().Msgf("health data added success for %v", data.CowID)
	} else {
		logger.Wr.Warn().Err(err).Msgf("adding health data error %v", data.CowID)
		return
	}
}
