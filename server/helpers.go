package server

import "fmt"

func one[T any](items []T, err error, def *T) (*T, error) {
	if err != nil {
		return def, err
	}
	if len(items) < 1 {
		return def, fmt.Errorf("not found")
	}
	if len(items) > 1 {
		return def, fmt.Errorf("unique match not found")
	}
	return &items[0], nil
}
