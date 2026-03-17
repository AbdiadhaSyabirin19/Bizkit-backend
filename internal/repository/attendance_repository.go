package repository

import (
	"time"

	"bizkit-backend/config"
	"bizkit-backend/internal/model"
)

func GetTodayAttendance(userID uint) (*model.Attendance, error) {
	var attendance model.Attendance
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	err := config.DB.Preload("User").
		Where("user_id = ? AND check_in >= ? AND check_in < ?", userID, today, tomorrow).
		First(&attendance).Error
	return &attendance, err
}

func GetAttendanceByID(id uint) (*model.Attendance, error) {
	var attendance model.Attendance
	err := config.DB.Preload("User").First(&attendance, id).Error
	return &attendance, err
}

func CreateAttendance(attendance *model.Attendance) error {
	return config.DB.Create(attendance).Error
}

func UpdateAttendance(attendance *model.Attendance) error {
	return config.DB.Save(attendance).Error
}

func GetAttendanceHistory(userID uint, limit int) ([]model.Attendance, error) {
	var list []model.Attendance
	err := config.DB.Preload("User").
		Where("user_id = ?", userID).
		Order("check_in DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// GetAttendanceByDate — untuk laporan admin: semua absensi pada tanggal tertentu
func GetAttendanceByDate(date time.Time) ([]model.Attendance, error) {
	var list []model.Attendance
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	err := config.DB.
		Preload("User").
		Preload("User.Role").
		Preload("User.Outlet").
		Where("check_in >= ? AND check_in < ?", start, end).
		Order("check_in ASC").
		Find(&list).Error
	return list, err
}