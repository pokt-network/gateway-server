package common

import (
	"math/rand"
)

func GetRandomElement[T any](elements []T) (T, bool) {
	if len(elements) == 0 {
		return *new(T), false
	}
	randomIndex := rand.Intn(len(elements))
	return elements[randomIndex], true
}
