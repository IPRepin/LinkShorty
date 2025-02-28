package di

import "LinkShorty/internal/user"

type IStatRepository interface {
	AddClick(linkID uint)
}

type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}
