package core

type Uint32 uint32

func (v Uint32) Copy() StateVar {
	return v
}

func (v Uint32) Equals(s StateVar) bool {
	v2, ok := s.(Uint32)
	if !ok {
		return ok
	}
	return v == v2
}
