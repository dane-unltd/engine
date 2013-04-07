package helpers

type Uint32 uint32

func (v Uint32) Copy() interface{} {
	return v
}

func (v Uint32) Equals(s interface{}) bool {
	v2, ok := s.(Uint32)
	if !ok {
		return ok
	}
	return v == v2
}

type Float64 float64

func (v Float64) Copy() interface{} {
	return v
}

func (v Float64) Equals(s interface{}) bool {
	v2, ok := s.(Float64)
	if !ok {
		return ok
	}
	return v == v2
}
