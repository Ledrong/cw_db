package usecase

import (
	"Kursach_project/apiii/models"
	powerboxRep "Kursach_project/apiii/src/powerbox/repository"
)

type UseCaseI interface {
	CreatePowerbox(powerbox *models.Powerbox) error
	UpdatePowerbox(powerbox *models.Powerbox, id int64) error
	SelectPowerbox(id int64) (*models.Powerbox, error)
	ShowFullPowerbox() ([]*models.Powerbox, error)
	ShowPartPowerbox() ([]*models.Powerbox, error)
	DeletePowerboxById(id int64) error
}

type useCase struct {
	powerboxRepository powerboxRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(powerboxRepository powerboxRep.RepositoryI) UseCaseI {
	return &useCase{
		powerboxRepository: powerboxRepository,
		//userRepository:  userRepository,
	}
}

// нужен ли user и нафига existpowerbox
func (u *useCase) CreatePowerbox(powerbox *models.Powerbox) error {
	existpowerbox, e := u.powerboxRepository.SelectPowerboxById(powerbox.Powerboxid)

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		powerbox.Powerboxid = existpowerbox.Powerboxid
		powerbox.Name = existpowerbox.Name
		powerbox.Brand = existpowerbox.Brand
		powerbox.Power = existpowerbox.Power
		powerbox.FormFactor = existpowerbox.FormFactor
		powerbox.Price = existpowerbox.Price
		return models.ErrConflict

	}

	e = u.powerboxRepository.CreatePowerbox(powerbox)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectPowerbox(id int64) (*models.Powerbox, error) {
	powerbox, e := u.powerboxRepository.SelectPowerboxById(id)

	if e != nil {
		return nil, e
	}

	return powerbox, nil
}

func (u *useCase) UpdatePowerbox(powerbox *models.Powerbox, id int64) error {
	_, e := u.powerboxRepository.SelectPowerboxById(id)
	if e != nil {
		return e
	}

	e = u.powerboxRepository.UpdatePowerbox(powerbox, id)
	if e != nil {
		return e
	}

	return nil
}

// почему-то какая-то параша как будто бы, мб не стоит передавать параметры в showfullpowerbox
func (u *useCase) ShowFullPowerbox() ([]*models.Powerbox, error) {
	powerboxfulltable, e := u.powerboxRepository.ShowFullPowerbox()
	if e != nil {
		return nil, e
	}
	return powerboxfulltable, nil
}

func (u *useCase) ShowPartPowerbox() ([]*models.Powerbox, error) {
	powerboxparttable, e := u.powerboxRepository.ShowPartPowerbox()
	if e != nil {
		return nil, e
	}
	return powerboxparttable, nil
}

func (u *useCase) DeletePowerboxById(id int64) error {
	e := u.powerboxRepository.DeletePowerboxById(id)
	if e != nil {
		return e
	}

	return nil
}
