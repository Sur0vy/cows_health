package dataprocessor

import (
	"context"
	"github.com/Sur0vy/cows_health.git/logger"
	"math"
	"sync"
	"time"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type Processor struct {
	log        *logger.Logger
	cowStorage models.CowStorage
	mdStorage  models.MonitoringDataStorage
}

func NewProcessor(ms models.MonitoringDataStorage, cs models.CowStorage, log *logger.Logger) *Processor {
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
	mds, err := dp.GetHealthData(c, data.CowID)
	if err != nil {
		return
	}

	//запустим алгоритмы определения здоровья
	health := dp.CalculateHealth(mds)

	//сохраним данные
	dp.saveHealthData(c, health)
}

func (dp *Processor) GetHealthData(c context.Context, cowID int) (models.DataSlice, error) {
	var (
		avgPH       float64
		avgTemp     float64
		avgMovement float64
	)
	var res models.DataSlice
	mds, err := dp.mdStorage.Get(c, cowID, 10)
	if err != nil {
		dp.log.Info().Msgf("monitoring data get success for %v", cowID)
		return res, err
	} else {
		dp.log.Warn().Err(err).Msgf("getting monitoring data error %v", cowID)
		for i, md := range mds {
			res.Data = append(res.Data, md)
			avgPH += md.PH
			avgTemp += md.Temperature
			avgMovement += md.Movement
			dp.log.Info().Msgf("%d = %v", i, md)
		}
		res.AvgPH = avgPH / float64(len(mds))
		res.AvgTemp = avgTemp / float64(len(mds))
		res.AvgMovement = avgMovement / float64(len(mds))
	}
	return res, nil
}

func (dp *Processor) CalculateHealth(mds models.DataSlice) models.Health {
	//запустим алгоримт определения здоровья и половой активности
	health := models.Health{
		Estrus:    false,
		Ill:       "",
		UpdatedAt: time.Now(),
	}
	for _, md := range mds.Data {
		if (math.Abs(md.Movement-mds.AvgMovement)/md.Movement > 0.1) &&
			(math.Abs(md.Temperature-mds.AvgTemp)/md.Temperature > 0.1) {
			health.Estrus = true
			break
		}
	}
	for _, md := range mds.Data {
		if (math.Abs(mds.AvgMovement-md.Movement)/md.Movement > 0.2) &&
			(math.Abs(md.Temperature-mds.AvgTemp)/md.Temperature > 0.1) &&
			(mds.AvgPH > 6) {
			health.Ill = "Инфекционное заболевание"
			break
		}
	}
	flag := false
	if health.Ill == "" {
		for _, md := range mds.Data {
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
