package physics

import (
	"fmt"
	. "github.com/dane-unltd/linalg/matrix"
	"sort"
)

type RigidBody struct {
	Pos, Vel VecD
	Rot      *DenseD
	MassInv  float64

	LinOpt func(c VecD) VecD
}

type Contact struct {
	Normal           VecD
	Dist             float64
	A, B             *RigidBody
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

func LinOptPoly(A *DenseD) func(VecD) VecD {
	_, n := A.Size()
	w := NewVecD(n)
	A.Tr()
	return func(c VecD) VecD {
		w.Mul(A, c)
		maxIx := w.Idmax()
		return A.ColView(maxIx).Copy().(VecD)
	}
}

func (c *Contact) Update() *Contact {
	s := NewVecD(3)
	sNeg := NewVecD(3)
	s.Sub(c.B.Pos, c.A.Pos)
	sNeg.Neg(s)

	Y := NewDenseD(3, 4)
	y := NewVecD(3)

	pA := NewVecD(3)
	pB := NewVecD(3)

	if len(c.Sup) == 4 {
		c.Sup = c.Sup[1:]
	}

	for _, sp := range c.Sup {
		pA.Mul(c.A.Rot, c.PointsA.ColView(sp))
		pB.Mul(c.B.Rot, c.PointsB.ColView(sp))
		y.Add(pA, c.A.Pos)
		y.Sub(y, pB)
		y.Sub(y, c.B.Pos)
		Y.SetCol(sp, y)
	}

	sTr, sNegTr := NewVecD(3), NewVecD(3)
	pArot := NewVecD(3)
	pBrot := NewVecD(3)
	for iter := 0; ; iter++ {
		sTr.Mul(c.A.Rot, s)
		sNegTr.Mul(c.B.Rot, sNeg)

		pA = c.A.LinOpt(sTr)
		pB = c.B.LinOpt(sNegTr)

		pArot.Mul(c.A.Rot, pA)
		pBrot.Mul(c.B.Rot, pB)

		y.Add(pArot, c.A.Pos)
		y.Sub(y, pBrot)
		y.Sub(y, c.B.Pos)

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
		fmt.Println("sup", c.Sup, i)
		c.Sup = append(c.Sup, i)

		fmt.Println("sup", c.Sup)
		fmt.Println("Y", Y)

		if iter > 0 && Dot(s, NewVecD(3).Sub(y, Y.ColView(c.Sup[0]))) < 1e-3 {
			break
		}

		s, c.Sup = MinPoly(Y, c.Sup)
		sNeg.Neg(s)
		if len(c.Sup) == 4 {
			break
		}
	}

	if s == nil {
		c.Dist = 0
		return c
	}
	c.Normal = s.Normalize(s)
	c.Dist = -Dot(s, y)
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

			if fCurr > f0-0.5*(gA.Dot(tempA)+gB.Dot(tempB)) {
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
