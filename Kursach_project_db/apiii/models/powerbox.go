package models

type Powerbox struct {
	Powerboxid int64  `json:"powerboxid,omitempty" db:"column:powerboxid;primaryKey"`
	Name       string `json:"name,omitempty" db:"column:name"`
	Brand      string `json:"brand,omitempty" db:"column:brand"`
	Power      int64  `json:"power,omitempty" db:"column:power"`
	FormFactor string `json:"form_factor,omitempty" db:"column:form_factor"`
	Price      int64  `json:"price,omitempty" db:"column:price"`
}
