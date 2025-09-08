package common

import (
	"math/rand"
	"slices"
)

type UniqueRand struct {
	size int
	used []int
}

func NewUniqueRand(size int) *UniqueRand {
	return &UniqueRand{size: size, used: []int{}}
}

func (u *UniqueRand) Int() int {
	for {
		i := rand.Intn(u.size)
		if !slices.Contains(u.used, i) {
			u.used = append(u.used, i)
			return i
		}
	}
}
