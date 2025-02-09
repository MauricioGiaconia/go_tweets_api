package models

type PaginationWithID struct {
	ID     int64
	Limit  int64
	Offset int64
}
