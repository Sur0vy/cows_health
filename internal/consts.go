package internal

const (
	//tables names
	TBreed string = "breeds"
	TCow   string = "cows"
	//TUser           string = "users"
	TFarm           string = "farms"
	THealth         string = "health"
	TMonitoringData string = "monitoring_data"

	//fields
	FBreedID string = "breed_id"

	//FUserID   string = "user_id"
	//FLogin    string = "login"
	//FPassword string = "password"
	//FDeleted  string = "deleted"

	FFarmID  string = "farm_id"
	FAddress string = "address"
	FName    string = "name"

	FCowID      string = "cow_id"
	FBolus      string = "bolus_sn"
	FDateOfBorn string = "date_of_born"
	FBolusType  string = "bolus_type"

	FEstrus    string = "estrus"
	FIll       string = "ill"
	FUpdatedAt string = "updated_at"

	FMDID        string = "md_id"
	FAddedAt     string = "added_at"
	FPH          string = "ph"
	FTemperature string = "temperature"
	FMovement    string = "movement"
	FCharge      string = "charge"
)
