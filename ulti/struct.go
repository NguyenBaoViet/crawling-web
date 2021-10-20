package ulti

import (
	"github.com/google/uuid"
)

type Category struct {
	ID       uuid.UUID `json:"ID"`
	Name     string    `json:"Name"`
	URL      string    `json:"URL"`
	SubLevel int       `json:"SubLevel"`
	ParentID uuid.UUID `json:"ParentID"`
}

type ProducInfo struct {
	Name       string    `json:"Name"`
	URL        string    `json:"URL"`
	CategoryID uuid.UUID `json:"CategoryID"`
}

type ProductDetail struct {
	ID       uuid.UUID `json:"ID"`
	SKU      string    `json:"SKU"`
	Name     string    `json:"Name"`
	Price    int64     `json:"Price"`
	OldPrice int64     `json:"OldPrice"`
	Color    string    `json:"Color"`
	Img      []string  `json:"Img"`
	Size     string    `json:"Size"`
}
