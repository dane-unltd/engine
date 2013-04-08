package physics

import (
	"fmt"
	. "github.com/dane-unltd/linalg/matrix"
	"math"
	"sort"
)

type RigidBody struct {
	Pos, Vel VecD
	Rot      *MatD
	MassInv  float64

	LinOpt func(c VecD) VecD
}

type Contact struct {
	Normal           VecD
	Dist             float64
	A, B             *RigidBody
	PointsA, PointsB *MatD
}

func (c *Contact) Update() *Contact {
	fmt.Println("creating contact")
	v := ZeroVec(len(c.A.Pos))
	vNeg := ZeroVec(len(c.A.Pos))
	v.Sub(c.A.Pos, c.B.Pos)
	vNeg.Neg(v)

	if c.PointsA == nil {
		pA := c.A.LinOpt(c.A.Rot.ApplyTo(vNeg))
		pB := c.B.LinOpt(c.B.Rot.ApplyTo(v))
		c.PointsA = FromArrayD([]float64(pA), false, len(pA), 1)
		c.PointsB = FromArrayD([]float64(pB), false, len(pB), 1)
	}

	m, n := c.PointsA.Size()

	Y := Zeros(m, n)
	y := ZeroVec(m)

	for i := 0; i < n; i++ {
		pA := c.A.Rot.ApplyTo(c.PointsA.Col(i))
		pB := c.B.Rot.ApplyTo(c.PointsB.Col(i))
		y.Add(pA, c.A.Pos)
		y.Sub(y, pB)
		y.Sub(y, c.B.Pos)
		Y.SetCol(i, y)
	}

	for {
		v = MinPoly(Y)

		vNeg.Neg(v)

		pA := c.A.LinOpt(c.A.Rot.ApplyTo(vNeg))
		pB := c.B.LinOpt(c.B.Rot.ApplyTo(v))

		pArot := c.A.Rot.ApplyTo(pA)
		pBrot := c.B.Rot.ApplyTo(pB)
		y.Add(pArot, c.A.Pos)
		y.Sub(y, pBrot)
		y.Sub(y, c.B.Pos)

		if v.Norm2Sq()+vNeg.Dot(y) < 1e-3 {
			break
		}

		c.PointsA.AddCol(pA)
		c.PointsB.AddCol(pB)
		Y.AddCol(y)
	}

	c.Dist = v.Norm2Sq()
	c.Normal = v.Normalize(v)
	return c
}

func MinPoly(Y *MatD) VecD {
	_, n := Y.Size()
	if n == 1 {
		return Y.Col(0).Copy()
	}
	if n == 2 {
		l := ZeroVec(2)
		d := Y.Col(0).Dot(Y.Col(1))
		l[0] = Y.Col(0).Norm2Sq() - d
		l[1] = Y.Col(1).Norm2Sq() - d
		l.Mul(1/l.Norm1(), l)

		if l[0] <= 0 {
			Y.SetCol(0, nil)
			return Y.Col(0).Copy()
		} else if l[1] <= 0 {
			Y.SetCol(1, nil)
			return Y.Col(0).Copy()
		} else {
			return Y.ApplyTo(l)
		}
	}
	if n == 3 {
		v, l := CheckTri(Y.Col(0), Y.Col(1), Y.Col(2))
		ix := 0
		for la := range l {
			if la == 0 {
				Y.SetCol(ix, nil)
				ix--
			}
			ix++
		}
		return v
	}
	if n == 4 {
		y0 := Y.Col(0)
		y1 := Y.Col(1)
		y2 := Y.Col(2)
		y3 := Y.Col(3)

		y30 := ZeroVec(m).Sub(y3, y0)
		y31 := ZeroVec(m).Sub(y3, y1)
		y32 := ZeroVec(m).Sub(y3, y2)

		y01 := ZeroVec(m).Sub(y0, y1)
		y02 := ZeroVec(m).Sub(y0, y2)

		n0 := ZeroVec(m).Cross(y30, y31)
		if y32.Dot(n0) > 0 {
			n0.Neg(n0)
		}

		n1 := ZeroVec(m).Cross(y30, y32)
		if y31.Dot(n1) > 0 {
			n1.Neg(n1)
		}

		n2 := ZeroVec(m).Cross(y31, y32)
		if y30.Dot(n2) > 0 {
			n2.Neg(n2)
		}

		n3 := ZeroVec(m).Cross(y01, y02)
		if y30.Dot(n2) <= 0 {
			n3.Neg(n3)
		}
	}
	return nil
}

func ProjectOnSimplex(v VecD) VecD {
	n := len(v)
	vSort := v.Copy()
	sort.Float64s([]float64(vSort))

	var wl float64
	for i := range vSort {
		wl = vSort[i]
		S := vSort[i:].Norm1()
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

func closestPoint(A, B *MatD, sA, sB VecD) (a, b VecD) {
	s := 1.0
	n, _ := A.Size()

	d := ZeroVec(n)
	tempA := ZeroVec(len(sA))
	tempB := ZeroVec(len(sB))

	sAhat := ZeroVec(len(sA))
	sBhat := ZeroVec(len(sB))

	ProjectOnSimplex(sA)
	ProjectOnSimplex(sB)

	a = A.ApplyTo(sA)
	b = B.ApplyTo(sB)

	nIter := 0
	for {
		d.Sub(b, a)
		fmt.Println("matrix vector mult")

		gA := A.ApplyToTr(d)
		gB := B.ApplyToTr(d.Neg(d))
		fmt.Println("after")

		if terminationCheck(sA, gA) && terminationCheck(sB, gB) {
			break
		}
		nIter++

		f0 := d.Norm2Sq()
		fPrev := f0

		s = s * 2
		for {
			ProjectOnSimplex(sAhat.Add(sA, tempA.Mul(s, gA)))
			ProjectOnSimplex(sBhat.Add(sB, tempB.Mul(s, gB)))

			fCurr := d.Sub(B.ApplyTo(sBhat),
				A.ApplyTo(sAhat)).Norm2Sq()
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

func LinOptPoly(A *MatD) func(VecD) VecD {
	return func(c VecD) VecD {
		w := A.ApplyToTr(c)
		maxIx := w.MaxIx()
		return A.Col(maxIx).Copy()
	}
}

func terminationCheck(s, g VecD) bool {
	if g.Norm2Sq() < 1e-6 {
		return true
	}
	q0 := ZeroVec(len(s))
	q := ZeroVec(len(s))
	for i := range q0 {
		q0[i] = 1
	}

	var Q *MatD

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
		diff := ZeroVec(m).Sub(gproj, g).Norm2Sq()
		if diff > 1e-3 {
			terminate = false
		}
	}
	return terminate
}
