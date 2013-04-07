package core

type Empty struct{}

func (e Empty) Copy() interface{} {
	return e
}

func (e Empty) Equals(v interface{}) bool {
	_, ok := v.(Empty)
	return ok
}
