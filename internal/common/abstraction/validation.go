package abstractions

type IValidator interface {
	Validate(data interface{}) (err error)
}
