package stat

import (
	"LinkShorty/pkg/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{db}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	result := repo.Db.Where("link_id = ? AND data = ?", linkId, currentDate).First(&stat)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		// Если запись не найдена - создаем новую
		repo.Db.Create(&Stat{
			LinkID: linkId,
			Clicks: 1,
			Data:   currentDate,
		})
	} else {
		stat.Clicks++
		repo.Db.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string
	switch by {
	case GroupByDay:
		selectQuery = "to_char(data, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(data, 'YYYY-MM') as period, sum(clicks)"
	}
	query := repo.DB.Table("stats").
		Select(selectQuery).
		Session(&gorm.Session{})

	if true {
		query.Where("count > 10")
	}

	query.Where("data BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	return stats
}
