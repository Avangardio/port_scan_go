package router_scanner

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

func scan(protocol string, hostname string, maxWorkers int) []int {
	var wg sync.WaitGroup
	ports := make(chan int, 100)
	openedPorts := make(chan int, 65535)

	// Создаем N воркеров, го рантайм сам им распределит значения
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go createScannerWorker(protocol, hostname, ports, openedPorts, &wg)
	}

	// Записываем сколько возможно значений в канал портов
	for port := 1; port <= 65535; port++ {
		ports <- port
	}
	// Закрываем канал, записи значений будут недоступны, чтения - доступны
	close(ports)
	wg.Wait()
	close(openedPorts)

	var result []int
	for openedPort := range openedPorts {
		fmt.Printf("%d", openedPort)
		result = append(result, openedPort)
	}
	fmt.Print(result)
	return result
}

func createScannerWorker(protocol string, hostname string, ports <-chan int, openedPorts chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Тут на самом деле хорошо будет видно, что го рантайм хорошо распределяет значения между подписчиками, примерно равномерно порты будут распределяться
	// Горутины воркеров будут читать канал пока он не опустеет и закроется
	for port := range ports {
		if protocol == "tcp" {
			address := fmt.Sprintf("%s:%d", hostname, port)
			conn, err := net.DialTimeout(protocol, address, 1*time.Second)
			if err == nil {
				_ = conn.Close()
				openedPorts <- port
			}
		} else {
			client := http.Client{
				Timeout: time.Second * 1, // Таймаут запроса 1 секунда
			}

			// Отправляем GET запрос на указанный порт (http)
			_, err := client.Get(fmt.Sprintf("http://%s:%d", hostname, port))
			if err == nil {
				openedPorts <- port
			}
			// Отправляем GET запрос на указанный порт (https)
			_, err = client.Get(fmt.Sprintf("https://%s:%d", hostname, port))
			if err == nil {
				openedPorts <- port
			}
		}
	}
}
