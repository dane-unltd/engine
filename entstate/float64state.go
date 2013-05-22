package entstate

type Float64State []float64

func NewFloat64State() Float64State {
	return make(Float64State, 0, 10)
}

func (st *Float64State) Clone() interface{} {
	res := make(Float64State, len(*st))
	copy(res, *st)
	return &res
}

func (st Float64State) Zero(id EntId) {
	st[id] = 0
}

func (dst *Float64State) Copy(src interface{}) {
	s := src.(*Float64State)
	if len(*dst) < len(*s) {
		*dst = make(Float64State, len(*s))
	}
	copy(*dst, *s)
}

func (st Float64State) Val(id EntId) interface{} {
	return st[id]
}

func (st Float64State) Equal(v interface{}, id EntId) bool {
	f64 := v.(float64)
	return f64 == st[id]
}

func (st *Float64State) Append() {
	(*st) = append((*st), 0)
}
