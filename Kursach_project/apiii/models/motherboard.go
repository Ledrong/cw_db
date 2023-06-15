package models

type Motherboard struct {
	Motherboardid     int64  `json:"motherboardid,omitempty" db:"column:motherboardid;primaryKey"`
	Name              string `json:"name,omitempty" db:"column:name"`
	Brand             string `json:"brand,omitempty" db:"column:brand"`
	CountSlots        int64  `json:"count_slots,omitempty" db:"column:count_slots"`
	SupportedDdr      string `json:"supported_ddr,omitempty" db:"column:supported_ddr"`
	FormFactorOfRam   string `json:"form_factor_of_ram,omitempty" db:"column:form_factor_of_ram"`
	MinFrequencyOfRam int64  `json:"min_frequency_of_ram,omitempty" db:"column:min_frequency_of_ram"`
	MaxFrequencyOfRam int64  `json:"max_frequency_of_ram,omitempty" db:"column:max_frequency_of_ram"`
	Price             int64  `json:"price,omitempty" db:"column:price"`
}
