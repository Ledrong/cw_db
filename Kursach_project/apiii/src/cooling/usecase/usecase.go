package usecase

import (
	"Kursach_project/apiii/models"
	coolingRep "Kursach_project/apiii/src/cooling/repository"
)

type UseCaseI interface {
	CreateCooling(cooling *models.Cooling) error
	UpdateCooling(cpu *models.Cooling, id int64) error
	SelectCooling(id int64) (*models.Cooling, error)
	ShowFullCooling() ([]*models.Cooling, error)
	ShowPartCooling() ([]*models.Cooling, error)
	DeleteCoolingById(id int64) error
}

type useCase struct {
	coolingRepository coolingRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(coolingRepository coolingRep.RepositoryI) UseCaseI {
	return &useCase{
		coolingRepository: coolingRepository,
		//userRepository:  userRepository,
	}
}

func (u *useCase) CreateCooling(cooling *models.Cooling) error {
	existcooling, e := u.coolingRepository.SelectCoolingById(cooling.Coolingid)

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		cooling.Coolingid = existcooling.Coolingid
		cooling.Name = existcooling.Name
		cooling.Brand = existcooling.Brand
		cooling.MaxSpeed = existcooling.MaxSpeed
		cooling.CountVentilators = existcooling.CountVentilators
		cooling.Price = existcooling.Price
		return models.ErrConflict

	}

	e = u.coolingRepository.CreateCooling(cooling)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdateCooling(cooling *models.Cooling, id int64) error {
	_, e := u.coolingRepository.SelectCoolingById(id)
	if e != nil {
		return e
	}

	e = u.coolingRepository.UpdateCooling(cooling, id)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectCooling(id int64) (*models.Cooling, error) {
	cooling, e := u.coolingRepository.SelectCoolingById(id)

	if e != nil {
		return nil, e
	}

	return cooling, nil
}

func (u *useCase) ShowFullCooling() ([]*models.Cooling, error) {
	coolingfulltable, e := u.coolingRepository.ShowFullCooling()
	if e != nil {
		return nil, e
	}
	return coolingfulltable, nil
}

func (u *useCase) ShowPartCooling() ([]*models.Cooling, error) {
	coolingparttable, e := u.coolingRepository.ShowPartCooling()
	if e != nil {
		return nil, e
	}
	return coolingparttable, nil
}

func (u *useCase) DeleteCoolingById(id int64) error {
	e := u.coolingRepository.DeleteCoolingById(id)
	if e != nil {
		return e
	}

	return nil
}
