package htick

type Expression string

func (e Expression) ToConfigure() string {
	return string(e)
}

func (e Expression) Validate() error {
	return nil
}
