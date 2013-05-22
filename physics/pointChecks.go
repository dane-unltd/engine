package physics

import (
	. "github.com/dane-unltd/linalg/matrix"
)

func CheckLine(Y *Dense, ixs []int) (s Vec, sup []int) {
	a := ixs[1]
	b := ixs[0]
	m, _ := Y.Size()

	yab := NewVec(m).Sub(Y.ColView(b), Y.ColView(a))
	ya0 := NewVec(m).Neg(Y.ColView(a))

	if Ddot(ya0, yab) > 0 {
		sup = []int{a, b}
		s = NewVec(m).Cross(NewVec(m).Cross(yab, ya0), yab)
		return
	}
	sup = []int{a}
	s = ya0
	return
}

func CheckTri(Y *Dense, ixs []int) (s Vec, sup []int) {

	a := ixs[2]
	b := ixs[1]
	c := ixs[0]

	ya0 := NewVec(3).Neg(Y.ColView(a))
	yab := NewVec(3).Add(Y.ColView(b), ya0)
	yac := NewVec(3).Add(Y.ColView(c), ya0)

	normal := NewVec(3).Cross(yab, yac)

	edge := NewVec(3).Cross(normal, yac)
	if Ddot(edge, ya0) > 0 {
		if Ddot(yac, ya0) > 0 {
			sup = []int{a, c}
			s = NewVec(3).Cross(NewVec(3).Cross(yac, ya0), yac)
			return
		} else {
			goto abtest
		}
	} else {
		edge.Cross(yab, normal)
		if Ddot(edge, ya0) > 0 {
			goto abtest
		} else {
			if Ddot(normal, ya0) > 0 {
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
	if Ddot(yab, ya0) > 0 {
		sup = []int{a, b}
		s = NewVec(3).Cross(NewVec(3).Cross(yab, ya0), yab)
		return
	}
	sup = []int{a}
	s = ya0

	return
}

func checkTetra(Y *Dense, ixs []int) (s Vec, sup []int) {
	a := ixs[3]
	b := ixs[2]
	c := ixs[1]
	d := ixs[0]

	ya0 := NewVec(3).Neg(Y.ColView(a))
	yab := NewVec(3).Add(Y.ColView(b), ya0)
	yac := NewVec(3).Add(Y.ColView(c), ya0)
	yad := NewVec(3).Add(Y.ColView(d), ya0)

	face := NewVec(3).Cross(yad, yac)

	var sup1, sup2 []int
	var s1 Vec
	inside := true

	if Ddot(face, ya0) > 0 {
		inside = false
		s, sup = CheckTri(Y, []int{d, c, a})
		if len(sup) == 3 {
			return
		}
		s1 = s
		sup1 = sup
	}
	face.Cross(yab, yad)
	if Ddot(face, ya0) > 0 {
		inside = false
		s, sup = CheckTri(Y, []int{b, d, a})
		if len(sup) == 3 {
			return
		}
		sup2 = sup
	}
	face.Cross(yac, yab)
	if Ddot(face, ya0) > 0 {
		s, sup = CheckTri(Y, []int{c, b, a})
		if len(sup) == 3 {
			return
		}
		if inside {
			return
		}
		if len(sup1) == 0 {
			sup1, sup2 = sup2, sup1
		}

		if len(sup) == 1 && len(sup1) == 1 {
			return
		}
		if len(sup1) != 2 || len(sup2) != 2 {
			return
		}
		if len(sup) != 2 {
			sup = sup1
			s = s1
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
