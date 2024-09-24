package cosmos

import "time"

type CommercialUnit struct {
	TypePackaging     string `json:"type_packaging"`
	QuantityPackaging int    `json:"quantity_packaging"`
	Ballast           *int   `json:"ballast"`
	Layer             *int   `json:"layer"`
}

type GtinData struct {
	Gtin           int64          `json:"gtin"`
	CommercialUnit CommercialUnit `json:"commercial_unit"`
}

type NcmData struct {
	Code            string  `json:"code"`
	Description     string  `json:"description"`
	FullDescription string  `json:"full_description"`
	Ex              *string `json:"ex"`
}

type GpcData struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type BrandData struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type CestData struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ParentID    int    `json:"parent_id"`
}

type CosmosProduct struct {
	Description  string     `json:"description"`
	Gtin         int64      `json:"gtin"`
	Thumbnail    string     `json:"thumbnail"`
	Width        *float64   `json:"width"`
	Height       *float64   `json:"height"`
	Length       *float64   `json:"length"`
	NetWeight    *int       `json:"net_weight"`
	GrossWeight  *int       `json:"gross_weight"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ReleaseDate  *time.Time `json:"release_date"`
	Price        *float64   `json:"price"`
	AvgPrice     *float64   `json:"avg_price"`
	MaxPrice     float64    `json:"max_price"`
	MinPrice     float64    `json:"min_price"`
	Gtins        []GtinData `json:"gtins"`
	Origin       string     `json:"origin"`
	BarcodeImage string     `json:"barcode_image"`
	Gpc          *GpcData   `json:"gpc"`
	Ncm          *NcmData   `json:"ncm"`
	Brand        *BrandData `json:"brand"`
	Cest         *CestData  `json:"cest"`
}
