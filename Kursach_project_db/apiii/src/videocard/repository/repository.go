package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateVideocard(videocard *models.Videocard) error
	UpdateVideocard(videocard *models.Videocard, id int64) error
	SelectVideocardById(id int64) (*models.Videocard, error)
	ShowFullVideocard() ([]*models.Videocard, error)
	ShowPartVideocard() ([]*models.Videocard, error)
	DeleteVideocardById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbVideocard dataBase) CreateVideocard(videocard *models.Videocard) error {
	var id int64
	e := dbVideocard.db.QueryRow(`INSERT INTO videocard (name, brand, series, vmemory, price) VALUES ($1, $2, $3, $4, $5) RETURNING videocardid`, videocard.Name, videocard.Brand, videocard.Series, videocard.Vmemory, videocard.Price).Scan(&id)
	if e != nil {
		return e
	}

	videocard.Videocardid = id
	return nil
}

func (dbVideocard dataBase) UpdateVideocard(videocard *models.Videocard, id int64) error {
	res, e := dbVideocard.db.Exec(`UPDATE videocard SET name = $1, brand = $2, series = $3, vmemory = $4, price = $5 WHERE videocardid = $6`, videocard.Name, videocard.Brand, videocard.Series, videocard.Vmemory, videocard.Price, id)
	if e != nil {
		return e
	}

	rows, e := res.RowsAffected()

	if e != nil {
		return e
	}
	if rows != 1 {
		return models.ErrNotFound
	}
	//videocard.Videocardid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	return nil
}

// спросить у Настиииииииии что делать approved
func (dbVideocard dataBase) SelectVideocardById(id int64) (*models.Videocard, error) {
	videocard := models.Videocard{}

	q := dbVideocard.db.QueryRow(`SELECT * FROM videocard WHERE videocardid = $1`, id).Scan(&videocard.Videocardid, &videocard.Name, &videocard.Brand, &videocard.Series, &videocard.Vmemory, &videocard.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &videocard, nil
}

func (dbVideocard dataBase) ShowFullVideocard() ([]*models.Videocard, error) {
	rows, e := dbVideocard.db.Query(`SELECT * FROM videocard`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Videocard, 0)
	for rows.Next() {
		var videocard models.Videocard
		e = rows.Scan(&videocard.Videocardid, &videocard.Name, &videocard.Brand, &videocard.Series, &videocard.Vmemory, &videocard.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &videocard)
	}

	return messiv, e
}

func (dbVideocard dataBase) ShowPartVideocard() ([]*models.Videocard, error) {
	rows, e := dbVideocard.db.Query(`SELECT videocardid, name, price FROM videocard`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Videocard, 0)
	for rows.Next() {
		var videocard models.Videocard
		e = rows.Scan(&videocard.Videocardid, &videocard.Name, &videocard.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &videocard)
	}

	return messiv, e
}

// /спросить у Насти надо ли передавать что-то кроме id
func (dbVideocard dataBase) DeleteVideocardById(id int64) error {
	//_, e := dbVideocard.db.Exec(`DELETE FROM videocard WHERE videocardid = $1`, id)
	res, e := dbVideocard.db.Exec(`DELETE FROM videocard WHERE videocardid = $1`, id)

	if e != nil {
		return e
	}

	rows, e := res.RowsAffected()

	if e != nil {
		return e
	}
	if rows != 1 {
		return models.ErrNotFound
	}
	//videocard.Videocardid, e = q.LastInsertId()

	return nil
}
