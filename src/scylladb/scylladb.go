package scylladb

import (
	"fmt"
	"github.com/gocql/gocql"
	"os"
	"sync"
)

var (
	session *gocql.Session
	once    sync.Once
)

// GetScyllaSession возвращает единственный экземпляр синглтона сессии сциллы
func GetScyllaSession() *gocql.Session {
	once.Do(func() {
		newSession, err := initScyllaSession()
		// В случае ошибки выходим с кодом 1
		if err != nil {
			fmt.Fprintf(os.Stderr, "Проблема с подключением к БД: %v\n", err)
			os.Exit(1)
		}
		// Присваиваем значение нашей переменной синглтона
		session = newSession
	})
	return session
}

// initScyllaSession Создает сессию сциллы или возвращает ошибку
func initScyllaSession() (*gocql.Session, error) {
	// Настройка кластера
	cluster := gocql.NewCluster("127.0.0.1") // Todo посадить на энвы
	cluster.Keyspace = "example"             // Todo посадить на энвы
	cluster.Consistency = gocql.Quorum

	// Создание сессии
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
