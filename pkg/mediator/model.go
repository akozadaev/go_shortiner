package mediator

type Mediator interface {
	Generate(str string)
}
type ShortinerMediator struct {
	*Md5
	*Base62
}

type Md5 struct {
	Mediator Mediator
}

type Base62 struct {
	Mediator Mediator
}
