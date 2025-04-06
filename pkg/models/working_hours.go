package models

type WorkingHoursRequest struct {
	StartTime string `json:"start_time"`
	Timezone1 string `json:"timezone_1"`
	Timezone2 string `json:"timezone_2"`
	Country1  string `json:"country_1"`
	Country2  string `json:"country_2"`
}

type WorkingHoursResponse struct {
	IsOverlap  bool   `json:"is_overlap"`
	LocalTime1 string `json:"local_time_1"`
	LocalTime2 string `json:"local_time_2"`
	Holiday1   bool   `json:"holiday_1"`
	Holiday2   bool   `json:"holiday_2"`
}
