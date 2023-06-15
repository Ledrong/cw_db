package models

type Videocard struct {
	Videocardid int64  `json:"videocardid,omitempty" db:"column:videocardid; primaryKey"`
	Name        string `json:"name,omitempty" db:"column:name"`
	Brand       string `json:"brand,omitempty" db:"column:brand"`
	Series      int64  `json:"series,omitempty" db:"column:series"`
	Vmemory     int64  `json:"vmemory,omitempty" db:"column:vmemory"`
	Price       int64  `json:"price,omitempty" db:"column:price"`
}
