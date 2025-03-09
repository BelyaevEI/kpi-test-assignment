package receiver

import (
	"context"
	"fmt"
	"time"

	"github.com/BelyaevEI/kpi-test-assignment/internal/buffer"
	"github.com/BelyaevEI/kpi-test-assignment/internal/model"
	utilsrandom "github.com/BelyaevEI/kpi-test-assignment/internal/utils_random"

	"github.com/google/uuid"
)

// Интерфейс для работы с сущностью получателя данных
type Receiver interface {
	GetData(ctx context.Context)
}

type receiver struct {
	urlGetData, bearerToken string
	buffer                  buffer.Buffer
}

func New(urlGetData, bearerToken string, buffer buffer.Buffer) Receiver {
	return &receiver{
		urlGetData:  urlGetData,
		bearerToken: bearerToken,
		buffer:      buffer,
	}
}

func (r *receiver) GetData(ctx context.Context) {

	// тут мы должны делать запрос для получения данных
	// с последующим разложением json и упаковкой в удобный вид
	// схемы json в задании не было, через postman раскрывать схему не очень
	// поэтому я принял решение для источника данных генерировать числа
	// ниже пример, как можно было бы сделать, будь схема в задании :)

	// client := &http.Client{}

	// jsonData, err := json.Marshal(formData)
	// if err != nil {
	// 	return err
	// }

	// req, err := http.NewRequest("POST", r.urlGetData, bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	return err
	// }

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+r.bearerToken)

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("failed to save fact, status code: %d", resp.StatusCode)
	// }

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// получаем и преобразуем данные
			readyData := r.convertToOutput(utilsrandom.GenerateRandomSlice(1000))

			// добавим данные на сохранение
			r.buffer.AddData(readyData)
			fmt.Printf("получили на обработку %v элементов\n", len(readyData))
			time.Sleep(55 * time.Second)
		}
	}
}

func (r *receiver) convertToOutput(inputData []int) []model.Data {
	outputData := make([]model.Data, len(inputData))

	for i := 0; i < len(inputData); i++ {
		outputData[i].Data = inputData[i]
		outputData[i].ID = uuid.New()
		outputData[i].InProgress = 0
	}

	return outputData
}
