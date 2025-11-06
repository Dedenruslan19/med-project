package logs

import (
	"log/slog"

	"Dedenruslan19/med-project/service/logs"

	"gorm.io/gorm"
)

type logRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewLogsRepo(db *gorm.DB, logger *slog.Logger) logs.LogRepo {
	return &logRepo{db: db, logger: logger}
}

func (r *logRepo) Create(l *logs.ExerciseLog) (int64, error) {
	if err := r.db.Create(&l).Error; err != nil {
		r.logger.Error("failed to create exercise log",
			"user_id", l.UserID,
			"exercise_id", l.ExerciseID,
			"error", err)
		return 0, err
	}
	return l.ID, nil
}

func (r *logRepo) GetByUserID(userID int64) ([]logs.ExerciseLog, error) {
	var userLogs []logs.ExerciseLog
	err := r.db.Where("user_id = ?", userID).Find(&userLogs).Error
	if err != nil {
		r.logger.Error("failed to get exercise logs by user_id",
			"user_id", userID,
			"error", err)
		return nil, err
	}
	return userLogs, nil
}

func (r *logRepo) GetByID(logID int64) (*logs.ExerciseLog, error) {
	var logEntry logs.ExerciseLog
	err := r.db.First(&logEntry, logID).Error
	if err != nil {
		r.logger.Error("failed to get exercise log by id",
			"log_id", logID,
			"error", err)
		return nil, err
	}
	return &logEntry, nil
}

func (r *logRepo) GetAll() ([]logs.ExerciseLog, error) {
	var allLogs []logs.ExerciseLog
	err := r.db.Find(&allLogs).Error
	if err != nil {
		r.logger.Error("failed to get all exercise logs", "error", err)
		return nil, err
	}
	return allLogs, nil
}

func (r *logRepo) Delete(logID int64) error {
	result := r.db.Delete(&logs.ExerciseLog{}, logID)
	if result.Error != nil {
		r.logger.Error("failed to delete exercise log",
			"log_id", logID,
			"error", result.Error)
		return result.Error
	}
	return nil
}
