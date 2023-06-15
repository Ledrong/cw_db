package usecase

import (
	"Kursach_project/apiii/models"
	cpuRep "Kursach_project/apiii/src/cpu/repository"
	ramRep "Kursach_project/apiii/src/ram/repository"
)

type UseCaseI interface {
	CreateCpu(cpu *models.Cpu) error
	UpdateCpu(cpu *models.Cpu, id int64) error
	SelectCpu(id int64) (*models.Cpu, error)
	ShowFullCpu() ([]*models.Cpu, error)
	ShowCompatibilityRam(id int64) ([]*models.Ram, error)
	ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error)
	ShowPartCpu() ([]*models.Cpu, error)
	DeleteCpuById(id int64) error
}

type useCase struct {
	cpuRepository cpuRep.RepositoryI
	ramRepository ramRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(cpuRepository cpuRep.RepositoryI, ramRepository ramRep.RepositoryI) UseCaseI {
	return &useCase{
		cpuRepository: cpuRepository,
		ramRepository: ramRepository,
		//userRepository:  userRepository,
	}
}

// нужен ли user и нафига existcpu
func (u *useCase) CreateCpu(cpu *models.Cpu) error {
	existcpu, e := u.cpuRepository.SelectCpuById(cpu.Cpuid)

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		cpu.Cpuid = existcpu.Cpuid
		cpu.Name = existcpu.Name
		cpu.Brand = existcpu.Brand
		cpu.Series = existcpu.Series
		cpu.Model = existcpu.Model
		cpu.SupportedDdr = existcpu.SupportedDdr
		cpu.Price = existcpu.Price
		return models.ErrConflict

	}

	e = u.cpuRepository.CreateCpu(cpu)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdateCpu(cpu *models.Cpu, id int64) error {
	_, e := u.cpuRepository.SelectCpuById(id)
	if e != nil {
		return e
	}

	e = u.cpuRepository.UpdateCpu(cpu, id)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectCpu(id int64) (*models.Cpu, error) {
	cpu, e := u.cpuRepository.SelectCpuById(id)

	if e != nil {
		return nil, e
	}

	return cpu, nil
}

// почему-то какая-то параша как будто бы, мб не стоит передавать параметры в showfullcpu
func (u *useCase) ShowFullCpu() ([]*models.Cpu, error) {
	cpufulltable, e := u.cpuRepository.ShowFullCpu()
	if e != nil {
		return nil, e
	}
	return cpufulltable, nil
}

func (u *useCase) ShowCompatibilityRam(id int64) ([]*models.Ram, error) {
	ramfulltable, e := u.cpuRepository.ShowCompatibilityRam(id)
	if e != nil {
		return nil, e
	}
	return ramfulltable, nil
}

func (u *useCase) ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error) {
	motherboardfulltable, e := u.cpuRepository.ShowCompatibilityMotherboard(id)
	if e != nil {
		return nil, e
	}
	return motherboardfulltable, nil
}

func (u *useCase) ShowPartCpu() ([]*models.Cpu, error) {
	cpuparttable, e := u.cpuRepository.ShowPartCpu()
	if e != nil {
		return nil, e
	}
	return cpuparttable, nil
}

func (u *useCase) DeleteCpuById(id int64) error {
	e := u.cpuRepository.DeleteCpuById(id)
	if e != nil {
		return e
	}

	return nil
}
