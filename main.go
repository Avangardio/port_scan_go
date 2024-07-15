package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"portScanGo/src/scanner"
	"syscall"
	"time"
)

func main() {
	// Грузим возможный env
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Ошибка при загрузке энвайромента")
	}
	// Ищем порт
	var port string
	if os.Getenv("PORT") == "" {
		port = ":8080"
	} else {
		port = os.Getenv("PORT")
	}

	// Создание нового ServeMux для маршрутов
	mux := http.NewServeMux()

	// Вынес логику хендлеры в другой файл
	router_scanner.RegisterScannerHandlers(mux)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Запуск сервера в горутине чтобы ожидать сиги
	go func() {
		fmt.Printf("Сервер запускается на http://localhost%s\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Ожидаем сигналы остановки для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Блокируем выполнение до получения сигналов
	<-quit
	fmt.Printf("Сервер выключается...\n")

	// Настройка контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Вызов Shutdown для graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при выключении сервера: %v\n", err)
	}

	fmt.Println("Сервер остановлен!")
}
