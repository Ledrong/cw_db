package usecase

import (
	"Kursach_project/apiii/models"
	videocardRep "Kursach_project/apiii/src/videocard/repository"
)

type UseCaseI interface {
	CreateVideocard(videocard *models.Videocard) error
	UpdateVideocard(videocard *models.Videocard, id int64) error
	SelectVideocard(id int64) (*models.Videocard, error)
	ShowFullVideocard() ([]*models.Videocard, error)
	ShowPartVideocard() ([]*models.Videocard, error)
	DeleteVideocardById(id int64) error
}

type useCase struct {
	videocardRepository videocardRep.RepositoryI
	//userRepository  userRep.RepositoryI
}

func New(videocardRepository videocardRep.RepositoryI) UseCaseI {
	return &useCase{
		videocardRepository: videocardRepository,
		//userRepository:  userRepository,
	}
}

// нужен ли user и нафига existvideocard
func (u *useCase) CreateVideocard(videocard *models.Videocard) error {
	existvideocard, e := u.videocardRepository.SelectVideocardById(videocard.Videocardid)

	if e != models.ErrNotFound && e != nil {
		return e
	} else if e == nil {
		//можно не заполнять хуйню эту
		videocard.Videocardid = existvideocard.Videocardid
		videocard.Name = existvideocard.Name
		videocard.Brand = existvideocard.Brand
		videocard.Series = existvideocard.Series
		videocard.Vmemory = existvideocard.Vmemory
		videocard.Price = existvideocard.Price
		return models.ErrConflict

	}

	e = u.videocardRepository.CreateVideocard(videocard)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) UpdateVideocard(videocard *models.Videocard, id int64) error {
	_, e := u.videocardRepository.SelectVideocardById(id)
	if e != nil {
		return e
	}

	e = u.videocardRepository.UpdateVideocard(videocard, id)
	if e != nil {
		return e
	}

	return nil
}

func (u *useCase) SelectVideocard(id int64) (*models.Videocard, error) {
	videocard, e := u.videocardRepository.SelectVideocardById(id)

	if e != nil {
		return nil, e
	}

	return videocard, nil
}

// почему-то какая-то параша как будто бы, мб не стоит передавать параметры в showfullvideocard
func (u *useCase) ShowFullVideocard() ([]*models.Videocard, error) {
	videocardfulltable, e := u.videocardRepository.ShowFullVideocard()
	if e != nil {
		return nil, e
	}
	return videocardfulltable, nil
}

func (u *useCase) ShowPartVideocard() ([]*models.Videocard, error) {
	videocardparttable, e := u.videocardRepository.ShowPartVideocard()
	if e != nil {
		return nil, e
	}
	return videocardparttable, nil
}

func (u *useCase) DeleteVideocardById(id int64) error {
	e := u.videocardRepository.DeleteVideocardById(id)
	if e != nil {
		return e
	}

	return nil
}
