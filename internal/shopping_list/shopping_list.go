package shoppinglist

type ShoppingList struct {
	CustomerId  string   `json:"customer_id"`
	Name        string   `json:"name"`
	ProductsIds []string `json:"products_ids"`
}
