package buffer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/BelyaevEI/kpi-test-assignment/internal/model"
	"github.com/gofrs/uuid"
)

// контракт для работы с буффером
type Buffer interface {
	AddData(newData []model.Data)
	GetChannel() (chan model.Data, chan model.Data)
	ProcesingData(ctx context.Context)
}

type buffer struct {
	storage        []model.Data
	mu             sync.Mutex
	sendDataToProc chan model.Data
	resDataOutProc chan model.Data
}

func New() Buffer {
	buf := make([]model.Data, 0)
	s := make(chan model.Data, 1000)
	r := make(chan model.Data, 1000)

	return &buffer{
		storage:        buf,
		sendDataToProc: s,
		resDataOutProc: r,
	}
}

// каналы для общения с буффером
func (b *buffer) GetChannel() (chan model.Data, chan model.Data) {
	return b.sendDataToProc, b.resDataOutProc
}

// по условию раз в минуту добавляем пачку данных
// поэтому мьютексом не сильно будем тормозить воркеров
func (b *buffer) AddData(newData []model.Data) {
	b.mu.Lock()
	b.storage = append(b.storage, newData...)
	b.mu.Unlock()
}

// 0 - элемент готов к обработке
// 1 - элемент уже отправлен на обработку
// 3 - эдемент сохранен и его можно удлать из буфера
// 4 - ошибка сохранения из буфера не удаляем
// отправляем воркерам данные и получаем от них элементы на удаление
func (b *buffer) ProcesingData(ctx context.Context) {
	// запускаем три горутины на отправку данных, обработку и удаление
	// предполагается, что сохранение на стороннем ресурсе может закончится
	// провалом, поэтому необходимо сделать гарантированную доставку

	listSuccess := make(map[uuid.UUID]bool)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func(ctx context.Context) {
		defer wg.Done()
		defer close(b.sendDataToProc)
		for {
			// если контекст отменен мы продолжаем работать пока не обработаем все элементы которые запросили
			if ctx.Err() == nil || len(b.storage) != 0 {
				b.mu.Lock()
				for i := 0; i < len(b.storage); i++ {
					if b.storage[i].InProgress == 1 || b.storage[i].InProgress == 3 {
						continue
					}
					// меняем статус, потому что отправляем данный элемент
					b.storage[i].InProgress = 1
					fmt.Printf("Отправили на обработку %v\n", b.storage[i])

					// отправка данных на обработку
					b.sendDataToProc <- b.storage[i]
				}
				b.mu.Unlock()
				time.Sleep(1 * time.Second)
			} else {
				return
			}
		}

	}(ctx)

	// обработка сохраненных элементов
	go func() {
		defer wg.Done()
		for v := range b.resDataOutProc {
			if v.InProgress == 3 {
				listSuccess[uuid.UUID(v.ID)] = true
			}
			// меняем статус после обработки
			for i := 0; i < len(b.storage); i++ {
				//находим нужный элемент по ID
				if b.storage[i].ID == v.ID {
					b.storage[i].InProgress = v.InProgress
				}
			}
		}
	}()

	// удаление из буффера успешно сохраненных записей
	go func(ctx context.Context) {
		var newStorage []model.Data
		defer wg.Done()
		for {
			// если контекст отменен, то мы продолжаем работать пока не сохраним все данные
			if ctx.Err() == nil || len(b.storage) != 0 || len(listSuccess) != 0 {
				if len(b.storage) >= len(listSuccess) {
					newStorage = make([]model.Data, 0, len(b.storage)-len(listSuccess))
				} else {
					newStorage = make([]model.Data, 0)
				}
				b.mu.Lock()
				for _, v := range b.storage {
					_, ok := listSuccess[uuid.UUID(v.ID)]
					if !ok {
						newStorage = append(newStorage, v)
					}
				}
				//чистим словарь
				listSuccess = make(map[uuid.UUID]bool)

				// заменяем хранилище на новое без сохранненых элементов
				fmt.Printf("Сохранили %v элементов\n", len(b.storage)-len(newStorage))
				b.storage = newStorage
				fmt.Printf("После удаления элементов в буффере %v\n", len(b.storage))
				b.mu.Unlock()

				// удаляем каждые 15 секунд
				// я думаю можно вынести переменную в конфиг,
				// чтобы можно было управлять скоростью обработки
				time.Sleep(15 * time.Second)
			} else {
				return
			}
		}

	}(ctx)
	wg.Wait()
}
