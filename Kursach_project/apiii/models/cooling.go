package models

type Cooling struct {
	Coolingid        int64  `json:"coolingid,omitempty" db:"column:coolingid;primaryKey"`
	Name             string `json:"name,omitempty" db:"column:name"`
	Brand            string `json:"brand,omitempty" db:"column:brand"`
	MaxSpeed         int64  `json:"max_speed,omitempty" db:"column:max_speed"`
	CountVentilators int64  `json:"count_ventilators,omitempty" db:"column:count_ventilators"`
	Price            int64  `json:"price,omitempty" db:"column:price"`
}
