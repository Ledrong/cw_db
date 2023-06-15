package usecase

import (
	"Kursach_project/apiii/models"
	coolingRep "Kursach_project/apiii/src/cooling/repository"
	cpuRep "Kursach_project/apiii/src/cpu/repository"
	motherboardRep "Kursach_project/apiii/src/motherboard/repository"
	pcbuildRep "Kursach_project/apiii/src/pcbuild/repository"
	powerboxRep "Kursach_project/apiii/src/powerbox/repository"
	ramRep "Kursach_project/apiii/src/ram/repository"
	userRep "Kursach_project/apiii/src/user/repository"
	videocardRep "Kursach_project/apiii/src/videocard/repository"
)

type UseCaseI interface {
	CreatePcbuild(pcbuild *models.Pcbuild, userid int64) error
	UpdatePcbuild(pcbuild *models.Pcbuild, userid int64, pcbuildid int64) error
	SelectPcbuildById(id int64) (*models.Pcbuild, error)
	ShowPcbuildById(id int64) (*models.Pcbuildinfo, error)
	ShowFullPcbuild() ([]*models.Pcbuildinfo, error)
	ShowMyPcbuild(id int64) ([]*models.Pcbuildinfo, error)
	DeletePcbuildById(userid int64, id int64) error
}

type useCase struct {
	pcbuildRepository     pcbuildRep.RepositoryI
	cpuRepository         cpuRep.RepositoryI
	ramRepository         ramRep.RepositoryI
	videocardRepository   videocardRep.RepositoryI
	powerboxRepository    powerboxRep.RepositoryI
	motherboardRepository motherboardRep.RepositoryI
	coolingRepository     coolingRep.RepositoryI
	userRepository        userRep.RepositoryI
}

func New(pcbuildRepository pcbuildRep.RepositoryI, cpuRepository cpuRep.RepositoryI, ramRepository ramRep.RepositoryI, videocardRepository videocardRep.RepositoryI,
	powerboxRepository powerboxRep.RepositoryI, motherboardRepository motherboardRep.RepositoryI, coolingRepository coolingRep.RepositoryI, userRepository userRep.RepositoryI) UseCaseI {
	return &useCase{
		pcbuildRepository:     pcbuildRepository,
		cpuRepository:         cpuRepository,
		ramRepository:         ramRepository,
		videocardRepository:   videocardRepository,
		powerboxRepository:    powerboxRepository,
		motherboardRepository: motherboardRepository,
		coolingRepository:     coolingRepository,
		userRepository:        userRepository,
	}
}

func (u useCase) CreatePcbuild(pcbuild *models.Pcbuild, userid int64) error {

	//_, e := u.userRepository.SelectUserById(userid)
	//if e != nil && userid != 0 {
	//	return e
	//}

	_, e := u.cpuRepository.SelectCpuById(pcbuild.Cpuid)
	if e != nil {
		return e
	}

	_, e = u.ramRepository.SelectRamById(pcbuild.Ramid)
	if e != nil {
		return e
	}

	_, e = u.powerboxRepository.SelectPowerboxById(pcbuild.Powerboxid)
	if e != nil {
		return e
	}

	_, e = u.motherboardRepository.SelectMotherboardById(pcbuild.Motherboardid)
	if e != nil {
		return e
	}

	_, e = u.coolingRepository.SelectCoolingById(pcbuild.Coolingid)
	if e != nil {
		return e
	}

	_, e = u.videocardRepository.SelectVideocardById(pcbuild.Videocardid)
	if e != nil {
		return e
	}

	existpcbuild, e := u.pcbuildRepository.SelectPcbuildById(pcbuild.Pcbuildid)
	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		pcbuild.Pcbuildid = existpcbuild.Pcbuildid
		pcbuild.Userid = existpcbuild.Userid
		pcbuild.Cpuid = existpcbuild.Cpuid
		pcbuild.Ramid = existpcbuild.Ramid
		pcbuild.Powerboxid = existpcbuild.Powerboxid
		pcbuild.Motherboardid = pcbuild.Motherboardid
		pcbuild.Coolingid = existpcbuild.Coolingid
		pcbuild.Videocardid = existpcbuild.Videocardid
		//pcbuild.Compatibility = existpcbuild.Compatibility
		return models.ErrConflict
	}

	e = u.pcbuildRepository.CreatePcbuild(pcbuild, userid)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdatePcbuild(pcbuild *models.Pcbuild, userid int64, pcbuildid int64) error {
	_, e := u.pcbuildRepository.SelectPcbuildById(pcbuildid)
	if e != nil {
		return e
	}

	e = u.pcbuildRepository.UpdatePcbuild(pcbuild, userid, pcbuildid)
	if e != nil {
		return e
	}

	return nil
}

func (u useCase) SelectPcbuildById(id int64) (*models.Pcbuild, error) {
	pcbuild, e := u.pcbuildRepository.SelectPcbuildById(id)

	if e != nil {
		return nil, e
	}

	return pcbuild, nil
}

func (u useCase) ShowPcbuildById(id int64) (*models.Pcbuildinfo, error) {
	pcbuildinfo, e := u.pcbuildRepository.ShowPcbuildById(id)

	if e != nil {
		return nil, e
	}

	return pcbuildinfo, nil
}

func (u useCase) ShowFullPcbuild() ([]*models.Pcbuildinfo, error) {
	pcbuildfulltable, e := u.pcbuildRepository.ShowFullPcbuild()
	if e != nil {
		return nil, e
	}
	return pcbuildfulltable, nil
}

func (u useCase) ShowMyPcbuild(id int64) ([]*models.Pcbuildinfo, error) {
	pcbuildfulltable, e := u.pcbuildRepository.ShowMyPcbuild(id)
	if e != nil {
		return nil, e
	}
	return pcbuildfulltable, nil
}

func (u useCase) DeletePcbuildById(userid int64, id int64) error {
	e := u.pcbuildRepository.DeletePcbuildById(userid, id)
	if e != nil {
		return e
	}

	return nil
}
