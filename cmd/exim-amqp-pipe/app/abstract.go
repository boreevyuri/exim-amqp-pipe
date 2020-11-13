package app

import (
	"io/ioutil"
	"log"
)

const (
	configFail = "unable to read config"
)

type Abstract struct {
	configFilename string
	services       []interface{}
	events         chan *Event
	done           chan bool
}

func (a *Abstract) SetConfigFilename(configFilename string) {
	a.configFilename = configFilename
}

func (a *Abstract) IsValidConfigFilename(filename string) bool {
	return len(filename) > 0
}

func (a *Abstract) SetEvents(events chan *Event) {
	a.events = events
}

func (a *Abstract) Events() chan *Event {
	return a.events
}

func (a *Abstract) SetDone(done chan bool) {
	a.done = done
}

func (a *Abstract) Done() chan bool {
	return a.done
}

func (a *Abstract) Services() []interface{} {
	return a.services
}

func (a *Abstract) FireInit(event *Event, abstractService interface{}) {
	//service := abstractService.(Service)
	//service.OnInit(event)
}

func (a *Abstract) FireRun(event *Event, abstractService interface{})    {}
func (a *Abstract) FireFinish(event *Event, abstractService interface{}) {}
func (a *Abstract) Init(event *Event)                                    {}
func (a *Abstract) Run()                                                 {}

func (a *Abstract) run(app Application, event *Event) {
	app.SetDone(make(chan bool))
	app.SetEvents(make(chan *Event, 1))
	go func() {
		for {
			select {
			case event := <-app.Events():
				// Если получаем в канал эвент инициализации - читаем конфиг
				if event.Type == InitAppEvent {
					bytes, err := ioutil.ReadFile(a.configFilename)
					if err != nil {
						log.Fatal(configFail, err)
					}
					event.Data = bytes
					app.Init(event)
				}

				for _, service := range app.Services() {
					switch event.Type {
					// Пробегаемся тем же сигналом по сервисам приложения
					case InitAppEvent:
						app.FireInit(event, service)
					case RunAppEvent:
						app.FireRun(event, service)
					case FinishAppEvent:
						app.FireFinish(event, service)
					}
				}

				switch event.Type {
				case InitAppEvent:
					//Все уже инициализировано. Значит меняем сигнал на запуск и отправляем в канал
					event.Type = RunAppEvent
					app.Events() <- event
				case FinishAppEvent:
					//Получили эвент финиша. Закрываем канал эвентов, отсылаем Done
					close(app.Events())
					app.Done() <- true
				}
			}
		}
	}()

	app.Events() <- event
	<-app.Done()
}
