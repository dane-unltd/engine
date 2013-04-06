package core

type Empty struct{}

func (e Empty) Copy() StateVar {
	return e
}

func (e Empty) Equals(v StateVar) bool {
	_, ok := v.(Empty)
	return ok
}
