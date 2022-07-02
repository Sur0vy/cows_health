package monitoringData

//
//import (
//	"context"
//	"github.com/Sur0vy/cows_health.git/internal/models"
//	"github.com/rs/zerolog/log"
//	"math"
//	"sync"
//	"time"
//
//	"github.com/Sur0vy/cows_health.git/internal/logger"
//)
//
//type monitoringData struct {
//	data        []models.MonitoringData
//	avgPH       float64
//	avgTemp     float64
//	avgMovement float64
//}
//
//type DataProcessor struct {
//	log         *logger.Logger
//	FarmStorage farm.FarmStorage
//}
//
//func NewDataProcessor(cs farm.FarmStorage, log *logger.Logger) *DataProcessor {
//	return &DataProcessor{
//		log:         log,
//		FarmStorage: cs,
//	}
//}
//
//func (dp *DataProcessor) Run(c context.Context, data farm.MonitoringData, wg *sync.WaitGroup) {
//	defer wg.Done()
//	//запросим, есть ли болюс
//	cowID := dp.FarmStorage.HasBolus(c, data.BolusNum)
//	if cowID != -1 {
//		data.CowID = cowID
//		log.Info().Msgf("cow id = %v", data.CowID)
//	} else {
//		log.Warn().Msgf("no cow data for bolus %v", data.BolusNum)
//		return
//	}
//	data.AddedAt = time.Now()
//	//сохраним данные
//	if err := dp.saveMonitoringData(c, data); err != nil {
//		return
//	}
//	//запросим данные за последние 10 минут
//	mds, err := dp.getHealthData(c, data.CowID)
//	if err != nil {
//		return
//	}
//
//	//запустим алгоритмы определения здоровья
//	health := dp.calculateHealth(mds)
//
//	//сохраним данные
//	dp.saveHealthData(c, health)
//}
//
//func (dp *DataProcessor) getHealthData(c context.Context, cowID int) (monitoringData, error) {
//	var (
//		avgPH       float64
//		avgTemp     float64
//		avgMovement float64
//	)
//	var monitoringData monitoringData
//	mds, err := dp.FarmStorage.GetMonitoringData(c, cowID, 10)
//	if err != nil {
//		log.Info().Msgf("monitoring data get success for %v", cowID)
//		return monitoringData, err
//	} else {
//		log.Warn().Err(err).Msgf("getting monitoring data error %v", cowID)
//		for i, md := range mds {
//			monitoringData.data = append(monitoringData.data, md)
//			avgPH += md.PH
//			avgTemp += md.Temperature
//			avgMovement += md.Movement
//			log.Info().Msgf("%d = %v", i, md)
//		}
//		monitoringData.avgPH = avgPH / float64(len(mds))
//		monitoringData.avgTemp = avgTemp / float64(len(mds))
//		monitoringData.avgMovement = avgMovement / float64(len(mds))
//	}
//	return monitoringData, nil
//}
//
//func (dp *DataProcessor) calculateHealth(mds monitoringData) farm.Health {
//	//запустим алгоримт определения здоровья и половой активности
//	health := farm.Health{
//		Estrus:    false,
//		Ill:       "",
//		UpdatedAt: time.Now(),
//	}
//	for _, md := range mds.data {
//		if (math.Abs(md.Movement-mds.avgMovement)/md.Movement > 0.1) &&
//			(math.Abs(md.Temperature-mds.avgTemp)/md.Temperature > 0.1) {
//			health.Estrus = true
//			break
//		}
//	}
//	for _, md := range mds.data {
//		if (math.Abs(mds.avgMovement-md.Movement)/md.Movement > 0.2) &&
//			(math.Abs(md.Temperature-mds.avgTemp)/md.Temperature > 0.1) &&
//			(mds.avgPH > 6) {
//			health.Ill = "Инфекционное заболевание"
//			break
//		}
//	}
//	flag := false
//	if health.Ill == "" {
//		for _, md := range mds.data {
//			if (md.Temperature < 40) || (md.Temperature > 41) ||
//				(md.PH > 5.5) {
//				flag = true
//				break
//			}
//		}
//		if !flag {
//			health.Ill = "Нервное заболевание"
//		}
//	}
//	return health
//}
//
//func (dp *DataProcessor) saveHealthData(c context.Context, health farm.Health) error {
//	var err error
//	if err := dp.FarmStorage.UpdateHealth(c, health); err == nil {
//		log.Info().Msgf("health data added success for %v", health.CowID)
//	} else {
//		log.Warn().Err(err).Msgf("adding health data error %v", health.CowID)
//	}
//	return err
//}
//
//func (dp *DataProcessor) saveMonitoringData(c context.Context, data farm.MonitoringData) error {
//	var err error
//	if err = dp.FarmStorage.AddMonitoringData(c, data); err == nil {
//		log.Info().Msgf("monitoring data added success for %v", data.CowID)
//	} else {
//		log.Warn().Err(err).Msgf("adding monitoring data error %v", data.CowID)
//	}
//	return err
//}
