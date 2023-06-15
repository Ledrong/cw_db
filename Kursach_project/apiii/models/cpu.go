package models

type Cpu struct {
	Cpuid             int64   `json:"cpuid,omitempty" db:"cpuid"`
	Name              string  `json:"name,omitempty" db:"name"`
	Brand             string  `json:"brand,omitempty" db:"brand"`
	Series            string  `json:"series,omitempty" db:"series"`
	Model             string  `json:"model,omitempty" db:"model"`
	BasicFrequency    float64 `json:"basic_frequency,omitempty" db:"basic_frequency"`
	SupportedDdr      string  `json:"supported_ddr,omitempty" db:"supported_ddr"`
	MinFrequencyOfRam int64   `json:"min_frequency_of_ram,omitempty" db:"column:min_frequency_of_ram"`
	MaxFrequencyOfRam int64   `json:"max_frequency_of_ram,omitempty" db:"column:max_frequency_of_ram"`
	Price             int64   `json:"price,omitempty" db:"column:price"`
}
