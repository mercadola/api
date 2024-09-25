package utils

type UpdateResult struct {
	ModifiedCount int64 `json:"modified_count"`
}

type DeleteResult struct {
	DeletedCount int64 `json:"deleted_count"`
}
