package dto

type ScheduleImagesDTO struct {
	ID         string `json:"id"`
	ScheduleID string `json:"scheduleId"`
	FileName   string `json:"filename"`
}
