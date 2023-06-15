package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreatePcbuild(pcbuild *models.Pcbuild, userid int64) error
	UpdatePcbuild(pcbuild *models.Pcbuild, userid int64, pcbuildid int64) error
	SelectPcbuildById(id int64) (*models.Pcbuild, error)   //строка таблицы по факту
	ShowPcbuildById(id int64) (*models.Pcbuildinfo, error) //строка джоинов
	ShowFullPcbuild() ([]*models.Pcbuildinfo, error)
	ShowMyPcbuild(id int64) ([]*models.Pcbuildinfo, error)
	//ShowPartPcbuild(pcbuild *models.Pcbuild) ([]*models.Pcbuild, error)
	DeletePcbuildById(userid int64, id int64) error
	//ShowPcbuild() (*[]models.Pcbuild, error)
}

type dataBase struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbPcbuild dataBase) CreatePcbuild(pcbuild *models.Pcbuild, userid int64) error {
	var id int64
	e := dbPcbuild.db.QueryRow(`INSERT INTO pc_build (user_id, cpu_id, ram_id, powerbox_id, motherboard_id, cooling_id, videocard_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING pc_buildid`, userid, pcbuild.Cpuid, pcbuild.Ramid, pcbuild.Powerboxid, pcbuild.Motherboardid, pcbuild.Coolingid, pcbuild.Videocardid).Scan(&id)
	if e != nil {
		return e
	}

	pcbuild.Pcbuildid = id
	return nil
}

func (dbCpu dataBase) UpdatePcbuild(pcbuild *models.Pcbuild, userid int64, pcbuildid int64) error {
	res, e := dbCpu.db.Exec(`UPDATE pc_build SET cpu_id = $1, ram_id = $2, powerbox_id = $3, motherboard_id = $4, cooling_id = $5, videocard_id = $6 WHERE user_id = $7 AND pc_buildid = $8`, pcbuild.Cpuid, pcbuild.Ramid, pcbuild.Powerboxid, pcbuild.Motherboardid, pcbuild.Coolingid, pcbuild.Videocardid, userid, pcbuildid)
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
	//cpu.Cpuid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	return nil
}

// спросить у Настиииииииии что делать approved
func (dbPcbuild dataBase) SelectPcbuildById(id int64) (*models.Pcbuild, error) {
	pcbuild := models.Pcbuild{}

	q := dbPcbuild.db.QueryRow(`select * from pc_build WHERE pc_buildid = $1`, id).Scan(&pcbuild.Pcbuildid, &pcbuild.Userid, &pcbuild.Cpuid, &pcbuild.Ramid, &pcbuild.Powerboxid, &pcbuild.Motherboardid, &pcbuild.Coolingid, &pcbuild.Videocardid, &pcbuild.Compatibility, &pcbuild.Price)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &pcbuild, nil
}

func (dbPcbuild dataBase) ShowPcbuildById(id int64) (*models.Pcbuildinfo, error) {
	pcbuildinfo := models.Pcbuildinfo{}

	q := dbPcbuild.db.QueryRow(`select pc_buildid, users.login, cpu.name as "Процессор", ram.name as "Оперативная Память", cooling.name as "Охлаждение", powerbox.name as "Блок питания", motherboard.name as "Материнская плата", videocard.name as "Видеокарта", compatibility, pc_build.price as "Стоимость"
		from pc_build
		join users on pc_build.user_id = users.userid
		join cpu on pc_build.cpu_id = cpu.cpuid
		join ram on pc_build.ram_id = ram.ramid
		join cooling on pc_build.cooling_id = cooling.coolingid
		join powerbox on pc_build.powerbox_id = powerbox.powerboxid
		join motherboard on pc_build.motherboard_id = motherboard.motherboardid
		join videocard on pc_build.videocard_id = videocard.videocardid WHERE pc_buildid = $1`, id).Scan(&pcbuildinfo.Pcbuildid, &pcbuildinfo.Userlogin, &pcbuildinfo.Cpuname, &pcbuildinfo.Ramname, &pcbuildinfo.Coolingname, &pcbuildinfo.Powerboxname, &pcbuildinfo.Motherboardname, &pcbuildinfo.Videocardname, &pcbuildinfo.Compatibility, &pcbuildinfo.Price)

	if q != nil {
		return nil, q
	}

	return &pcbuildinfo, nil
}

//func (dbPcbuild dataBase) ShowFullPcbuild() ([]*models.Pcbuild, error) {
//	rows, e := dbPcbuild.db.Query(`SELECT * FROM pc_build`)
//	if e != nil {
//		//fmt.Println(e.Error())
//		return nil, e
//	}
//	defer rows.Close()
//
//	messiv := make([]*models.Pcbuild, 0)
//	for rows.Next() {
//		var pcbuild models.Pcbuild
//		e = rows.Scan(&pcbuild.Pcbuildid, &pcbuild.CpuId, &pcbuild.RamId, &pcbuild.PowerboxId, &pcbuild.Motherboardid, &pcbuild.Coolingid, &pcbuild.Videocardid, &pcbuild.Compatibility, &pcbuild.Price)
//		if e != nil {
//			return nil, e
//		}
//		messiv = append(messiv, &pcbuild)
//	}
//
//	return messiv, e
//}

func (dbPcbuild dataBase) ShowFullPcbuild() ([]*models.Pcbuildinfo, error) {
	rows, e := dbPcbuild.db.Query(`select pc_buildid, users.login, cpu.name as "Процессор", ram.name as "Оперативная Память", cooling.name as "Охлаждение", powerbox.name as "Блок питания", motherboard.name as "Материнская плата", videocard.name as "Видеокарта", compatibility, pc_build.price as "Стоимость"
from pc_build
join users on pc_build.user_id = users.userid
join cpu on pc_build.cpu_id = cpu.cpuid
join ram on pc_build.ram_id = ram.ramid
join cooling on pc_build.cooling_id = cooling.coolingid
join powerbox on pc_build.powerbox_id = powerbox.powerboxid
join motherboard on pc_build.motherboard_id = motherboard.motherboardid
join videocard on pc_build.videocard_id = videocard.videocardid`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Pcbuildinfo, 0)
	for rows.Next() {
		var pcbuildinfo models.Pcbuildinfo
		e = rows.Scan(&pcbuildinfo.Pcbuildid, &pcbuildinfo.Userlogin, &pcbuildinfo.Cpuname, &pcbuildinfo.Ramname, &pcbuildinfo.Coolingname, &pcbuildinfo.Powerboxname, &pcbuildinfo.Motherboardname, &pcbuildinfo.Videocardname, &pcbuildinfo.Compatibility, &pcbuildinfo.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &pcbuildinfo)
	}

	return messiv, e
}

func (dbPcbuild dataBase) ShowMyPcbuild(id int64) ([]*models.Pcbuildinfo, error) {
	rows, e := dbPcbuild.db.Query(`select pc_buildid, users.login, cpu.name as "Процессор", ram.name as "Оперативная Память", cooling.name as "Охлаждение", powerbox.name as "Блок питания", motherboard.name as "Материнская плата", videocard.name as "Видеокарта", compatibility, pc_build.price as "Стоимость"
from pc_build
join users on pc_build.user_id = users.userid
join cpu on pc_build.cpu_id = cpu.cpuid
join ram on pc_build.ram_id = ram.ramid
join cooling on pc_build.cooling_id = cooling.coolingid
join powerbox on pc_build.powerbox_id = powerbox.powerboxid
join motherboard on pc_build.motherboard_id = motherboard.motherboardid
join videocard on pc_build.videocard_id = videocard.videocardid WHERE users.userid = $1`, id)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.Pcbuildinfo, 0)
	for rows.Next() {
		var pcbuildinfo models.Pcbuildinfo
		e = rows.Scan(&pcbuildinfo.Pcbuildid, &pcbuildinfo.Userlogin, &pcbuildinfo.Cpuname, &pcbuildinfo.Ramname, &pcbuildinfo.Coolingname, &pcbuildinfo.Powerboxname, &pcbuildinfo.Motherboardname, &pcbuildinfo.Videocardname, &pcbuildinfo.Compatibility, &pcbuildinfo.Price)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &pcbuildinfo)
	}

	return messiv, e
}

func (dbPcbuild dataBase) DeletePcbuildById(userid int64, id int64) error {
	//_, e := dbPcbuild.db.Exec(`DELETE FROM pc_build WHERE pc_buildid = $1`, id)
	res, e := dbPcbuild.db.Exec(`DELETE FROM pc_build WHERE user_id = $1 AND pc_buildid = $2`, userid, id)

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

	return nil
}
