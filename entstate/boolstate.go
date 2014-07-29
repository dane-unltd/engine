package entstate

type BoolState []bool

func NewBoolState() BoolState {
	return make(BoolState, 0, 10)
}

func (dst *BoolState) Copy(src interface{}) {
	s := src.(*BoolState)
	if len(*dst) < len(*s) {
		*dst = make(BoolState, len(*s))
	}
	copy(*dst, *s)
}

func (st BoolState) Zero(id EntId) {
	st[id] = false
}

func (st *BoolState) Clone() interface{} {
	res := make(BoolState, len(*st))
	copy(res, *st)
	return &res
}

func (st BoolState) Val(id EntId) interface{} {
	return st[id]
}

func (st BoolState) Equal(v interface{}, id EntId) bool {
	b := v.(bool)
	return b == st[id]
}

func (st *BoolState) Append(n uint32) {
	for len(*st) <= int(n) {
		(*st) = append((*st), false)
	}
}
