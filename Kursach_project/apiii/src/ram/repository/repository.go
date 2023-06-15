package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateRam(ram *models.Ram) error
	UpdateRam(ram *models.Ram, id int64) error
	SelectRamById(id int64) (*models.Ram, error)
	ShowFullRam() ([]*models.Ram, error)
	ShowCompatibilityCpu(id int64) ([]*models.Cpu, error)
	ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error)
	ShowPartRam() ([]*models.Ram, error)
	DeleteRamById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbRam dataBase) CreateRam(ram *models.Ram) error {
	var id int64
	e := dbRam.db.QueryRow(`INSERT INTO ram (name, brand, rammemory, ddr, form_factor, clock_frequency, price) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ramid`, ram.Name, ram.Brand, ram.Rammemory, ram.Ddr, ram.FormFactor, ram.ClockFrequency, ram.Price).Scan(&id)
	if e != nil {
		return e
	}

	ram.Ramid = id

	_, err := dbRam.db.Exec(`INSERT INTO motherboard_ram(motherboard_id, ram_id) SELECT motherboard.motherboardid, ram.ramid FROM motherboard, ram WHERE ram.ramid = $1 AND position(ram.ddr in motherboard.supported_ddr) > 0 AND ram.clock_frequency >= motherboard.min_frequency_of_ram AND ram.clock_frequency <= motherboard.max_frequency_of_ram AND ram.form_factor = motherboard.form_factor_of_ram`, id)
	if err != nil {
		return nil
	}

	_, err = dbRam.db.Exec(`INSERT INTO cpu_ram(cpu_id, ram_id) SELECT cpu.cpuid, ram.ramid FROM cpu, ram WHERE ram.ramid = $1 AND position(ram.ddr in cpu.supported_ddr) > 0 AND ram.clock_frequency >= cpu.min_frequency_of_ram AND ram.clock_frequency <= cpu.max_frequency_of_ram`, id)
	if err != nil {
		return nil
	}

	return nil
}

func (dbRam dataBase) UpdateRam(ram *models.Ram, id int64) error {
	_, e := dbRam.db.Exec(`DELETE FROM motherboard_ram WHERE ram_id = $1`, id)
	if e != nil {
		return nil
	}

	_, e = dbRam.db.Exec(`DELETE FROM cpu_ram WHERE ram_id = $1`, id)
	if e != nil {
		return nil
	}
	res, e := dbRam.db.Exec(`UPDATE ram SET  name = $1, brand = $2, rammemory = $3, ddr = $4, form_factor = $5, clock_frequency = $6, price = $7 WHERE ramid = $8`, ram.Name, ram.Brand, ram.Rammemory, ram.Ddr, ram.FormFactor, ram.ClockFrequency, ram.Price, id)
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
	//ram.Ramid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	_, err := dbRam.db.Exec(`INSERT INTO motherboard_ram(motherboard_id, ram_id) SELECT motherboard.motherboardid, ram.ramid FROM motherboard, ram WHERE ram.ramid = $1 AND position(ram.ddr in motherboard.supported_ddr) > 0 AND ram.clock_frequency >= motherboard.min_frequency_of_ram AND ram.clock_frequency <= motherboard.max_frequency_of_ram AND ram.form_factor = motherboard.form_factor_of_ram`, id)
	if err != nil {
		return nil
	}

	_, err = dbRam.db.Exec(`INSERT INTO cpu_ram(cpu_id, ram_id) SELECT cpu.cpuid, ram.ramid FROM cpu, ram WHERE ram.ramid = $1 AND position(ram.ddr in cpu.supported_ddr) > 0 AND ram.clock_frequency >= cpu.min_frequency_of_ram AND ram.clock_frequency <= cpu.max_frequency_of_ram`, id)
	if err != nil {
		return nil
	}

	return nil
}

// спросить у Настиииииииии что делать approved
func (dbRam dataBase) SelectRamById(id int64) (*models.Ram, error) {
	ram := models.Ram{}

	q := dbRam.db.QueryRow(`SELECT * FROM ram WHERE ramid = $1`, id).Scan(&ram.Ramid, &ram.Name, &ram.Brand, &ram.Rammemory, &ram.Ddr, &ram.FormFactor, &ram.ClockFrequency, &ram.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &ram, nil
}

func (dbRam dataBase) ShowFullRam() ([]*models.Ram, error) {
	rows, e := dbRam.db.Query(`SELECT * FROM ram`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Ram, 0)
	for rows.Next() {
		var ram models.Ram
		e = rows.Scan(&ram.Ramid, &ram.Name, &ram.Brand, &ram.Rammemory, &ram.Ddr, &ram.FormFactor, &ram.ClockFrequency, &ram.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &ram)
	}

	return messiv, e
}

func (dbRam dataBase) ShowCompatibilityCpu(id int64) ([]*models.Cpu, error) {
	//_, e := dbCpu.db.Exec(`INSERT INTO cpu_ram (cpu_id, ram_id) SELECT cpuid, ramid FROM cpu, ram ON CONFLICT DO NOTHING`)
	//if e != nil {
	//	return nil, e
	//}
	//
	//rows, e := dbCpu.db.Query(`SELECT ram.ramid, ram.name, ram.brand, ram.rammemory, ram.ddr, ram.price FROM ram JOIN cpu_ram ON ram.ramid = cpu_ram.ram_id JOIN cpu ON cpu.cpuid = cpu_ram.cpu_id WHERE (position(ddr in supported_ddr) > 0) AND cpu.cpuid = $1`, id)
	//if e != nil {
	//	//fmt.Println(e.Error())
	//	return nil, e
	//}
	//defer rows.Close()
	//
	//messiv := make([]*models.Ram, 0)
	//for rows.Next() {
	//	var ram models.Ram
	//	e = rows.Scan(&ram.Ramid, &ram.Name, &ram.Brand, &ram.Rammemory, &ram.Ddr, &ram.Price)
	//	if e != nil {
	//		return nil, e
	//	}
	//	messiv = append(messiv, &ram)
	//}
	//
	//_, e = dbCpu.db.Exec(`DELETE FROM cpu_ram`)
	//if e != nil {
	//	return nil, e
	//}
	//
	//return messiv, e
	rows, e := dbRam.db.Query(`SELECT cpu.cpuid, cpu.name, cpu.brand, cpu.series, cpu.model, cpu.supported_ddr, cpu.price FROM cpu JOIN cpu_ram ON cpu.cpuid = cpu_ram.cpu_id JOIN ram ON ram.ramid = cpu_ram.ram_id WHERE ram.ramid = $1`, id)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Cpu, 0)
	for rows.Next() {
		var cpu models.Cpu
		e = rows.Scan(&cpu.Cpuid, &cpu.Name, &cpu.Brand, &cpu.Series, &cpu.Model, &cpu.SupportedDdr, &cpu.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &cpu)
	}

	return messiv, e
}

func (dbRam dataBase) ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error) {
	//_, e := dbCpu.db.Exec(`INSERT INTO cpu_ram (cpu_id, ram_id) SELECT cpuid, ramid FROM cpu, ram ON CONFLICT DO NOTHING`)
	//if e != nil {
	//	return nil, e
	//}
	//
	//rows, e := dbCpu.db.Query(`SELECT ram.ramid, ram.name, ram.brand, ram.rammemory, ram.ddr, ram.price FROM ram JOIN cpu_ram ON ram.ramid = cpu_ram.ram_id JOIN cpu ON cpu.cpuid = cpu_ram.cpu_id WHERE (position(ddr in supported_ddr) > 0) AND cpu.cpuid = $1`, id)
	//if e != nil {
	//	//fmt.Println(e.Error())
	//	return nil, e
	//}
	//defer rows.Close()
	//
	//messiv := make([]*models.Ram, 0)
	//for rows.Next() {
	//	var ram models.Ram
	//	e = rows.Scan(&ram.Ramid, &ram.Name, &ram.Brand, &ram.Rammemory, &ram.Ddr, &ram.Price)
	//	if e != nil {
	//		return nil, e
	//	}
	//	messiv = append(messiv, &ram)
	//}
	//
	//_, e = dbCpu.db.Exec(`DELETE FROM cpu_ram`)
	//if e != nil {
	//	return nil, e
	//}
	//
	//return messiv, e
	rows, e := dbRam.db.Query(`SELECT motherboard.motherboardid, motherboard.name, motherboard.brand, motherboard.count_slots, motherboard.supported_ddr, motherboard.form_factor_of_ram, motherboard.min_frequency_of_ram, motherboard.max_frequency_of_ram, motherboard.price FROM motherboard JOIN motherboard_ram ON motherboard.motherboardid = motherboard_ram.motherboard_id JOIN ram ON motherboard_ram.ram_id = ram.ramid WHERE ram.ramid = $1`, id)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Motherboard, 0)
	for rows.Next() {
		var motherboard models.Motherboard
		e = rows.Scan(&motherboard.Motherboardid, &motherboard.Name, &motherboard.Brand, &motherboard.CountSlots, &motherboard.SupportedDdr, &motherboard.FormFactorOfRam, &motherboard.MinFrequencyOfRam, &motherboard.MaxFrequencyOfRam, &motherboard.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &motherboard)
	}

	return messiv, e
}

func (dbRam dataBase) ShowPartRam() ([]*models.Ram, error) {
	rows, e := dbRam.db.Query(`SELECT ramid, name, price FROM ram`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Ram, 0)
	for rows.Next() {
		var ram models.Ram
		e = rows.Scan(&ram.Ramid, &ram.Name, &ram.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &ram)
	}

	return messiv, e
}

// /спросить у Насти
func (dbRam dataBase) DeleteRamById(id int64) error {
	_, e := dbRam.db.Exec(`DELETE FROM motherboard_ram WHERE ram_id = $1`, id)
	if e != nil {
		return nil
	}

	_, e = dbRam.db.Exec(`DELETE FROM cpu_ram WHERE ram_id = $1`, id)
	if e != nil {
		return nil
	}

	res, e := dbRam.db.Exec(`DELETE FROM ram WHERE ramid = $1`, id)

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
	//ram.Ramid, e = q.LastInsertId()

	return nil
}
