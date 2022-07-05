package dataprocessor

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
)

type dataSlice struct {
	data        []models.MonitoringData
	avgPH       float64
	avgTemp     float64
	avgMovement float64
}

type Processor struct {
	log        *logger.Logger
	cowStorage models.CowStorage
	mdStorage  models.MonitoringDataStorage
}

func NewDataProcessor(ms models.MonitoringDataStorage, cs models.CowStorage, log *logger.Logger) *Processor {
	return &Processor{
		log:        log,
		cowStorage: cs,
		mdStorage:  ms,
	}
}

func (dp *Processor) Run(c context.Context, data models.MonitoringData, wg *sync.WaitGroup) {
	defer wg.Done()
	//запросим, есть ли болюс
	cowID := dp.cowStorage.HasBolus(c, data.BolusNum)
	if cowID != -1 {
		data.CowID = cowID
		dp.log.Info().Msgf("cow id = %v", data.CowID)
	} else {
		dp.log.Warn().Msgf("no cow data for bolus %v", data.BolusNum)
		return
	}
	data.AddedAt = time.Now()
	//сохраним данные
	if err := dp.Save(c, data); err != nil {
		return
	}
	//запросим данные за последние 10 минут
	mds, err := dp.getHealthData(c, data.CowID)
	if err != nil {
		return
	}

	//запустим алгоритмы определения здоровья
	health := dp.calculateHealth(mds)

	//сохраним данные
	dp.saveHealthData(c, health)
}

func (dp *Processor) getHealthData(c context.Context, cowID int) (dataSlice, error) {
	var (
		avgPH       float64
		avgTemp     float64
		avgMovement float64
	)
	var res dataSlice
	mds, err := dp.mdStorage.Get(c, cowID, 10)
	if err != nil {
		dp.log.Info().Msgf("monitoring data get success for %v", cowID)
		return res, err
	} else {
		dp.log.Warn().Err(err).Msgf("getting monitoring data error %v", cowID)
		for i, md := range mds {
			res.data = append(res.data, md)
			avgPH += md.PH
			avgTemp += md.Temperature
			avgMovement += md.Movement
			dp.log.Info().Msgf("%d = %v", i, md)
		}
		res.avgPH = avgPH / float64(len(mds))
		res.avgTemp = avgTemp / float64(len(mds))
		res.avgMovement = avgMovement / float64(len(mds))
	}
	return res, nil
}

func (dp *Processor) calculateHealth(mds dataSlice) models.Health {
	//запустим алгоримт определения здоровья и половой активности
	health := models.Health{
		Estrus:    false,
		Ill:       "",
		UpdatedAt: time.Now(),
	}
	for _, md := range mds.data {
		if (math.Abs(md.Movement-mds.avgMovement)/md.Movement > 0.1) &&
			(math.Abs(md.Temperature-mds.avgTemp)/md.Temperature > 0.1) {
			health.Estrus = true
			break
		}
	}
	for _, md := range mds.data {
		if (math.Abs(mds.avgMovement-md.Movement)/md.Movement > 0.2) &&
			(math.Abs(md.Temperature-mds.avgTemp)/md.Temperature > 0.1) &&
			(mds.avgPH > 6) {
			health.Ill = "Инфекционное заболевание"
			break
		}
	}
	flag := false
	if health.Ill == "" {
		for _, md := range mds.data {
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
	return health
}

func (dp *Processor) Save(c context.Context, data models.MonitoringData) error {
	var err error
	if err = dp.mdStorage.Add(c, data); err == nil {
		dp.log.Info().Msgf("monitoring data added success for %v", data.CowID)
	} else {
		dp.log.Warn().Err(err).Msgf("adding monitoring data error %v", data.CowID)
	}
	return err
}

func (dp *Processor) saveHealthData(c context.Context, health models.Health) {
	if err := dp.cowStorage.UpdateHealth(c, health); err == nil {
		dp.log.Info().Msgf("health data added success for %v", health.CowID)
	} else {
		dp.log.Warn().Err(err).Msgf("adding health data error %v", health.CowID)
	}
}
