package storage

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Farm struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  int    `json:"-"`
}

type CowBreed struct {
	ID    int    `json:"breed_id"`
	Breed string `json:"breed"`
}

type Cow struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name"`
	BolusNum        string `json:"bolusNum"`
	LactationDay    int    `json:"lactationDay"`
	Breed           string `json:"type"`
	BreedID         int    `json:"typeID,omitempty"`
	DateOfBorn      string `json:"dateOfBorn"`
	BolusID         int    `json:"bolusID"`
	Age             int    `json:"age,omitempty"`
	FarmID          int    `json:"farmID"`
	Farm            string `json:"farm,omitempty"`
	Calf            bool   `json:"calf,omitempty"`            //*bool
	InseminationDay int    `json:"inseminationDay,omitempty"` //*int
}

type Bolus struct {
	ID           int    `json:"id"`
	SerialNumber string `json:"number"`
	Type         string `json:"type"`
	TypeID       int    `json:"typeID,omitempty"`
	CowName      string `json:"cowName"`
	Status       string `json:"status"`
	Charge       int    `json:"chargeLevel"`
	CowID        int    `json:"cowID"`
	FarmID       int    `json:"-"`
}

type Health struct {
	Drink       bool    `json:"drink"`
	Stress      bool    `json:"stress"`
	Temperature float32 `json:"temperature"`
	Activity    float32 `json:"activity"`
	CowID       int     `json:"-"`
}

type MonitoringData struct {
	BolusID      int     `json:"-"`
	SerialNumber string  `json:"num"`
	Time         string  `json:"dateTime"`
	PH           float32 `json:"ph"`
	Temperature  float32 `json:"temperature"`
	Movement     float32 `json:"movement"`
	Humidity     float32 `json:"humidity"`
	Charge       float32 `json:"charge"`
}

type CowInfo struct {
	Health  Health           `json:"health"`
	Summary Cow              `json:"summary"`
	History []MonitoringData `json:"history"`
}

type BolusType struct {
	ID   int
	Name string
}
