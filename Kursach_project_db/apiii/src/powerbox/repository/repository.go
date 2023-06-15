package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreatePowerbox(powerbox *models.Powerbox) error
	UpdatePowerbox(powerbox *models.Powerbox, id int64) error
	SelectPowerboxById(id int64) (*models.Powerbox, error)
	ShowFullPowerbox() ([]*models.Powerbox, error)
	ShowPartPowerbox() ([]*models.Powerbox, error)
	DeletePowerboxById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbPowerbox dataBase) CreatePowerbox(powerbox *models.Powerbox) error {
	var id int64
	e := dbPowerbox.db.QueryRow(`INSERT INTO powerbox (name, brand, power, form_factor, price) VALUES ($1, $2, $3, $4, $5) RETURNING powerboxid`, powerbox.Name, powerbox.Brand, powerbox.Power, powerbox.FormFactor, powerbox.Price).Scan(&id)
	if e != nil {
		return e
	}

	powerbox.Powerboxid = id
	return nil
}

func (dbPowerbox dataBase) UpdatePowerbox(powerbox *models.Powerbox, id int64) error {
	res, e := dbPowerbox.db.Exec(`UPDATE powerbox SET  name = $1, brand = $2, power = $3, form_factor = $4, price = $5 WHERE powerboxid = $6`, powerbox.Name, powerbox.Brand, powerbox.Power, powerbox.FormFactor, powerbox.Price, id)
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
	//powerbox.Powerboxid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	return nil
}

// спросить у Настиииииииии что делать approved
func (dbPowerbox dataBase) SelectPowerboxById(id int64) (*models.Powerbox, error) {
	powerbox := models.Powerbox{}

	q := dbPowerbox.db.QueryRow(`SELECT * FROM powerbox WHERE powerboxid = $1`, id).Scan(&powerbox.Powerboxid, &powerbox.Name, &powerbox.Brand, &powerbox.Power, &powerbox.FormFactor, &powerbox.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &powerbox, nil
}

func (dbPowerbox dataBase) ShowFullPowerbox() ([]*models.Powerbox, error) {
	rows, e := dbPowerbox.db.Query(`SELECT * FROM powerbox`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Powerbox, 0)
	for rows.Next() {
		var powerbox models.Powerbox
		e = rows.Scan(&powerbox.Powerboxid, &powerbox.Name, &powerbox.Brand, &powerbox.Power, &powerbox.FormFactor, &powerbox.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &powerbox)
	}

	return messiv, e
}

func (dbPowerbox dataBase) ShowPartPowerbox() ([]*models.Powerbox, error) {
	rows, e := dbPowerbox.db.Query(`SELECT powerboxid, name, price FROM powerbox`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Powerbox, 0)
	for rows.Next() {
		var powerbox models.Powerbox
		e = rows.Scan(&powerbox.Powerboxid, &powerbox.Name, &powerbox.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &powerbox)
	}

	return messiv, e
}

// /спросить у Насти
func (dbPowerbox dataBase) DeletePowerboxById(id int64) error {
	//_, e := dbPowerbox.db.Exec(`DELETE FROM powerbox WHERE powerboxid = $1`, id)
	res, e := dbPowerbox.db.Exec(`DELETE FROM powerbox WHERE powerboxid = $1`, id)

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
	//powerbox.Powerboxid, e = q.LastInsertId()

	return nil
}
