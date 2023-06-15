package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateMotherboard(motherboard *models.Motherboard) error
	UpdateMotherboard(motherboard *models.Motherboard, id int64) error
	SelectMotherboardById(id int64) (*models.Motherboard, error)
	ShowFullMotherboard() ([]*models.Motherboard, error)
	ShowCompatibilityRam(id int64) ([]*models.Ram, error)
	ShowCompatibilityCpu(id int64) ([]*models.Cpu, error)
	ShowPartMotherboard() ([]*models.Motherboard, error)
	DeleteMotherboardById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbMotherboard dataBase) CreateMotherboard(motherboard *models.Motherboard) error {
	var id int64
	e := dbMotherboard.db.QueryRow(`INSERT INTO motherboard (name, brand, count_slots, supported_ddr, form_factor_of_ram, min_frequency_of_ram, max_frequency_of_ram, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING motherboardid`, motherboard.Name, motherboard.Brand, motherboard.CountSlots, motherboard.SupportedDdr, motherboard.FormFactorOfRam, motherboard.MinFrequencyOfRam, motherboard.MaxFrequencyOfRam, motherboard.Price).Scan(&id)
	if e != nil {
		return e
	}
	motherboard.Motherboardid = id

	_, err := dbMotherboard.db.Exec(`INSERT INTO motherboard_ram(motherboard_id, ram_id) SELECT motherboard.motherboardid, ram.ramid FROM motherboard, ram WHERE motherboard.motherboardid = $1 AND position(ram.ddr in motherboard.supported_ddr) > 0 AND ram.clock_frequency >= motherboard.min_frequency_of_ram AND ram.clock_frequency <= motherboard.max_frequency_of_ram AND ram.form_factor = motherboard.form_factor_of_ram`, id)
	if err != nil {
		return nil
	}

	_, err = dbMotherboard.db.Exec(`INSERT INTO cpu_motherboard(cpu_id, motherboard_id) SELECT cpu.cpuid, motherboard.motherboardid FROM cpu, motherboard WHERE motherboard.motherboardid = $1 AND (string_to_array($2, ', ') && string_to_array(cpu.supported_ddr, ', ')) AND motherboard.max_frequency_of_ram >= cpu.min_frequency_of_ram AND motherboard.min_frequency_of_ram <= cpu.max_frequency_of_ram`, id, motherboard.SupportedDdr)
	if err != nil {
		return nil
	}

	return nil
}

func (dbMotherboard dataBase) UpdateMotherboard(motherboard *models.Motherboard, id int64) error {
	_, e := dbMotherboard.db.Exec(`DELETE FROM motherboard_ram WHERE motherboard_id = $1`, id)
	if e != nil {
		return nil
	}

	_, e = dbMotherboard.db.Exec(`DELETE FROM cpu_motherboard WHERE motherboard_id = $1`, id)
	if e != nil {
		return nil
	}

	res, e := dbMotherboard.db.Exec(`UPDATE motherboard SET  name = $1, brand = $2, count_slots = $3, supported_ddr = $4, form_factor_of_ram = $5, min_frequency_of_ram = $6, max_frequency_of_ram = $7, price = $8 WHERE  motherboardid = $9`, motherboard.Name, motherboard.Brand, motherboard.CountSlots, motherboard.SupportedDdr, motherboard.FormFactorOfRam, motherboard.MinFrequencyOfRam, motherboard.MaxFrequencyOfRam, motherboard.Price, id)
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
	//motherboard.motherboardid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	_, err := dbMotherboard.db.Exec(`INSERT INTO motherboard_ram(motherboard_id, ram_id) SELECT motherboard.motherboardid, ram.ramid FROM motherboard, ram WHERE motherboard.motherboardid = $1 AND position(ram.ddr in motherboard.supported_ddr) > 0 AND ram.clock_frequency >= motherboard.min_frequency_of_ram AND ram.clock_frequency <= motherboard.max_frequency_of_ram AND ram.form_factor = motherboard.form_factor_of_ram`, id)
	if err != nil {
		return nil
	}

	_, err = dbMotherboard.db.Exec(`INSERT INTO cpu_motherboard(cpu_id, motherboard_id) SELECT cpu.cpuid, motherboard.motherboardid FROM cpu, motherboard WHERE motherboard.motherboardid = $1 AND (string_to_array($2, ', ') && string_to_array(cpu.supported_ddr, ', ')) AND motherboard.max_frequency_of_ram >= cpu.min_frequency_of_ram AND motherboard.min_frequency_of_ram <= cpu.max_frequency_of_ram`, id, motherboard.SupportedDdr)
	if err != nil {
		return nil
	}

	return nil
}

// спросить у Настиииииииии что делать approved
func (dbMotherboard dataBase) SelectMotherboardById(id int64) (*models.Motherboard, error) {
	motherboard := models.Motherboard{}
	//надо ли сканить в айди (надо)
	q := dbMotherboard.db.QueryRow(`SELECT * FROM motherboard WHERE motherboardid = $1`, id).Scan(&motherboard.Motherboardid, &motherboard.Name, &motherboard.Brand, &motherboard.CountSlots, &motherboard.SupportedDdr, &motherboard.FormFactorOfRam, &motherboard.MinFrequencyOfRam, &motherboard.MaxFrequencyOfRam, &motherboard.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &motherboard, nil
}

func (dbMotherboard dataBase) ShowFullMotherboard() ([]*models.Motherboard, error) {
	rows, e := dbMotherboard.db.Query(`SELECT * FROM motherboard`)
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

func (dbMotherboard dataBase) ShowCompatibilityRam(id int64) ([]*models.Ram, error) {
	rows, e := dbMotherboard.db.Query(`SELECT ram.ramid, ram.name, ram.brand, ram.rammemory, ram.ddr, ram.form_factor, ram.clock_frequency, ram.price FROM ram JOIN motherboard_ram ON ram.ramid = motherboard_ram.ram_id JOIN motherboard ON motherboard_ram.motherboard_id = motherboard.motherboardid WHERE motherboard.motherboardid = $1`, id)
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

func (dbMotherboard dataBase) ShowCompatibilityCpu(id int64) ([]*models.Cpu, error) {
	rows, e := dbMotherboard.db.Query(`SELECT cpu.cpuid,cpu.name, cpu.brand, cpu.series, cpu.model, cpu.basic_frequency, cpu.supported_ddr, cpu.min_frequency_of_ram, cpu.max_frequency_of_ram, cpu.price FROM cpu JOIN cpu_ram ON cpu.cpuid = cpu_ram.cpu_id JOIN ram ON cpu_ram.ram_id = ram.ramid WHERE cpu.cpuid = $1`, id)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Cpu, 0)
	for rows.Next() {
		var cpu models.Cpu
		e = rows.Scan(&cpu.Cpuid, &cpu.Name, &cpu.Brand, &cpu.Series, &cpu.Model, &cpu.BasicFrequency, &cpu.SupportedDdr, &cpu.MinFrequencyOfRam, &cpu.MaxFrequencyOfRam, &cpu.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &cpu)
	}

	return messiv, e
}

func (dbMotherboard dataBase) ShowPartMotherboard() ([]*models.Motherboard, error) {
	rows, e := dbMotherboard.db.Query(`SELECT motherboardid, name, price FROM motherboard`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Motherboard, 0)
	for rows.Next() {
		var motherboard models.Motherboard
		e = rows.Scan(&motherboard.Motherboardid, &motherboard.Name, &motherboard.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &motherboard)
	}

	return messiv, e
}

// /спросить у Насти approved
func (dbMotherboard dataBase) DeleteMotherboardById(id int64) error {
	_, e := dbMotherboard.db.Exec(`DELETE FROM motherboard_ram WHERE motherboard_id = $1`, id)
	if e != nil {
		return nil
	}

	_, e = dbMotherboard.db.Exec(`DELETE FROM cpu_motherboard WHERE motherboard_id = $1`, id)
	if e != nil {
		return nil
	}

	res, e := dbMotherboard.db.Exec(`DELETE FROM motherboard WHERE motherboardid = $1`, id)

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

	//motherboard.Motherboardid, e = q.LastInsertId()

	return nil
}
