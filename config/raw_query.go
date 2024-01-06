package config

const (
	InsertTrainer      = `INSERT INTO trainers (phone_number, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	ListTrainers       = `SELECT id, phone_number, user_id, created_at, updated_at FROM trainers ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	GetTrainerByID     = `SELECT id,phone_number,user_id, created_at, updated_at FROM trainers WHERE id= $1`
	GetTrainerByUserID = `SELECT id,phone_number,user_id, created_at, updated_at FROM trainers WHERE user_id= $1`
	DeleteTrainerByID  = `DELETE FROM trainers WHERE id = $1`
	UpdateTrainerByID  = `UPDATE trainers
	SET phone_number = $1, user_id = $2, updated_at = $3
	WHERE id = $4 RETURNING id,phone_number, user_id`
	SelectUserAll  = "SELECT * FROM users LIMIT $1 OFFSET $2"
	SelectUserByID = "SELECT * FROM users WHERE id = $1"
	// CRUD User
	ListUsers                                  = `Select id,name,email,username,address,role,created_at,updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	GetUserByID                                = `SELECT id,name,email,username,address,role,created_at,updated_at FROM users WHERE id = $1`
	InsertUser                                 = `INSERT INTO users(name,email,username,address,hash_password,role,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at`
	InsertAndGetUserRole                       = `INSERT INTO users (name, email, username, address, hash_password, role, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	InsertUserToParticipant                    = `INSERT INTO participants (user_id, created_at, updated_at) VALUES ($1, $2, $3)`
	InsertUserToTrainer                        = `INSERT INTO trainers (user_id, created_at, updated_at) VALUES ($1, $2, $3)`
	UpdatedUserByName                          = `UPDATE users SET name = $2  WHERE id = $1`
	UpdatedUserByEmail                         = `UPDATE users SET email = $2  WHERE id = $1`
	UpdatedUserByUsername                      = `UPDATE users SET username = $2  WHERE id = $1`
	UpdatedUserByAddress                       = `UPDATE users SET address = $2  WHERE id = $1`
	UpdatedUserByHashPassword                  = `UPDATE users SET hash_password = $2  WHERE id = $1`
	UpdatedUserByRole                          = `UPDATE users SET role = $2  WHERE id = $1`
	GetUserIDbyName                            = `SELECT id FROM users WHERE name = $1`
	UpdatedUserAll                             = `UPDATE users SET name = $2,email = $3,username = $4,address = $5,hash_password=$6,role =$7  WHERE id = $1`
	DeleteUserByID                             = `DELETE FROM users WHERE id = $1`
	SelectTaskList                             = `SELECT id, title, content, author_id, created_at, updated_at FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	SelectQuestionList                         = `SELECT id, question, status, participant_id,trainer_id,schedule_id, created_at, updated_at FROM questions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	SelectQuestionByID                         = `SELECT id, question, status, participant_id,trainer_id,schedule_id, created_at, updated_at FROM questions WHERE id = $1`
	InsertQuestion                             = `INSERT INTO questions ( question, status, participant_id,trainer_id,schedule_id, updated_at) VALUES ($1, $2, $3, $4,$5,$6) RETURNING id, created_at`
	InsertQuestionNew                          = `INSERT INTO questions ( question, answer, status, participant_id, trainer_id, schedule_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	UpdateQuestion                             = `UPDATE questions SET status = $2 WHERE id = $1`
	DeleteQuestion                             = `DELETE FROM questions WHERE id = $1`
	DeleteSchedule                             = `DELETE FROM Schedule WHERE date = $1`
	SelectQuestionByTrainerID                  = `SELECT id, question, status, participant_id,trainer_id,schedule_id, created_at, updated_at FROM questions WHERE trainer_id = $1 limit $2 offset $3`
	SelectQuestionByScheduleIDandParticipantID = `SELECT id, question, status, participant_id FROM questions WHERE schedule_id = $1 AND participant_id = $2 AND status = $3`

	UpdateQuestionStatusByTrainer = `
	UPDATE
		questions
	SET
		answer = $3,
		status = $4,
		updated_at = $5,
		trainer_id = $6
	WHERE
		participant_id = $1 AND schedule_id = $2
	RETURNING
	id, question, answer, status, participant_id, trainer_id, schedule_id, created_at, updated_at`

	InsertParticipant = `
	INSERT INTO
		participants(date_of_birth, place_of_birth, last_education, user_id,role)
	VALUES
		($1, $2, $3, $4,$5)
	RETURNING
		id, date_of_birth, place_of_birth, last_education, user_id,role;
	`
	ListParticipant = `
	SELECT
		id,
		date_of_birth,
		place_of_birth,
		last_education,
		user_id,
		role
	FROM
		participants
	ORDER BY
		created_at DESC
	LIMIT $1 OFFSET $2;
	`

	ListParticipantByIDandRole = `
	SELECT
		id,
		role
	FROM
		participants;`
	GetParticipantByID = `
	SELECT
		id,
		date_of_birth,
		place_of_birth,
		last_education,
		user_id,
		role
	FROM
		participants
	WHERE
		id = $1
	ORDER BY
		created_at DESC;
	`
	UpdateParticipantByID = `
	UPDATE
		participants
	SET
		date_of_birth = $2,
		place_of_birth = $3,
		last_education = $4,
		user_id = $5,
		role = $6,
		updated_at = $7
	WHERE
		id = $1
	RETURNING
	id, date_of_birth, place_of_birth, last_education, user_id, role;
	`

	UpdateParticipantByRole = `
	UPDATE
		participants
	SET
		Role = $2
	WHERE
		id = $1
	`

	DeleteParticipantByID = `
	DELETE FROM
		participants
	WHERE
		id = $1;
	`

	GetParticipantBUserID = `
	SELECT
		id,
		date_of_birth,
		place_of_birth,
		last_education,
		user_id
	FROM
		participants
	WHERE
		user_id = $1`
	GetParticipantByUserId        = `SELECT id,date_of_birth,place_of_birth,last_education,user_id,role FROM participants WHERE user_id = $1 ORDER BY created_at DESC;`
	GetScheduleByParticipantID    = `select s.trainer_id, s.participant_id FROM participants p JOIN schedules s ON p.id = s.participant_id WHERE p.id =$1 ;`
	SelectParticipantRoleBasic    = `SELECT id, date_of_birth, place_of_birth,last_education,user_id, role FROM participants WHERE id=$1 ORDER BY created_at DESC`
	SelectParticipantRoleAdvance  = `SELECT id, date_of_birth, place_of_birth,last_education,user_id, role FROM participants WHERE id=$1 ORDER BY created_at DESC`
	GetScheduleByParticipanId     = `select s.trainer_id, s.participant_id FROM participants p JOIN schedules s ON p.id = s.participant_id WHERE p.id =$1 ;`
	SelectParticipantWithSchedule = `SELECT s.id,s.activity, to_char(s.date, 'Day') AS day_of_weeks,s.date,s.trainer_id, s.created_at FROM participants p JOIN schedules s ON p.id = s.participant_id WHERE p.id = $1;`
	InsertQuestionByParticipants  = `INSERT INTO questions (participant_id, question_text, created_at) VALUES ($1, $2, $3) RETURNING id, trainer_id, schedule_id, ;`
	CreateQuestionQuery           = `INSERT INTO questions (question, status, participant_id, schedule_id) VALUES ($1, $2, $3, $4) RETURNING id, participant_id, schedule_id, created_at, updated_at`
	ScheduleIDByParticipantId     = `SELECT id FROM schedules WHERE participant_id = $1;`
	// SelectQuestionList = `SELECT id, question, status, participant_id, created_at, updated_at FROM questions ORDER BY created_at DESC`
	// SelectQuestionByID = `SELECT id, question, status, participant_id, created_at, updated_at FROM questions WHERE id = $1`
	// InsertQuestion     = `INSERT INTO questions (id, question, status, participant_id, created_at, updated_at) VALUES ,($1, $2, $3, $4, $5, $6)`
	// UpdateQuestion     = `UPDATE questions SET question = $1, status = $2, participant_id = $3, updated_at = $4 WHERE id = $5`
	// DeleteQuestion     = `DELETE FROM questions WHERE id = $1`

	InsertAbsence               = `INSERT INTO absences (date, participant_id, trainer_id, schedule_id) VALUES ($1, $2, $3, $4)`
	ListAbsence                 = `SELECT id, date, information, absence_status, absence_time, participant_id, created_at, updated_at FROM absences ORDER BY created_at desc limit $1 offset $2`
	ListAbsencebyDate           = `SELECT id, date, information, absence_status, absence_time, participant_id, created_at, updated_at FROM absences WHERE date >= $1 AND date <= $2 ORDER BY created_at asc limit $3 offset $4`
	GetAbsencesById             = `SELECT id, date, information, absence_status, absence_time, created_at, updated_at FROM absences WHERE participant_id = $1 ORDER BY created_at desc`
	GetAbsencesByScheduleId     = `SELECT id, date, participant_id, trainer_id, schedule_id FROM absences WHERE schedule_id = $1`
	DeleteByParticipantId       = `DELETE FROM absences WHERE participant_id = $1`
	InsertSchedule              = `INSERT INTO schedules (activity, date, trainer_id, participant_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	ListSchedule                = `SELECT id, activity, date, trainer_id, participant_id, created_at, updated_at FROM schedules ORDER BY created_at desc limit $1 offset $2`
	ListScheduleByDate          = `SELECT id, activity, date, trainer_id, participant_id, created_at, updated_at FROM schedules WHERE date >= $1 AND date <= $2 ORDER BY created_at desc limit $3 offset $4`
	ListScheduleByTrainerId     = `SELECT id, activity, date, trainer_id, participant_id, created_at, updated_at FROM schedules WHERE trainer_id = $1 limit $2 offset $3`
	ListScheduleByParticipantId = `SELECT id, activity, date, trainer_id, participant_id, created_at, updated_at FROM schedules WHERE participant_id = $1 limit $2 offset $3`

	InsertScheduleImage = `
	INSERT INTO
		schedule_images(schedule_id, file_name)
	VALUES($1, $2)
	RETURNING
		id,
		schedule_id,
		file_name;
	`
	ScheduleIDByTrainerId           = `SELECT id FROM schedules WHERE trainer_id = $1`
	ScheduleIdByDay                 = `Select id from schedules WHERE EXTRACT(DOW FROM date) = $1 ORDER BY date asc`
	DateByDay                       = `Select date from schedules WHERE EXTRACT(DOW FROM date) = $1 ORDER BY date asc`
	DateByScheduleId                = `Select date from schedules WHERE id = $1`
	ScheduleIDbyDate                = `Select id from schedules WHERE date = $1`
	UpdateScheduleByAdmin           = `Update schedules SET trainer_id = $2 WHERE date = $1 Returning id, activity, date, trainer_id, participant_id, updated_at`
	UpdateAbsencesByParticipantName = `
	UPDATE
		absences
	SET
		information = $3,
		absence_status = $4,
		absence_time = $5,
		updated_at = $6
	WHERE
		participant_id = $1 AND schedule_id = $2
	RETURNING
	id, information, absence_status, absence_time, updated_at`
)
