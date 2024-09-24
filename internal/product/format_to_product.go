package product

import "github.com/mercadola/api/internal/cosmos"

func FormatToProduct(data cosmos.CosmosProduct) *Product {
	var typePackaging, ncm string
	var quantityPackaging int

	if len(data.Gtins) > 0 && data.Gtins[0].CommercialUnit.TypePackaging != "" {
		typePackaging = data.Gtins[0].CommercialUnit.TypePackaging
		quantityPackaging = data.Gtins[0].CommercialUnit.QuantityPackaging
	}

	if data.Ncm != nil {
		ncm = data.Ncm.Code
	}
	product := Product{}
	return product.New(data.Description, data.Thumbnail, ncm, typePackaging, data.BarcodeImage, data.Gtin,
		getFloat64Value(data.Width, 0.0), getFloat64Value(data.Height, 0.0),
		getFloat64Value(data.Length, 0.0), quantityPackaging)
}

func getFloat64Value(value *float64, defaultValue float64) float64 {
	if value != nil {
		return *value
	}
	return defaultValue
}
