package models

type Ram struct {
	Ramid          int64  `json:"ramid,omitempty" db:"column:ramid;primaryKey"`
	Name           string `json:"name,omitempty" db:"column:name"`
	Brand          string `json:"brand,omitempty" db:"column:brand"`
	Rammemory      int64  `json:"rammemory,omitempty" db:"column:rammemory"`
	Ddr            string `json:"ddr,omitempty" db:"column:ddr"`
	FormFactor     string `json:"form_factor,omitempty" db:"column:form_factor"`
	ClockFrequency int64  `json:"clock_frequency,omitempty" db:"column:clock_frequency"`
	Price          int64  `json:"price,omitempty" db:"column:price"`
}
