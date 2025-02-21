package user

import "LinkShorty/pkg/db"

type UserRepository struct {
	DataBase *db.Db
}

func NewUserRepository(dataBase *db.Db) *UserRepository {
	return &UserRepository{
		DataBase: dataBase,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.DataBase.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.DataBase.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
