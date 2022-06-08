package storage

const (
	//tables names
	TBreed          string = "breeds"
	TCow            string = "cows"
	TUser           string = "users"
	TFarm           string = "farms"
	THealth         string = "health"
	TMonitoringData string = "monitoring_data"

	//fields
	FBreed   string = "breed"
	FBreedID string = "breed_id"

	FUserID   string = "user_id"
	FLogin    string = "login"
	FPassword string = "password"

	FFarmID  string = "farm_id"
	FAddress string = "farm_address"
	FName    string = "name"

	FCowID      string = "cow_id"
	FBolus      string = "bolus_sn"
	FDateOfBorn string = "date_of_born"
	FBolusType  string = "bolus_type"

	FDrink     string = "drink"
	FStress    string = "stress"
	FIll       string = "ill"
	FUpdatedAt string = "updated_at"

	FMDID        string = "md_id"
	FAddedAt     string = "added_at"
	FPH          string = "ph"
	FTemperature string = "temperature"
	FMovement    string = "movement"
	FCharge      string = "charge"
)
