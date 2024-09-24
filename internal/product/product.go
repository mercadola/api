package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                string    `json:"id" bson:"id"`
	Description       string    `json:"description" bson:"description"`
	Ean               int64     `json:"ean" bson:"ean"`
	Thumbnail         string    `json:"thumbnail" bson:"thumbnail"`
	Width             float64   `json:"width" bson:"width"`
	Height            float64   `json:"height" bson:"height"`
	Length            float64   `json:"length" bson:"length"`
	Ncm               string    `json:"ncm" bson:"ncm"`
	TypePackaging     string    `json:"type_packaging" bson:"type_packaging"`
	QuantityPackaging int       `json:"quantity_packaging" bson:"quantity_packaging"`
	BarcodeImage      string    `json:"barcode_image" bson:"barcode_image"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
}

func (p *Product) New(description, thumbnail, ncm, typePackaging, barcodeImage string, ean int64, width, height, length float64, quantityPackaging int) *Product {
	now := time.Now()
	p.ID = uuid.New().String()
	p.Description = description
	p.Ean = ean
	p.Thumbnail = thumbnail
	p.Width = width
	p.Height = height
	p.Length = length
	p.Ncm = ncm
	p.TypePackaging = typePackaging
	p.QuantityPackaging = quantityPackaging
	p.BarcodeImage = barcodeImage
	p.CreatedAt = now
	p.UpdatedAt = now
	return p
}

type CreateByEanDto struct {
	Ean string `json:"ean" validate:"required"`
}
