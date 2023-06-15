package usecase

import (
	"Kursach_project/apiii/models"
	motherboardRep "Kursach_project/apiii/src/motherboard/repository"
)

type UseCaseI interface {
	CreateMotherboard(motherboard *models.Motherboard) error
	UpdateMotherboard(motherboard *models.Motherboard, id int64) error
	SelectMotherboard(id int64) (*models.Motherboard, error)
	ShowFullMotherboard() ([]*models.Motherboard, error)
	ShowCompatibilityRam(id int64) ([]*models.Ram, error)
	ShowCompatibilityCpu(id int64) ([]*models.Cpu, error)
	ShowPartMotherboard() ([]*models.Motherboard, error)
	DeleteMotherboardById(id int64) error
}

type useCase struct {
	motherboardRepository motherboardRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(motherboardRepository motherboardRep.RepositoryI) UseCaseI {
	return &useCase{
		motherboardRepository: motherboardRepository,
	}
}

func (u *useCase) CreateMotherboard(motherboard *models.Motherboard) error {
	existmotherboard, e := u.motherboardRepository.SelectMotherboardById(motherboard.Motherboardid)

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		motherboard.Motherboardid = existmotherboard.Motherboardid
		motherboard.Name = existmotherboard.Name
		motherboard.Brand = existmotherboard.Brand
		motherboard.CountSlots = existmotherboard.CountSlots
		motherboard.Price = existmotherboard.Price
		return models.ErrConflict
	}

	e = u.motherboardRepository.CreateMotherboard(motherboard)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdateMotherboard(motherboard *models.Motherboard, id int64) error {
	_, e := u.motherboardRepository.SelectMotherboardById(id)
	if e != nil {
		return e
	}

	e = u.motherboardRepository.UpdateMotherboard(motherboard, id)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectMotherboard(id int64) (*models.Motherboard, error) {
	motherboard, e := u.motherboardRepository.SelectMotherboardById(id)

	if e != nil {
		return nil, e
	}

	return motherboard, nil
}

func (u *useCase) ShowFullMotherboard() ([]*models.Motherboard, error) {
	motherboardfulltable, e := u.motherboardRepository.ShowFullMotherboard()
	if e != nil {
		return nil, e
	}
	return motherboardfulltable, nil
}

func (u *useCase) ShowCompatibilityRam(id int64) ([]*models.Ram, error) {
	ramfulltable, e := u.motherboardRepository.ShowCompatibilityRam(id)
	if e != nil {
		return nil, e
	}
	return ramfulltable, nil
}

func (u *useCase) ShowCompatibilityCpu(id int64) ([]*models.Cpu, error) {
	cpufulltable, e := u.motherboardRepository.ShowCompatibilityCpu(id)
	if e != nil {
		return nil, e
	}
	return cpufulltable, nil
}

func (u *useCase) ShowPartMotherboard() ([]*models.Motherboard, error) {
	motherboardparttable, e := u.motherboardRepository.ShowPartMotherboard()
	if e != nil {
		return nil, e
	}
	return motherboardparttable, nil
}

func (u *useCase) DeleteMotherboardById(id int64) error {
	e := u.motherboardRepository.DeleteMotherboardById(id)
	if e != nil {
		return e
	}

	return nil
}
