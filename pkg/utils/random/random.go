package random

import (
	"crypto/rand"
	"math/big"
)

type Generator interface {
	GenerateRandomIntString(digit int) (string, error)
}

type GeneratorImpl struct{}

func NewGenerator() Generator {
	return &GeneratorImpl{}
}

func (GeneratorImpl) GenerateRandomIntString(digit int) (string, error) {
	//	Create random number with specified digit
	b := make([]rune, digit)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return "", err
		}
		b[i] = rune(num.Int64() + 48)
	}
	return string(b), nil
}
