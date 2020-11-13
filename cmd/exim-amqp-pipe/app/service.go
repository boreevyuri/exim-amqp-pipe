package app

type Service interface {
	OnInit(*Event)
}

type PublishService interface {
	Service
	OnPublish(*Event)
}
