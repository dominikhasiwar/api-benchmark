package models

type ListResult[T any] struct {
	LastEvaluatedKey string `json:"lastEvaluatedKey"`
	Items            []T    `json:"items"`
}
