package main

import (
	"dev11/app"
	"dev11/config"
	"dev11/storage"

	"github.com/sirupsen/logrus"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)
/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func main() {
	MyLogger := logrus.New()

	MyLogger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
		PrettyPrint:      true,
	}

	config, err := config.Get()
	if err != nil {
		MyLogger.WithFields(logrus.Fields{
			"func":    "config.Get",
			"package": "main",
		}).Fatal(err)
	}

	EventStorage := storage.New()

	app := app.NewApp(MyLogger, config, EventStorage)

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", app.CreateEvent)
	mux.HandleFunc("/update_event", app.UpdateEvent)
	mux.HandleFunc("/delete_event", app.DeleteEvent)
	mux.HandleFunc("/events_for_day", app.GetEventsForDay)
	mux.HandleFunc("/events_for_week", app.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", app.GetEventsForMonth)

	MWmux := app.LogMiddleware(mux)

	srv := &http.Server{
		Addr:         config.Host + ":" + config.Port,
		Handler:      MWmux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err = srv.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
			MyLogger.Info("Server at :8080 port Stopped")
		default:
			MyLogger.WithFields(logrus.Fields{
				"func":    "srv.ListenAndServe",
				"package": "main",
			}).Fatal(err)
		}
	}()

	MyLogger.Info("Server at :8080 port Start")

	signalChanel := make(chan os.Signal, 1)
	defer close(signalChanel)

	signal.Notify(signalChanel, syscall.SIGTERM, syscall.SIGINT)

	<-signalChanel

	MyLogger.Info("server at :8080 port Shutting Down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	err = srv.Shutdown(ctx)
	if err != nil {
		MyLogger.WithFields(logrus.Fields{
			"func":    "srv.Shutdown",
			"package": "main",
		}).Fatal(err)
	}

	cancel()
}
