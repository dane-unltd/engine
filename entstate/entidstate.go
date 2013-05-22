package entstate

type EntIdState []EntId

func NewEntIdState() EntIdState {
	return make(EntIdState, 0, 10)
}

func (st *EntIdState) Clone() interface{} {
	rst := make(EntIdState, len(*st))
	copy(rst, *st)
	return &rst
}

func (st EntIdState) Zero(id EntId) {
	st[id] = 0
}

func (dst *EntIdState) Copy(src interface{}) {
	s := src.(*EntIdState)
	if len(*dst) < len(*s) {
		*dst = make(EntIdState, len(*s))
	}
	copy(*dst, *s)
}

func (st EntIdState) Val(id EntId) interface{} {
	return st[id]
}

func (st EntIdState) Equal(v interface{}, id EntId) bool {
	entId := v.(EntId)
	return entId == st[id]
}

func (st *EntIdState) Append() {
	(*st) = append((*st), 0)
}
