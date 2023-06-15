package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateCpu(cpu *models.Cpu) error
	UpdateCpu(cpu *models.Cpu, id int64) error
	SelectCpuById(id int64) (*models.Cpu, error)
	ShowFullCpu() ([]*models.Cpu, error)
	ShowCompatibilityRam(id int64) ([]*models.Ram, error)
	ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error)
	ShowPartCpu() ([]*models.Cpu, error)
	DeleteCpuById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

//func (dbCpu dataBase) CreateCpu(cpu *models.Cpu) error {
//	q, e := dbCpu.db.Exec(`INSERT INTO cpu (name, brand, series, model, supported_ddr, price) VALUES ($1, $2, $3, $4, $5, $6)`, cpu.Name, cpu.Brand, cpu.Series, cpu.Model, cpu.SupportedDdr, cpu.Price)
//	if e != nil {
//		return e
//	}
//	cpu.Cpuid, e = q.LastInsertId() //драйвер не поддерживает, заменить на queryrow
//
//	if e != nil {
//		return e
//	}
//	return nil
//}

func (dbCpu dataBase) CreateCpu(cpu *models.Cpu) error {
	var id int64
	e := dbCpu.db.QueryRow(`INSERT INTO cpu (name, brand, series, model, basic_frequency, supported_ddr, min_frequency_of_ram, max_frequency_of_ram, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING cpuid`, cpu.Name, cpu.Brand, cpu.Series, cpu.Model, cpu.BasicFrequency, cpu.SupportedDdr, cpu.MinFrequencyOfRam, cpu.MaxFrequencyOfRam, cpu.Price).Scan(&id)
	if e != nil {
		return e
	}
	cpu.Cpuid = id
	//cpu.Cpuid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	_, err := dbCpu.db.Exec(`INSERT INTO cpu_ram(cpu_id, ram_id) SELECT cpu.cpuid, ram.ramid FROM cpu, ram WHERE cpu.cpuid = $1 AND position(ram.ddr in cpu.supported_ddr) > 0 AND ram.clock_frequency >= cpu.min_frequency_of_ram AND ram.clock_frequency <= cpu.max_frequency_of_ram`, id)
	if err != nil {
		return nil
	}

	_, err = dbCpu.db.Exec(`INSERT INTO cpu_motherboard(cpu_id, motherboard_id) SELECT cpu.cpuid, motherboard.motherboardid FROM cpu, motherboard WHERE cpu.cpuid = $1 AND (string_to_array($2, ', ') && string_to_array(motherboard.supported_ddr, ', ')) AND cpu.max_frequency_of_ram >= motherboard.min_frequency_of_ram AND cpu.min_frequency_of_ram <= motherboard.max_frequency_of_ram`, id, cpu.SupportedDdr)
	if err != nil {
		return nil
	}

	return nil
}

func (dbCpu dataBase) UpdateCpu(cpu *models.Cpu, id int64) error {
	_, e := dbCpu.db.Exec(`DELETE FROM cpu_ram WHERE cpu_id = $1`, id)
	if e != nil {
		return nil
	}

	_, e = dbCpu.db.Exec(`DELETE FROM cpu_motherboard WHERE cpu_id = $1`, id)
	if e != nil {
		return nil
	}

	res, e := dbCpu.db.Exec(`UPDATE cpu SET  name = $1, brand = $2, series = $3, model = $4, basic_frequency = $5, supported_ddr = $6, min_frequency_of_ram = $7, max_frequency_of_ram = $8, price = $9 WHERE cpuid = $10`, cpu.Name, cpu.Brand, cpu.Series, cpu.Model, cpu.BasicFrequency, cpu.SupportedDdr, cpu.MinFrequencyOfRam, cpu.MaxFrequencyOfRam, cpu.Price, id)
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

	_, err := dbCpu.db.Exec(`INSERT INTO cpu_ram(cpu_id, ram_id) SELECT cpu.cpuid, ram.ramid FROM cpu, ram WHERE cpu.cpuid = $1 AND position(ram.ddr in cpu.supported_ddr) > 0 AND ram.clock_frequency >= cpu.min_frequency_of_ram AND ram.clock_frequency <= cpu.max_frequency_of_ram`, id)
	if err != nil {
		return nil
	}

	_, err = dbCpu.db.Exec(`INSERT INTO cpu_motherboard(cpu_id, motherboard_id) SELECT cpu.cpuid, motherboard.motherboardid FROM cpu, motherboard WHERE cpu.cpuid = $1 AND (string_to_array($2, ', ') && string_to_array(motherboard.supported_ddr, ', ')) AND cpu.max_frequency_of_ram >= motherboard.min_frequency_of_ram AND cpu.min_frequency_of_ram <= motherboard.max_frequency_of_ram`, id, cpu.SupportedDdr)
	if err != nil {
		return nil
	}
	//cpu.Cpuid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	return nil
}

func (dbCpu dataBase) SelectCpuById(id int64) (*models.Cpu, error) {
	cpu := models.Cpu{}

	q := dbCpu.db.QueryRow(`SELECT * FROM cpu WHERE cpuid = $1`, id).Scan(&cpu.Cpuid, &cpu.Name, &cpu.Brand, &cpu.Series, &cpu.Model, &cpu.BasicFrequency, &cpu.SupportedDdr, &cpu.MinFrequencyOfRam, &cpu.MaxFrequencyOfRam, &cpu.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &cpu, nil
}

func (dbCpu dataBase) ShowFullCpu() ([]*models.Cpu, error) {
	rows, e := dbCpu.db.Query(`SELECT * FROM cpu`)
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

func (dbCpu dataBase) ShowCompatibilityRam(id int64) ([]*models.Ram, error) {
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
	rows, e := dbCpu.db.Query(`SELECT ram.ramid, ram.name, ram.brand, ram.rammemory, ram.ddr, ram.form_factor, ram.clock_frequency, ram.price FROM ram JOIN cpu_ram ON ram.ramid = cpu_ram.ram_id JOIN cpu ON cpu_ram.cpu_id = cpu.cpuid WHERE cpu.cpuid = $1`, id)
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

func (dbCpu dataBase) ShowCompatibilityMotherboard(id int64) ([]*models.Motherboard, error) {
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
	rows, e := dbCpu.db.Query(`SELECT motherboard.motherboardid, motherboard.name, motherboard.brand, motherboard.count_slots, motherboard.supported_ddr, motherboard.form_factor_of_ram, motherboard.min_frequency_of_ram, motherboard.max_frequency_of_ram, motherboard.price FROM motherboard JOIN cpu_motherboard ON motherboard.motherboardid = cpu_motherboard.motherboard_id JOIN cpu ON cpu_motherboard.cpu_id = cpu.cpuid WHERE cpu.cpuid = $1`, id)
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

func (dbCpu dataBase) ShowPartCpu() ([]*models.Cpu, error) {
	rows, e := dbCpu.db.Query(`SELECT cpuid, name, price FROM cpu`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Cpu, 0)
	for rows.Next() {
		var cpu models.Cpu
		e = rows.Scan(&cpu.Cpuid, &cpu.Name, &cpu.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &cpu)
	}

	return messiv, e
}

// /спросить у Насти надо ли передавать что-то кроме id
func (dbCpu dataBase) DeleteCpuById(id int64) error {
	_, e := dbCpu.db.Exec(`DELETE FROM cpu_ram WHERE cpu_id = $1`, id)
	if e != nil {
		return nil
	}

	_, e = dbCpu.db.Exec(`DELETE FROM cpu_motherboard WHERE cpu_id = $1`, id)
	if e != nil {
		return nil
	}

	//res, e := dbCpu.db.Exec(`DELETE FROM cpu WHERE cpuid = $1`, id)
	res, e := dbCpu.db.Exec(`DELETE FROM cpu WHERE cpuid = $1`, id)
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

	//cpu.Cpuid, e = q.LastInsertId()

	return nil
}

//func (dbForum *dataBase) CreateForum(forum *models.Forum) error {
//	tx := dbForum.db.Omit("posts", "threads").Create(forum)
//	if tx.Error != nil {
//		return errors.Wrap(tx.Error, "database error (table forums)")
//	}
//
//	return nil
//}
//
//func (dbForum *dataBase) SelectForumBySlug(slug string) (*models.Forum, error) {
//	forum := models.Forum{}
//
//	tx := dbForum.db.Where("slug = ?", slug).Take(&forum)
//	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
//		return nil, models.ErrNotFound
//	} else if tx.Error != nil {
//		return nil, errors.Wrap(tx.Error, "database error (table forums)")
//	}
//
//	return &forum, nil
//}
//
//func (dbForum *dataBase) SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error) {
//	users := make([]*models.User, 0, 10)
//
//	if desc {
//		if since != "" {
//			tx := dbForum.db.Limit(limit).Where("nickname IN (?)", dbForum.db.
//				Select("user_nickname").Table("forum_user").Where("forum = ? AND user_nickname < ?",
//				slug, since)).Order("nickname desc").Find(&users)
//			if tx.Error != nil {
//				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
//			}
//		} else {
//			tx := dbForum.db.Limit(limit).Where("nickname IN (?)", dbForum.db.
//				Select("user_nickname").Table("forum_user").Where("forum = ?", slug)).
//				Order("nickname desc").Find(&users)
//			if tx.Error != nil {
//				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
//			}
//		}
//	} else {
//		if since != "" {
//			tx := dbForum.db.Limit(limit).Where("nickname IN (?)", dbForum.db.
//				Select("user_nickname").Table("forum_user").Where("forum = ? AND user_nickname > ?",
//				slug, since)).Order("nickname").Find(&users)
//			if tx.Error != nil {
//				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
//			}
//		} else {
//			tx := dbForum.db.Limit(limit).Where("nickname IN (?)", dbForum.db.
//				Select("user_nickname").Table("forum_user").Where("forum = ?", slug)).
//				Order("nickname").Find(&users)
//			if tx.Error != nil {
//				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
//			}
//		}
//	}
//
//	return users, nil
//}
//
//func (dbForum *dataBase) CreateForumUser(forum string, user string) error {
//	fu := models.ForumUser{
//		Forum: forum,
//		User:  user,
//	}
//	tx := dbForum.db.Table("forum_user").Clauses(clause.OnConflict{DoNothing: true}).Create(&fu)
//	if tx.Error != nil {
//		return errors.Wrap(tx.Error, "database error (table forum_user)")
//	}
//
//	return nil
//}
