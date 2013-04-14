package physics

import . "github.com/dane-unltd/linalg/matrix"

type Quaternion []float64

func NewQuaternion(comps ...float64) Quaternion {
	q := make(Quaternion, 4)
	copy(q, comps)
	return q
}

func RotFromQuat(q Quaternion) *DenseD {
	Nq := VecD(q).Nrm2Sq()
	s := 0.0
	if Nq > 0.0 {
		s = 2 / Nq
	}

	X := q[1] * s
	Y := q[2] * s
	Z := q[3] * s
	wX := q[0] * X
	wY := q[0] * Y
	wZ := q[0] * Z
	xX := q[1] * X
	xY := q[1] * Y
	xZ := q[1] * Z
	yY := q[2] * Y
	yZ := q[2] * Z
	zZ := q[3] * Z
	RotArr := []float64{1.0 - (yY + zZ), xY + wZ, xZ - wY,
		xY - wZ, 1.0 - (xX + zZ), yZ + wX,
		xZ + wY, yZ - wX, 1.0 - (xX + yY),
	}
	return FromArrayD(RotArr, true, 3, 3)
}

func (res Quaternion) Add(a, b Quaternion) Quaternion {
	VecD(res).Add(VecD(a), VecD(b))
	return res
}
func (res Quaternion) Sub(a, b Quaternion) Quaternion {
	VecD(res).Sub(VecD(a), VecD(b))
	return res
}
func (res Quaternion) Scal(a float64) Quaternion {
	VecD(res).Scal(a)
	return res
}
func (res Quaternion) Axpy(a float64, x Quaternion) Quaternion {
	VecD(res).Axpy(a, VecD(x))
	return res
}
func (res Quaternion) Normalize(a Quaternion) Quaternion {
	VecD(res).Normalize(VecD(a))
	return res
}

func (res Quaternion) Mul(a, b Quaternion) Quaternion {
	res[0] = a[0]*b[0] - a[1]*b[1] - a[2]*b[2] - a[3]*b[3]
	res[1] = a[0]*b[1] + a[1]*b[0] + a[2]*b[3] - a[3]*b[2]
	res[2] = a[0]*b[2] + a[2]*b[0] + a[3]*b[1] - a[1]*b[3]
	res[3] = a[0]*b[3] + a[3]*b[0] + a[1]*b[2] - a[2]*b[1]
	return res
}
