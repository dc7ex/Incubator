package context

func NewContext() Context {
	return &context{
		store: make(Map),
	}
}
