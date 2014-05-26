//+build ignore

package physics

import . "github.com/dane-unltd/linalg/matrix"

func LinOptSphere(r float64) SupportFunc {
	w := NewVecD(3)
	return func(c VecD) VecD {
		w.Normalize(c)
		w.Scal(r)
		return w
	}
}

func LinOptPoly(A *DenseD) SupportFunc {
	_, n := A.Size()
	w := NewVecD(n)
	A.T()
	return func(c VecD) VecD {
		w.Mul(A, c)
		maxIx := w.Imax()
		return A.ColView(maxIx).Copy().(VecD)
	}
}

func AABB(v VecD) *DenseD {
	v.Scal(0.5)
	A := FromArrayD(v, false, 3, 1)
	for ix := 0; ix < 3; ix++ {
		v[ix] = -v[ix]
		A.AddCol(v)
		v[ix] = -v[ix]
	}
	for ix := 0; ix < 3; ix++ {
		v[ix] = -v[ix]
	}
	A.AddCol(v)
	for ix := 0; ix < 3; ix++ {
		v[ix] = -v[ix]
		A.AddCol(v)
		v[ix] = -v[ix]
	}
	return A
}
