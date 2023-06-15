package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateCooling(cooling *models.Cooling) error
	UpdateCooling(cpu *models.Cooling, id int64) error
	SelectCoolingById(id int64) (*models.Cooling, error)
	ShowFullCooling() ([]*models.Cooling, error)
	ShowPartCooling() ([]*models.Cooling, error)
	DeleteCoolingById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbCooling dataBase) CreateCooling(cooling *models.Cooling) error {
	var id int64
	e := dbCooling.db.QueryRow(`INSERT INTO cooling (name, brand, max_speed, count_ventilators, price) VALUES ($1, $2, $3, $4, $5) RETURNING coolingid`, cooling.Name, cooling.Brand, cooling.MaxSpeed, cooling.CountVentilators, cooling.Price).Scan(&id)
	if e != nil {
		return e
	}
	cooling.Coolingid = id

	return nil
}

func (dbCooling dataBase) UpdateCooling(cooling *models.Cooling, id int64) error {
	res, e := dbCooling.db.Exec(`UPDATE cooling SET  name = $1, brand = $2, max_speed = $3, count_ventilators = $4, price = $5 WHERE coolingid = $6`, cooling.Name, cooling.Brand, cooling.MaxSpeed, cooling.CountVentilators, cooling.Price, id)
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
	//cooling.coolingid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	return nil
}

func (dbCooling dataBase) SelectCoolingById(id int64) (*models.Cooling, error) {
	cooling := models.Cooling{}

	q := dbCooling.db.QueryRow(`SELECT * FROM cooling WHERE coolingid = $1`, id).Scan(&cooling.Coolingid, &cooling.Name, &cooling.Brand, &cooling.MaxSpeed, &cooling.CountVentilators, &cooling.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &cooling, nil
}

func (dbCooling dataBase) ShowFullCooling() ([]*models.Cooling, error) {
	rows, e := dbCooling.db.Query(`SELECT * FROM cooling`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Cooling, 0)
	for rows.Next() {
		var cooling models.Cooling
		e = rows.Scan(&cooling.Coolingid, &cooling.Name, &cooling.Brand, &cooling.MaxSpeed, &cooling.CountVentilators, &cooling.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &cooling)
	}

	return messiv, e
}

func (dbCooling dataBase) ShowPartCooling() ([]*models.Cooling, error) {
	rows, e := dbCooling.db.Query(`SELECT coolingid, name, price FROM cooling`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Cooling, 0)
	for rows.Next() {
		var cooling models.Cooling
		e = rows.Scan(&cooling.Coolingid, &cooling.Name, &cooling.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &cooling)
	}

	return messiv, e
}

func (dbCooling dataBase) DeleteCoolingById(id int64) error {
	//_, e := dbCooling.db.Exec(`DELETE FROM cooling WHERE coolingid = $1`, id)

	res, e := dbCooling.db.Exec(`DELETE FROM cooling WHERE coolingid = $1`, id)

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
	//cooling.coolingid, e = q.LastInsertId()

	return nil
}
