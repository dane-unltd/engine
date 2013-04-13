package physics

import (
	"github.com/dane-unltd/engine/core"
	. "github.com/dane-unltd/linalg/matrix"
	"sort"
)

type SupportFunc func(VecD) VecD

func (s SupportFunc) Copy() interface{} {
	return s
}

func (s SupportFunc) Equals(interface{}) bool {
	return true
}

type Contact struct {
	Normal           VecD
	Dist             float64
	A, B             core.EntId
	PointsA, PointsB *DenseD
	Sup              []int
}

func NewContact() *Contact {
	c := Contact{}
	c.Sup = make([]int, 0)
	c.PointsA = NewDenseD(3, 4)
	c.PointsB = NewDenseD(3, 4)
	return &c
}

func (c *Contact) Copy() interface{} {
	cn := Contact{}
	cn.Normal = c.Normal.Copy().(VecD)
	cn.Dist = c.Dist
	cn.A, cn.B = c.A, c.B
	cn.PointsA = c.PointsA.Copy().(*DenseD)
	cn.PointsB = c.PointsB.Copy().(*DenseD)
	cn.Sup = make([]int, len(c.Sup))
	copy(cn.Sup, c.Sup)
	return &cn
}

func (c *Contact) Update(posA, posB VecD, rotA, rotB *DenseD,
	supFunA, supFunB SupportFunc) *Contact {
	s := NewVecD(3)
	sNeg := NewVecD(3)
	if c.Normal == nil {
		s.Sub(posB, posA)
	} else {
		s = c.Normal
	}
	sNeg.Neg(s)

	Y := NewDenseD(3, 4)
	y := NewVecD(3)

	pA := NewVecD(3)
	pB := NewVecD(3)

	sTr, sNegTr := NewVecD(3), NewVecD(3)
	pArot := NewVecD(3)
	pBrot := NewVecD(3)

	if len(c.Sup) == 4 {
		c.Sup = c.Sup[1:]
	}

	for _, sp := range c.Sup {
		pA.Mul(rotA, c.PointsA.ColView(sp))
		pB.Mul(rotB, c.PointsB.ColView(sp))
		y.Add(pA, posA)
		y.Sub(y, pB)
		y.Sub(y, posB)
		Y.SetCol(sp, y)
	}
	if len(c.Sup) > 0 {
		s, c.Sup = MinPoly(Y, c.Sup)
	} else {
		sTr.Mul(rotA, s)
		sNegTr.Mul(rotB, sNeg)

		pA = supFunA(sTr)
		pB = supFunB(sNegTr)

		pArot.Mul(rotA, pA)
		pBrot.Mul(rotB, pB)

		y.Add(pArot, posA)
		y.Sub(y, pBrot)
		y.Sub(y, posB)

		c.PointsA.SetCol(0, pA)
		c.PointsB.SetCol(0, pB)
		Y.SetCol(0, y)
		c.Sup = append(c.Sup, 0)
		s.Neg(y)
		sNeg.Neg(s)
	}

	for iter := 0; ; iter++ {
		sTr.Mul(rotA, s)
		sNegTr.Mul(rotB, sNeg)

		pA = supFunA(sTr)
		pB = supFunB(sNegTr)

		pArot.Mul(rotA, pA)
		pBrot.Mul(rotB, pB)

		y.Add(pArot, posA)
		y.Sub(y, pBrot)
		y.Sub(y, posB)

		if Ddot(s, NewVecD(3).Sub(y, Y.ColView(c.Sup[0]))) < 1e-3 {
			break
		}

		supCheck := []bool{true, true, true, true}
		for _, sv := range c.Sup {
			supCheck[sv] = false
		}
		i := 0
		for ; i < 4; i++ {
			if supCheck[i] {
				break
			}
		}

		c.PointsA.SetCol(i, pA)
		c.PointsB.SetCol(i, pB)
		Y.SetCol(i, y)
		c.Sup = append(c.Sup, i)

		s, c.Sup = MinPoly(Y, c.Sup)
		if len(c.Sup) == 4 {
			break
		}
		sNeg.Neg(s)
	}

	if s == nil {
		c.Dist = 0
		s = NewVecD(3).Sub(posB, posA)
		c.Normal = s.Normalize(s)
		return c
	}
	c.Normal = s.Normalize(s)
	c.Dist = -Ddot(s, y)
	return c
}

func MinPoly(Y *DenseD, sup []int) (s VecD, supRes []int) {
	n := len(sup)
	if n == 1 {
		s = NewVecD(3).Neg(Y.ColView(sup[0]))
		supRes = sup
		return
	}
	if n == 2 {
		s, supRes = CheckLine(Y, sup)
		return
	}
	if n == 3 {
		s, supRes = CheckTri(Y, sup)
		return
	}
	if n == 4 {
		s, supRes = checkTetra(Y, sup)
		return
	}
	panic("too many points")
}

func ProjectOnSimplex(v VecD) VecD {
	n := len(v)
	vSort := v.Copy().(VecD)
	sort.Float64s([]float64(vSort))

	var wl float64
	for i := range vSort {
		wl = vSort[i]
		S := vSort[i:].Asum()
		P := S - float64(n-i)*wl
		if P < 1 {
			wl = (S - 1) / float64(n-i)
			break
		}
	}

	for i := range v {
		v[i] = v[i] - wl
	}
	return v.Max(v, 0)
}

/*
func closestPoint(A, B *DenseD, sA, sB VecD) (a, b VecD) {
	s := 1.0
	n, _ := A.Size()

	d := NewVecD(n)
	tempA := NewVecD(len(sA))
	tempB := NewVecD(len(sB))

	sAhat := NewVecD(len(sA))
	sBhat := NewVecD(len(sB))

	ProjectOnSimplex(sA)
	ProjectOnSimplex(sB)

	a = A.ApplyTo(sA)
	b = B.ApplyTo(sB)

	nIter := 0
	for {
		d.Sub(b, a)

		gA := A.ApplyToTr(d)
		gB := B.ApplyToTr(d.Neg(d))

		if terminationCheck(sA, gA) && terminationCheck(sB, gB) {
			break
		}
		nIter++

		f0 := d.Nrm2Sq()
		fPrev := f0

		s = s * 2
		for {
			ProjectOnSimplex(sAhat.Add(sA, tempA.Mul(s, gA)))
			ProjectOnSimplex(sBhat.Add(sB, tempB.Mul(s, gB)))

			fCurr := d.Sub(B.ApplyTo(sBhat),
				A.ApplyTo(sAhat)).Nrm2Sq()
			tempA.Sub(sAhat, sA)
			tempB.Sub(sBhat, sB)

			if fCurr > f0-0.5*(gA.Ddot(tempA)+gB.Ddot(tempB)) {
				s = s * 0.5
				fPrev = fCurr
			} else {
				if fPrev < fCurr {
					s = 2 * s
				}
				break
			}
		}

		ProjectOnSimplex(sA.Add(sA, tempA.Mul(s, gA)))
		ProjectOnSimplex(sB.Add(sB, tempB.Mul(s, gB)))

		a = A.ApplyTo(sA)
		b = B.ApplyTo(sB)
	}
	return
}



func terminationCheck(s, g VecD) bool {
	if g.Nrm2Sq() < 1e-6 {
		return true
	}
	q0 := NewVecD(len(s))
	q := NewVecD(len(s))
	for i := range q0 {
		q0[i] = 1
	}

	var Q *DenseD

	first := true
	for i, v := range s {
		if v < 0.001 {
			q[i] = -1
			if first {
				Q = FromArrayD([]float64(q), false, len(q), 1)
				first = false
			} else {
				Q.AddCol(q)
			}
			q[i] = 0
			q0[i] = 0
		}
	}
	if Q == nil {
		Q = FromArrayD([]float64(q0.Normalize(q0)),
			false, len(q), 1)
	} else {
		Q.AddCol(q0.Normalize(q0))
	}
	m, n := Q.Size()

	rn := OnesVec(n)
	rn.Mul(-1.0/math.Sqrt(float64(m-n+1)), rn)
	rn[n-1] = 1
	R := Eye(n)
	R.SetCol(n-1, rn)

	gtrans := Q.ApplyToTr(g)
	coeffs := R.SolveTriU(gtrans)

	terminate := true
	for i := 0; i < len(coeffs)-1; i++ {
		if coeffs[i] < -1e-3 {
			terminate = false
		}
	}
	if terminate {
		gproj := Q.ApplyTo(R.ApplyTo(coeffs))
		diff := NewVecD(m).Sub(gproj, g).Nrm2Sq()
		if diff > 1e-3 {
			terminate = false
		}
	}
	return terminate
}
*/
