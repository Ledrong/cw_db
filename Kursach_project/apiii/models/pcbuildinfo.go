package models

type Pcbuildinfo struct {
	Pcbuildid       int64  `db:"pc_buildid"`
	Userlogin       string `db:"users.login"`
	Cpuname         string `db:"cpu.name"`
	Ramname         string `db:"ram.name"`
	Powerboxname    string `db:"powerbox.name"`
	Motherboardname string `db:"motherboard.name"`
	Coolingname     string `db:"cooling.name"`
	Videocardname   string `db:"videocard.name"`
	Compatibility   string `db:"compatibility"`
	Price           int64  `db:"price"`
}
