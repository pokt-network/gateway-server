package common

import (
	"math/rand"
)

func GetRandomElement[T any](elements []T) T {
	if len(elements) == 0 {
		return *new(T)
	}
	randomIndex := rand.Intn(len(elements))
	return elements[randomIndex]
}
