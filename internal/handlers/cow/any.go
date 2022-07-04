package cow

//	GetCowBreeds(c echo.Context) error
//	GetCows(c echo.Context) error
//
//	AddCow(c echo.Context) error
//	DelCows(c echo.Context) error
//	GetCowInfo(c echo.Context) error
//
//	GetBolusesTypes(c echo.Context) error
//	AddMonitoringData(c echo.Context) error
//}

//
//func (h *FarmHandler) GetCows(c echo.Context) error {
//	farmIDStr := c.Param("id")
//	h.log.Info().Msgf("farm ID: %s", farmIDStr)
//	farmID, err := strconv.Atoi(farmIDStr)
//	if err != nil {
//		return c.NoContent(http.StatusBadRequest)
//	}
//
//	var cows string
//	cows, err = h.farmStorage.GetCows(c.Request().Context(), farmID)
//
//	if err != nil {
//		switch err.(type) {
//		case *errors.EmptyError:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
//			return c.NoContent(http.StatusNoContent)
//		default:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
//			return c.NoContent(http.StatusInternalServerError)
//		}
//	}
//	h.log.Info().Msg("Cows for user getting success")
//	return c.JSON(http.StatusOK, cows)
//}
//
//func (h *FarmHandler) GetCowInfo(c echo.Context) error {
//	cowIDStr := c.Param("id")
//	h.log.Info().Msgf("cow ID: %s", cowIDStr)
//	cowID, err := strconv.Atoi(cowIDStr)
//	if err != nil {
//		return c.NoContent(http.StatusBadRequest)
//	}
//
//	var cow string
//	cow, err = h.farmStorage.GetCowInfo(c.Request().Context(), cowID)
//
//	if err != nil {
//		switch err.(type) {
//		case *errors.EmptyError:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
//			return c.NoContent(http.StatusNoContent)
//		default:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
//			return c.NoContent(http.StatusInternalServerError)
//		}
//	}
//	h.log.Info().Msg("Cow info for user getting success")
//	return c.JSON(http.StatusOK, cow)
//}
//
//func (h *FarmHandler) AddMonitoringData(c echo.Context) error {
//	defer c.Request().Body.Close()
//	input, err := ioutil.ReadAll(c.Request().Body)
//	if err != nil {
//		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
//		return c.NoContent(http.StatusBadRequest)
//	}
//
//	var data []models.MonitoringData
//	if err := json.Unmarshal(input, &data); err != nil {
//		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
//		return c.NoContent(http.StatusBadRequest)
//	}
//
//	h.log.Info().Msg("Monitoring data will be added")
//	var wg sync.WaitGroup
//	for _, md := range data {
//		wg.Add(1)
//		dp := monitoringData.NewDataProcessor(h.farmStorage, h.log)
//		dp.Run(c.Request().Context(), md, &wg)
//	}
//	wg.Wait()
//	h.log.Info().Msg("All monitoring data was added")
//	return c.NoContent(http.StatusAccepted)
//}
//
//func (h *FarmHandler) GetBolusesTypes(c echo.Context) error {
//	types, err := h.farmStorage.GetBolusesTypes(c.Request().Context())
//	if err != nil {
//		switch err.(type) {
//		case *errors.EmptyError:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
//			return c.NoContent(http.StatusNoContent)
//		default:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
//			return c.NoContent(http.StatusInternalServerError)
//		}
//	}
//	h.log.Info().Msg("boluses types getting success")
//	return c.String(http.StatusOK, types)
//}
//
//func (h *FarmHandler) GetCowBreeds(c echo.Context) error {
//	breeds, err := h.farmStorage.GetCowBreeds(c.Request().Context())
//	if err != nil {
//		switch err.(type) {
//		case *errors.EmptyError:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusNoContent)
//			return c.NoContent(http.StatusNoContent)
//		default:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
//			return c.NoContent(http.StatusInternalServerError)
//		}
//	}
//	h.log.Info().Msg("Breeds getting success")
//	return c.JSON(http.StatusOK, breeds)
//}
//
//func (h *FarmHandler) AddCow(c echo.Context) error {
//	defer c.Request().Body.Close()
//	input, err := ioutil.ReadAll(c.Request().Body)
//	if err != nil {
//		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
//		return c.NoContent(http.StatusBadRequest)
//	}
//
//	var cow models.Cow
//	if err := json.Unmarshal(input, &cow); err != nil {
//		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
//		return c.NoContent(http.StatusBadRequest)
//	}
//
//	cow.AddedAt = time.Now()
//	//добавим в таблицу коров
//	err = h.farmStorage.AddCow(c.Request().Context(), cow)
//	if err != nil {
//		switch err.(type) {
//		case *errors.ExistError:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
//			return c.NoContent(http.StatusConflict)
//		default:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
//			return c.NoContent(http.StatusInternalServerError)
//		}
//	}
//	h.log.Info().Msg("Cow for user added success")
//	return c.NoContent(http.StatusCreated)
//}
//
//func (h *FarmHandler) DelCows(c echo.Context) error {
//	defer c.Request().Body.Close()
//	IDs, err := getIDFromJSON(c.Request().Body)
//	if err != nil {
//		h.log.Warn().Msgf("Error with code: %v", http.StatusBadRequest)
//		return c.NoContent(http.StatusBadRequest)
//	}
//	err = h.farmStorage.DeleteCows(c.Request().Context(), IDs)
//	if err != nil {
//		switch err.(type) {
//		case *errors.EmptyError:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusConflict)
//			return c.NoContent(http.StatusConflict)
//		default:
//			h.log.Warn().Msgf("Error with code: %v", http.StatusInternalServerError)
//			return c.NoContent(http.StatusInternalServerError)
//		}
//	}
//	return c.NoContent(http.StatusAccepted)
//}
//
//func getIDFromJSON(reader io.ReadCloser) ([]int, error) {
//	input, err := ioutil.ReadAll(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	var IDs []int
//	if err := json.Unmarshal(input, &IDs); err != nil {
//		return nil, err
//	}
//	return IDs, nil
//}
