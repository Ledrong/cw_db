package usecase

import (
	"Kursach_project/apiii/models"
	userRep "Kursach_project/apiii/src/user/repository"
)

type UseCaseI interface {
	CreateUser(user *models.User) error
	Login(user *models.User) error
	UpdateUser(user *models.User, id int64) error
	GetroleUser(id int64) (string, error)
	SelectUser(id int64) (*models.User, error)
	ShowFullUser() ([]*models.User, error)
	DeleteUserById(id int64) error
}

type useCase struct {
	userRepository userRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(userRepository userRep.RepositoryI) UseCaseI {
	return &useCase{
		userRepository: userRepository,
	}
}

func (u *useCase) CreateUser(user *models.User) error {
	existuser, e := u.userRepository.SelectUserById(user.Userid)

	//if e != models.ErrNotFound && e != nil {
	//	return e
	//} else if e == nil {
	//	//можно не заполнять хуйню эту
	//	user.Userid = existuser.Userid
	//	user.Login = existuser.Login
	//	user.Password = existuser.Password
	//	user.Role = existuser.Role
	//	return models.ErrConflict
	//}

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		user.Userid = existuser.Userid
		user.Login = existuser.Login
		user.Password = existuser.Password
		user.Role = existuser.Role
		return models.ErrConflict
	}

	e = u.userRepository.CreateUser(user)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdateUser(user *models.User, id int64) error {
	_, e := u.userRepository.SelectUserById(id)
	if e != nil {
		return e
	}

	e = u.userRepository.UpdateUser(user, id)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) Login(user *models.User) error {
	e := u.userRepository.Login(user)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectUser(id int64) (*models.User, error) {
	user, e := u.userRepository.SelectUserById(id)

	if e != nil {
		return nil, e
	}

	return user, nil
}

func (u *useCase) ShowFullUser() ([]*models.User, error) {
	userfulltable, e := u.userRepository.ShowFullUser()
	if e != nil {
		return nil, e
	}
	return userfulltable, nil
}

func (u *useCase) GetroleUser(id int64) (string, error) {
	user, e := u.userRepository.GetroleUserById(id)
	if e != nil {
		return "", e
	}
	return user, nil
}

func (u *useCase) DeleteUserById(id int64) error {
	e := u.userRepository.DeleteUserById(id)
	if e != nil {
		return e
	}

	return nil
}
