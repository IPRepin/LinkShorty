package link

import (
	"LinkShorty/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	DataBase *db.Db
}

func NewLinkRepository(dataBase *db.Db) *LinkRepository {
	return &LinkRepository{
		DataBase: dataBase,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.DataBase.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) GetLinksCount() (int64, error) {
	var count int64
	result := repo.DataBase.Table("links").
		Where("deleted_at is null").
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (repo *LinkRepository) GetAllLinks(limit, offset int) ([]Link, error) {
	var links []Link
	query := repo.DataBase.Table("links").
		Where("deleted_at is null").
		Session(&gorm.Session{})
	result := query.
		Order("id asc").
		Limit(limit).
		Offset(offset).
		Scan(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}
