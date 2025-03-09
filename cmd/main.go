package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/BelyaevEI/kpi-test-assignment/internal/buffer"
	"github.com/BelyaevEI/kpi-test-assignment/internal/config"
	"github.com/BelyaevEI/kpi-test-assignment/internal/gracefullshutdown"
	"github.com/BelyaevEI/kpi-test-assignment/internal/receiver"
	"github.com/BelyaevEI/kpi-test-assignment/internal/saver"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(3)

	//считываем данные из конфига
	cfg, err := config.Load("./../config.env")
	if err != nil {
		log.Fatalf("failed read config: %s", err.Error())
	}

	// создаем буффер
	buffer := buffer.New()

	//создаем "получателя данных"
	receiver := receiver.New(cfg.GetUrlGetData(), cfg.GetBearerToken(), buffer)

	// создаем "сохранителя данных"
	saver := saver.New(buffer.GetChannel())

	// опрашиваем источник данных
	go func() {
		defer wg.Done()
		receiver.GetData(ctx)
	}()

	// обрабатываем входящие данные
	go func() {
		time.Sleep(1 * time.Second)
		defer wg.Done()
		buffer.ProcesingData(ctx)
	}()

	// сохраняем входящие элементы
	go func() {
		time.Sleep(3 * time.Second)
		defer wg.Done()
		saver.SaveData()
	}()

	// если нужно выключиь программу
	gracefullshutdown.GracefulShutdown(ctx, cancel, wg)
	fmt.Println("Программа завершена, все данные сохранены! До свидания.")
}
