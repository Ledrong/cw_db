package repository

import (
	"Kursach_project/apiii/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateUser(user *models.User) error
	Login(user *models.User) error
	UpdateUser(user *models.User, id int64) error
	GetroleUserById(id int64) (string, error)
	SelectUserById(id int64) (*models.User, error)
	ShowFullUser() ([]*models.User, error)
	DeleteUserById(id int64) error
}

type dataBase struct {
	db *sqlx.DB
}

func (dbUser dataBase) GetroleUserById(id int64) (string, error) {
	var userrole string
	q := dbUser.db.QueryRow(`SELECT role FROM users WHERE userid = $1`, id).Scan(&userrole)

	if q != nil {
		if q == sql.ErrNoRows {
			return "", models.ErrNotFound
		} else {
			return "", q
		}
	}

	return userrole, nil
}

func New(db *sqlx.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbUser dataBase) CreateUser(user *models.User) error {
	var id int64
	_, err := dbUser.db.Exec(`CREATE USER $1 LOGIN PASSWORD $2 IN ROLE "user"`, user.Login, user.Password)
	if err != nil {
		return err
	}
	//e := dbUser.db.QueryRow(`INSERT INTO users (login, password, role) VALUES ($1, $2, $3) RETURNING userid`, user.Login, user.Password, user.Role).Scan(&id)
	userrole := "user"
	e := dbUser.db.QueryRow(`INSERT INTO users (login, password, role) VALUES ($1, $2, $3) RETURNING userid`, user.Login, user.Password, userrole).Scan(&id)
	if e != nil {
		return e
	}

	user.Userid = id
	user.Role = "user"

	return nil
}

func (dbUser dataBase) Login(user *models.User) error {
	e := dbUser.db.QueryRow(`SELECT userid, role FROM users WHERE login=$1 AND password=$2`, user.Login, user.Password).Scan(&user.Userid, &user.Role)
	if e != nil {
		if e == sql.ErrNoRows {
			return models.ErrNotFound
		} else {
			return e
		}
	}

	return nil
}

func (dbUser dataBase) UpdateUser(user *models.User, id int64) error {
	old_name, err := dbUser.db.Exec(`SELECT login FROM users WHERE users.userid = $1`, user.Userid)
	if err != nil {
		return err
	}

	_, err = dbUser.db.Exec(`ALTER USER $1 RENAME TO $2 WITH PASSWORD $3 SET ROLE new_role`, old_name, user.Login, user.Password, user.Role)

	res, e := dbUser.db.Exec(`UPDATE users SET login = $1, password = $2, role = $3 WHERE userid = $4`, user.Login, user.Password, user.Role, id)
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
	//user.Userid, e = q.LastInsertId() драйвер не поддерживает, заменить на queryrow

	return nil
}

func (dbUser dataBase) SelectUserById(id int64) (*models.User, error) {
	user := models.User{}

	q := dbUser.db.QueryRow(`SELECT * FROM users WHERE userid = $1`, id).Scan(&user.Userid, &user.Login, &user.Password, &user.Role)

	if q != nil {
		if q == sql.ErrNoRows {
			return nil, models.ErrNotFound
		} else {
			return nil, q
		}
	}

	return &user, nil
}

func (dbUser dataBase) ShowFullUser() ([]*models.User, error) {
	rows, e := dbUser.db.Query(`SELECT * FROM users`)
	if e != nil {
		//fmt.Println(e.Error())
		return nil, e
	}
	defer rows.Close()

	messiv := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		e = rows.Scan(&user.Userid, &user.Login, &user.Password, &user.Role)
		if e != nil {
			return nil, e
		}
		messiv = append(messiv, &user)
	}

	return messiv, e
}

func (dbUser dataBase) DeleteUserById(id int64) error {
	//_, e := dbUser.db.Exec(`DELETE FROM "user" WHERE userid = $1`, id)
	login, err := dbUser.db.Exec(`SELECT login FROM users WHERE userid = $1`, id)
	if err != nil {
		return err
	}

	_, err = dbUser.db.Exec(`DROP USER $1`, login)
	res, e := dbUser.db.Exec(`DELETE FROM users WHERE userid = $1`, id)

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
