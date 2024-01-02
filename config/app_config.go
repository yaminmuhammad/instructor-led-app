package config

const (
	APIGroup = "/api/v1"

	AdminGroup              = "/admin"
	TrainerGroup            = "/trainers"
	ParticipantsGroup       = "/participants"
	QuestionGetList         = "/questions"
	QuestionGetById         = "/questions/:id"
	QuestionPost            = "/questions"
	QuestionDelete          = "/questions/:id"
	QuestionTrainer         = "/questions/trainer"
	QuestionGetByTrainerId  = "/questions/trainer/"
	UpdateQuestionByTrainer = "/questions/trainer"

	MasterDataUsers              = "/master-data/users"
	MasterDataUsersCsv           = "/master-data/users/csv"
	MasterDataUserByID           = "/master-data/users/:id"
	MasterDataTrainers           = "/master-data/trainers"
	MasterDataTrainerByID        = "/master-data/trainers/:id"
	MasterDataTrainerByUserID    = "/master-data/trainer/:id"
	MasterDataParticipants       = "/master-data/participants"
	MasterDataParticipantsUpdate = "/master-data/participants/update/:id"
	MasterDataParticipantsRole   = "/master-data/participants/role"
	MasterDataParticipantByID    = "/master-data/participants/:id"

	AbsenceByTrainerScheduleId  = "/absence/trainer/"
	AbsenceParticipantByTrainer = "/absence/trainer/"

	UploadActivityProof = "/upload-activity-proof"

	//schedule
	SchedulePost   = "/schedule/"
	ScheduleList   = "/schedule/"
	ScheduleById   = "/schedule/:id"
	deleteSchedule = "/schedule/:date"

	//fitur participants
	ScheduleByParticipantId = "/partisipant/schedule"
	ParticipantNewQuetion   = "/participant/question"

	AbsencePost      = "/absence/"
	ListAbsences     = "/absence/"
	ListAbsencesById = "/absence/:id"
	DeleteAbsence    = "/absence/:id"

	ScheduleByTrainerId = "/schedule/trainer/"
)
