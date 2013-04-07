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
	Points   *MatD

	LinOpt func(c VecD) VecD
}

type Contact struct {
	Normal VecD
	Dist   float64
	A, B   *RigidBody
}

func CreateContact(A, B *RigidBody) *Contact {
	/*	gradA := ZeroVec(len(A.Pos))
		gradB := ZeroVec(len(B.Pos))
		gradA.Sub(B.Pos, A.Pos)
		gradB.Sub(A.Pos, B.Pos)

		pA := A.Rot.ApplyTo(A.LinOpt(A.Rot.ApplyTo(gradA)))
		pA.Add(pA, A.Pos)
		pB := B.Rot.ApplyTo(B.LinOpt(B.Rot.ApplyTo(gradB)))
		pB.Add(pB, B.Pos)

		pointsA := FromArrayD([]float64(pA), false, len(pA), 1)
		pointsB := FromArrayD([]float64(pB), false, len(pB), 1)

		alphaA := VecD{1}
		alphaB := VecD{1}

		tempA := ZeroVec(len(pA))
		tempB := ZeroVec(len(pB))

		for {
			gradA.Sub(pB, pA)
			gradB.Sub(pA, pB)

			pANew := A.Rot.ApplyTo(A.LinOpt(A.Rot.ApplyTo(gradA)))
			pANew.Add(pANew, A.Pos)
			pBNew := B.Rot.ApplyTo(B.LinOpt(B.Rot.ApplyTo(gradB)))
			pBNew.Add(pBNew, B.Pos)

			dA := gradA.Dot(tempB.Sub(pANew, pA))
			dB := gradB.Dot(tempA.Sub(pBNew, pB))

			if dA+dB < 0.03 {
				break
			}

			if dA > 0.01 {
				pointsA.AddCol(pANew)
				alphaA = append(alphaA, 1)
			}
			if dB > 0.01 {
				pointsB.AddCol(pBNew)
				alphaB = append(alphaB, 1)
			}

			pA, pB = closestPoint(pointsA, pointsB, alphaA, alphaB)
		}*/

	_, nA := A.Points.Size()
	_, nB := B.Points.Size()

	pointsA := A.Points.Copy()
	pointsB := B.Points.Copy()

	for i := 0; i < nA; i++ {
		pointsA.Col(i).Add(pointsA.Col(i), A.Pos)
	}
	for i := 0; i < nB; i++ {
		pointsB.Col(i).Add(pointsB.Col(i), B.Pos)
	}

	alphaA := OnesVec(nA)
	alphaB := OnesVec(nB)

	pA, pB := closestPoint(pointsA, pointsB, alphaA, alphaB)
	gradA := ZeroVec(3).Sub(pB, pA)

	c := &Contact{}
	c.Dist = gradA.Norm2()
	c.A = A
	c.B = B

	/*	fmt.Println("pA,pB", pA, pB)
		fmt.Println("shares", alphaA, alphaB)
		fmt.Println("points", pointsA, pointsB)
		fmt.Println("grad", gradA)
		fmt.Println("posA", A.Pos)
		fmt.Println("posB", B.Pos)
		fmt.Println("dist", c.Dist)*/

	if c.Dist > 0.01 {
		c.Normal = gradA.Normalize(gradA)
	} else {
		shares := append(alphaA, alphaB...)
		C := ConcatD(A.Points, B.Points)
		mean := C.ApplyTo(shares)
		_, nC := C.Size()
		for i := 0; i < nC; i++ {
			C.Col(i).Sub(C.Col(i), mean)
		}
		//C = MulD(C, Diag(shares))

		Ct := C.Copy()
		Ct.Tr()
		CtC := MulD(C, Ct)
		V, d := CtC.EigSy()
		ixMin := d.MinIx()
		c.Normal = V.Col(ixMin)
		fmt.Println("shares", shares)
		fmt.Println("mean", mean)
		fmt.Println("V", V)
		fmt.Println("pA,pB", pA, pB)
		fmt.Println("shares", alphaA, alphaB)
		fmt.Println("points", pointsA, pointsB)
		fmt.Println("grad", gradA)
		fmt.Println("posA", A.Pos)
		fmt.Println("posB", B.Pos)
		fmt.Println("dist", c.Dist)
	}
	return c
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

		gA := A.ApplyToTr(d)
		gB := B.ApplyToTr(d.Neg(d))

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
