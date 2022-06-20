package storage

import (
	"context"
	"math"
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
	var (
		avgPH       float64
		avgTemp     float64
		avgMovement float64
	)
	mds, err := s.GetMonitoringData(c, data.CowID, 10)
	if err != nil {
		logger.Wr.Info().Msgf("monitoring data get success for %v", data.CowID)
		return
	} else {
		logger.Wr.Warn().Err(err).Msgf("getting monitoring data error %v", data.CowID)
		for i, md := range mds {
			avgPH += md.PH
			avgTemp += md.Temperature
			avgMovement += md.Movement
			logger.Wr.Info().Msgf("%d = %v", i, md)
		}
		avgPH /= float64(len(mds))
		avgTemp /= float64(len(mds))
		avgMovement /= float64(len(mds))
	}

	//запустим алгоримт определения здоровья и половой активности
	health := Health{
		CowID:     cowID,
		Estrus:    false,
		Ill:       "",
		UpdatedAt: time.Now(),
	}
	for _, md := range mds {
		if (math.Abs(md.Movement-avgMovement)/md.Movement > 0.1) &&
			(math.Abs(md.Temperature-avgTemp)/md.Temperature > 0.1) {
			health.Estrus = true
			break
		}
	}
	for _, md := range mds {
		if (math.Abs(avgMovement-md.Movement)/md.Movement > 0.2) &&
			(math.Abs(md.Temperature-avgTemp)/md.Temperature > 0.1) &&
			(avgPH > 6) {
			health.Ill = "Инфекционное заболевание"
			break
		}
	}
	flag := false
	if health.Ill == "" {
		for _, md := range mds {
			if (md.Temperature < 40) || (md.Temperature > 41) ||
				(md.PH > 5.5) {
				flag = true
				break
			}
		}
		if !flag {
			health.Ill = "Нервное заболевание"
		}
	}

	//запишем в состояние о корове
	if err := s.UpdateHealth(c, health); err == nil {
		logger.Wr.Info().Msgf("health data added success for %v", data.CowID)
	} else {
		logger.Wr.Warn().Err(err).Msgf("adding health data error %v", data.CowID)
		return
	}
}
