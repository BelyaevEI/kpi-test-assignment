package saver

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/BelyaevEI/kpi-test-assignment/internal/model"
)

type Saver interface {
	SaveData()
}

type saver struct {
	sendDataToProc chan model.Data
	resDataOutProc chan model.Data
}

func New(sendDataToProc, resDataOutProc chan model.Data) Saver {
	return &saver{
		sendDataToProc: sendDataToProc,
		resDataOutProc: resDataOutProc,
	}
}

// если сохранение происходит долго можно
// добавить worker которые будут параллельно обрабатывать
// сохранение данных
func (s *saver) SaveData() {
	var wg sync.WaitGroup

	// я думаю можно вынести переменную кол-ва воркеров в конфиг,
	// чтобы можно было управлять скоростью обработки
	for i := 1; i < 5; i++ {
		wg.Add(1)
		go s.worker(&wg)
	}

	// Ожидаем завершения всех воркеров
	wg.Wait()
	close(s.resDataOutProc)
}

// сущность воркера который будет обрабатывать данные параллельно другим воркерам
func (s *saver) worker(wg *sync.WaitGroup) {
	var result int
	defer wg.Done()

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for v := range s.sendDataToProc {
		// тут нужно делать запрос на сохранение и получаеть ответ от сервера
		// но т.к. в задании не было написано какие данные мы отправляем
		// сделал имитацию этого действия
		// надеюсь на адекватность проверяющего

		// для теста эмулируем работу и отправляем ответ
		fmt.Printf("Сохраняем элемент %v\n", v)
		if rng.Intn(2) == 1 {
			result = 3
		} else {
			result = 4
		}
		fmt.Printf("Результат сохранения %v\n", result)
		v.InProgress = result
		s.resDataOutProc <- v
		time.Sleep(300 * time.Microsecond)
	}
}
