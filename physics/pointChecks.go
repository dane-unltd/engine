package physics

import (
	. "github.com/dane-unltd/linalg/matrix"
)

func CheckLine(Y *DenseD, ixs []int) (s VecD, sup []int) {
	a := ixs[1]
	b := ixs[0]
	m, _ := Y.Size()

	yab := NewVecD(m).Sub(Y.ColView(b), Y.ColView(a))
	ya0 := NewVecD(m).Neg(Y.ColView(a))

	if Dot(ya0, yab) > 0 {
		sup = []int{a, b}
		s = NewVecD(m).Cross(NewVecD(m).Cross(yab, ya0), yab)
		return
	}
	sup = []int{a}
	s = ya0
	return
}

func CheckTri(Y *DenseD, ixs []int) (s VecD, sup []int) {

	a := ixs[2]
	b := ixs[1]
	c := ixs[0]

	ya0 := NewVecD(3).Neg(Y.ColView(a))
	yab := NewVecD(3).Add(Y.ColView(b), ya0)
	yac := NewVecD(3).Add(Y.ColView(c), ya0)

	normal := NewVecD(3).Cross(yab, yac)

	edge := NewVecD(3).Cross(normal, yac)
	if Dot(edge, ya0) > 0 {
		if Dot(yac, ya0) > 0 {
			sup = []int{a, c}
			s = NewVecD(3).Cross(NewVecD(3).Cross(yac, ya0), yac)
			return
		} else {
			goto abtest
		}
	} else {
		edge.Cross(yab, normal)
		if Dot(edge, ya0) > 0 {
			goto abtest
		} else {
			if Dot(normal, ya0) > 0 {
				sup = []int{a, b, c}
				s = normal
				return
			}
			sup = []int{a, c, b}
			s = normal.Neg(normal)
			return
		}
	}

abtest:
	if Dot(yab, ya0) > 0 {
		sup = []int{a, b}
		s = NewVecD(3).Cross(NewVecD(3).Cross(yab, ya0), yab)
		return
	}
	sup = []int{a}
	s = ya0

	return
}

func checkTetra(Y *DenseD, ixs []int) (s VecD, sup []int) {
	a := ixs[3]
	b := ixs[2]
	c := ixs[1]
	d := ixs[0]

	ya0 := NewVecD(3).Neg(Y.ColView(a))
	yab := NewVecD(3).Add(Y.ColView(b), ya0)
	yac := NewVecD(3).Add(Y.ColView(c), ya0)
	yad := NewVecD(3).Add(Y.ColView(d), ya0)

	face := NewVecD(3).Cross(yad, yac)

	var sup1, sup2 []int
	var s1 VecD
	inside := true

	if Dot(face, ya0) > 0 {
		inside = false
		s, sup = CheckTri(Y, []int{d, c, a})
		if len(sup) == 3 {
			return
		}
		s1 = s
		sup1 = sup
	}
	face.Cross(yab, yad)
	if Dot(face, ya0) > 0 {
		inside = false
		s, sup = CheckTri(Y, []int{b, d, a})
		if len(sup) == 3 {
			return
		}
		sup2 = sup
	}
	face.Cross(yac, yab)
	if Dot(face, ya0) > 0 {
		inside = false
		s, sup = CheckTri(Y, []int{c, b, a})
		if len(sup) == 3 {
			return
		}
		if len(sup) == 1 && len(sup1) == 1 {
			return
		}
		if sup[1] == sup2[1] || sup[1] == sup1[1] {
			return
		}
		sup = sup1
		s = s1
		return
	}

	if inside {
		sup = []int{b, c, d, a}
		s = nil
		return
	}

	return
}
