package usecase

import (
	"Kursach_project/apiii/models"
	cpuRep "Kursach_project/apiii/src/cpu/repository"
	ramRep "Kursach_project/apiii/src/ram/repository"
)

type UseCaseI interface {
	CreateRam(ram *models.Ram) error
	UpdateRam(ram *models.Ram, id int64) error
	SelectRam(id int64) (*models.Ram, error)
	ShowFullRam() ([]*models.Ram, error)
	ShowCompatibilityCpu(id int64) ([]*models.Cpu, error)
	ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error)
	ShowPartRam() ([]*models.Ram, error)
	DeleteRamById(id int64) error
}

type useCase struct {
	ramRepository ramRep.RepositoryI
	cpuRepository cpuRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(ramRepository ramRep.RepositoryI, cpuRepository cpuRep.RepositoryI) UseCaseI {
	return &useCase{
		ramRepository: ramRepository,
		cpuRepository: cpuRepository,

		//userRepository:  userRepository,
	}
}

func (u *useCase) CreateRam(ram *models.Ram) error {
	existram, e := u.ramRepository.SelectRamById(ram.Ramid)

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		ram.Ramid = existram.Ramid
		ram.Name = existram.Name
		ram.Brand = existram.Brand
		ram.Rammemory = existram.Rammemory
		ram.Ddr = existram.Ddr
		ram.Price = existram.Price
		return models.ErrConflict

	}

	e = u.ramRepository.CreateRam(ram)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdateRam(ram *models.Ram, id int64) error {
	_, e := u.ramRepository.SelectRamById(id)
	if e != nil {
		return e
	}

	e = u.ramRepository.UpdateRam(ram, id)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectRam(id int64) (*models.Ram, error) {
	ram, e := u.ramRepository.SelectRamById(id)

	if e != nil {
		return nil, e
	}

	return ram, nil
}

// почему-то какая-то параша как будто бы, мб не стоит передавать параметры в showfullram
func (u *useCase) ShowFullRam() ([]*models.Ram, error) {
	ramfulltable, e := u.ramRepository.ShowFullRam()
	if e != nil {
		return nil, e
	}
	return ramfulltable, nil
}

func (u *useCase) ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error) {
	motherboardfulltable, e := u.ramRepository.ShowCompatibilityMotherboard(id)
	if e != nil {
		return nil, e
	}
	return motherboardfulltable, nil
}

func (u *useCase) ShowCompatibilityCpu(id int64) ([]*models.Cpu, error) {
	cpufulltable, e := u.ramRepository.ShowCompatibilityCpu(id)
	if e != nil {
		return nil, e
	}
	return cpufulltable, nil
}

func (u *useCase) ShowPartRam() ([]*models.Ram, error) {
	ramparttable, e := u.ramRepository.ShowPartRam()
	if e != nil {
		return nil, e
	}
	return ramparttable, nil
}

func (u *useCase) DeleteRamById(id int64) error {
	e := u.ramRepository.DeleteRamById(id)
	if e != nil {
		return e
	}

	return nil
}
