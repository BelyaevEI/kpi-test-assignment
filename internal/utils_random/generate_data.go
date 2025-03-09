package utilsrandom

import (
	"math/rand"
	"time"
)

func GenerateRandomSlice(size int) []int {
	// Создаем новый источник случайных чисел с использованием текущего времени.
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source) // Создаем новый генератор случайных чисел.

	// Создаем слайс для хранения случайных чисел.
	slice := make([]int, size)

	// Заполняем слайс случайными числами от 1 до 10.
	for i := 0; i < size; i++ {
		slice[i] = rng.Intn(10) + 1 // rng.Intn(10) дает число от 0 до 9, поэтому добавляем 1.
	}

	return slice
}
