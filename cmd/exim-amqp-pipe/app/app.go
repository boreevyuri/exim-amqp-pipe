package app

type EventType int

const (
	InitAppEvent EventType = iota
	RunAppEvent
	FinishAppEvent
)

var (
	//App текущего приложения
	App Application

	//Services для создания итератора по сервисам
	Services []interface{}
)

type Event struct {
	Type EventType
	Data []byte
	//Args map[string]interface{}
}

func NewEvent(kind EventType) *Event {
	return &Event{Type: kind}
}

type Application interface {
	SetConfigFilename(string)
	IsValidConfigFilename(string) bool

	SetEvents(chan *Event)
	Events() chan *Event

	SetDone(chan bool)
	Done() chan bool

	Services() []interface{}

	FireInit(*Event, interface{})
	FireRun(*Event, interface{})

	FireFinish(*Event, interface{})
	Init(*Event)
	Run()
	//Timeout() Timeout
}
