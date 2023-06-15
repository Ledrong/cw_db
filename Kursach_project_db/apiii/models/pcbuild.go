package models

type Pcbuild struct {
	Pcbuildid     int64  `json:"pc_buildid,omitempty" db:"column:pc_buildid;primaryKey"`
	Userid        int64  `json:"user_id,omitempty" db:"column:user_id"`
	Cpuid         int64  `json:"cpu_id,omitempty" db:"column:cpu_id"`
	Ramid         int64  `json:"ram_id,omitempty" db:"column:ram_id"`
	Powerboxid    int64  `json:"powerbox_id,omitempty" db:"column:powerbox_id"`
	Motherboardid int64  `json:"motherboard_id,omitempty" db:"column:motherboard_id"`
	Coolingid     int64  `json:"cooling_id,omitempty" db:"column:cooling_id"`
	Videocardid   int64  `json:"videocard_id,omitempty" db:"column:videocard_id"`
	Compatibility string `json:"compatibility,omitempty" db:"column:compatibility"`
	Price         int64  `json:"price,omitempty" db:"column:price"`
}
