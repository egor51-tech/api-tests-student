package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	charset      = lowerLetters + upperLetters + digits
)

// randomFrom возвращает строку из случайных символов указанного набора
func randomFrom(symbols string, length int) string {
	b := make([]byte, length)
	for i := range b {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		if err != nil {
			panic(fmt.Sprintf("Не удалось сгенерировать случайное число: %s", err.Error()))
		}
		b[i] = symbols[randIndex.Int64()]
	}

	return string(b)
}

// RandomString возвращает случайную строку из букв и цифр
func RandomString(length int) string {
	return randomFrom(charset, length)
}

// RandomPassword возвращает случайный пароль, в котором точно есть строчная буква, заглавная буква и цифра:
// сервер отклоняет пароли без цифры. Длина должна быть не меньше трёх.
func RandomPassword(length int) string {
	return randomFrom(lowerLetters, 1) + randomFrom(upperLetters, 1) + randomFrom(digits, 1) + randomFrom(charset, length-3)
}
