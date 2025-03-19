package array

import (
	"errors"
	"slices"
)

func ContainAny[T comparable](arr []T, targets []T) bool { // To check if the array contains any of the targets
	targetSet := make(map[T]struct{}, len(targets))
	for _, t := range targets {
		targetSet[t] = struct{}{}
	}

	for _, a := range arr {
		if _, found := targetSet[a]; found {
			return true
		}
	}
	return false
}

func RemoveOne[T comparable](arr []T, target T) ([]T, error) {
	for i, a := range arr {
		if a == target {
			return slices.Delete(arr, i, i+1), nil
		}
	}

	return nil, errors.New("target not found")
}
