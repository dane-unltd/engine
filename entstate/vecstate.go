package entstate

import "github.com/dane-unltd/engine/math3"

type VecState []math3.Vec

func NewVecState() VecState {
	return make(VecState, 0, 10)
}

func (dst *VecState) Copy(src interface{}) {
	s := src.(*VecState)
	if len(*dst) < len(*s) {
		*dst = make(VecState, len(*s))
	}
	copy(*dst, *s)
}

func (st VecState) Zero(id EntId) {
	st[id].Zero()
}

func (st *VecState) Clone() interface{} {
	res := make(VecState, len(*st))
	copy(res, *st)
	return &res
}

func (st VecState) Val(id EntId) interface{} {
	return &st[id]
}

func (st VecState) Equal(v interface{}, id EntId) bool {
	vec := v.(*math3.Vec)
	return vec.Equals(&st[id])
}

func (st *VecState) Append(n uint32) {
	for len(*st) <= int(n) {
		(*st) = append((*st), math3.Vec{})
	}
}
