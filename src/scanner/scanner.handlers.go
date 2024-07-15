package router_scanner

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ScanResult Структура для ответа
type ScanResult struct {
	Website    string `json:"website"`
	MaxWorkers int    `json:"max_workers"`
	Ports      []int  `json:"ports"`
}

func RegisterScannerHandlers(mux *http.ServeMux) {
	// Регистрация роутов и их обработчиков
	mux.HandleFunc("/scan-tcp/{website}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод недопустим!", http.StatusMethodNotAllowed)
			return
		}
		// Извлекаем имя вебсайта из URL
		website := r.PathValue("website")

		// Проверяем, что имя вебсайта не пустое
		if website == "" {
			http.Error(w, "Не указано имя вебсайта", http.StatusBadRequest)
			return
		}
		// Извлекаем максимальное количество рабочих из параметра запроса
		maxWorkersStr := r.URL.Query().Get("max-workers")
		maxWorkers, err := strconv.Atoi(maxWorkersStr)
		if err != nil {
			// Если параметр не задан или не может быть преобразован в число, используем значение по умолчанию
			maxWorkers = 10 // Значение по умолчанию
		}

		// Вызываем наш скан
		openedPorts := scan("tcp", website, maxWorkers)
		// Создаем объект ScanResult
		result := ScanResult{
			Website:    website,
			MaxWorkers: maxWorkers,
			Ports:      openedPorts,
		}
		// Сериализуем объект в JSON
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Не удалось сериализовать в JSON", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовок Content-Type для ответа
		w.Header().Set("Content-Type", "application/json")

		// Отправляем ответ с данными в формате JSON
		w.Write(jsonBytes)
	})

	mux.HandleFunc("/scan-http/{website}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод недопустим!", http.StatusMethodNotAllowed)
			return
		}
		// Извлекаем имя вебсайта из URL
		website := r.PathValue("website")

		// Проверяем, что имя вебсайта не пустое
		if website == "" {
			http.Error(w, "Не указано имя вебсайта", http.StatusBadRequest)
			return
		}
		// Извлекаем максимальное количество рабочих из параметра запроса
		maxWorkersStr := r.URL.Query().Get("max-workers")
		maxWorkers, err := strconv.Atoi(maxWorkersStr)
		if err != nil {
			// Если параметр не задан или не может быть преобразован в число, используем значение по умолчанию
			maxWorkers = 10 // Значение по умолчанию
		}

		// Вызываем наш скан
		openedPorts := scan("http", website, maxWorkers)
		// Создаем объект ScanResult
		result := ScanResult{
			Website:    website,
			MaxWorkers: maxWorkers,
			Ports:      openedPorts,
		}
		// Сериализуем объект в JSON
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Не удалось сериализовать в JSON", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовок Content-Type для ответа
		w.Header().Set("Content-Type", "application/json")

		// Отправляем ответ с данными в формате JSON
		w.Write(jsonBytes)
	})
}
