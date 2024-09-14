package product

import "time"

type CommercialUnit struct {
	TypePackaging     string      `json:"type_packaging"`
	QuantityPackaging int         `json:"quantity_packaging"`
	Ballast           interface{} `json:"ballast"`
	Layer             interface{} `json:"layer"`
}

type Gtin struct {
	Gtin           int64          `json:"gtin"`
	CommercialUnit CommercialUnit `json:"commercial_unit"`
}
type Gpc struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
type Ncm struct {
	Code            string      `json:"code"`
	Description     string      `json:"description"`
	FullDescription string      `json:"full_description"`
	Ex              interface{} `json:"ex"`
}

type Product struct {
	Description  string      `json:"description"`
	Gtin         int64       `json:"gtin"`
	Thumbnail    string      `json:"thumbnail"`
	Width        float64     `json:"width"`
	Height       float64     `json:"height"`
	Length       float64     `json:"length"`
	NetWeight    interface{} `json:"net_weight"`
	GrossWeight  interface{} `json:"gross_weight"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	ReleaseDate  interface{} `json:"release_date"`
	Price        interface{} `json:"price"`
	AvgPrice     interface{} `json:"avg_price"`
	MaxPrice     float64     `json:"max_price"`
	MinPrice     float64     `json:"min_price"`
	Gtins        []Gtin      `json:"gtins"`
	Origin       string      `json:"origin"`
	BarcodeImage string      `json:"barcode_image"`
	Gpc          Gpc         `json:"gpc"`
	Ncm          Ncm         `json:"ncm"`
}
