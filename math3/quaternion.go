package math3

import "math"

type Quaternion [4]float64

func (res *Mat) RotFromQuat(q *Quaternion) {
	Nq := q.Nrm2Sq()
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

	res[0] = 1.0 - (yY + zZ)
	res[1] = xY + wZ
	res[2] = xZ - wY
	res[3] = xY - wZ
	res[4] = 1.0 - (xX + zZ)
	res[5] = yZ + wX
	res[6] = xZ + wY
	res[7] = yZ - wX
	res[8] = 1.0 - (xX + yY)
}

func (res *Quaternion) Add(a, b *Quaternion) *Quaternion {
	res[0] = a[0] + b[0]
	res[1] = a[1] + b[1]
	res[2] = a[2] + b[2]
	res[3] = a[3] + b[3]
	return res
}

func (res *Quaternion) Sub(a, b *Quaternion) *Quaternion {
	res[0] = a[0] - b[0]
	res[1] = a[1] - b[1]
	res[2] = a[2] - b[2]
	res[3] = a[3] - b[3]
	return res
}

func (res *Quaternion) Scal(a float64) *Quaternion {
	res[0] = a * res[0]
	res[1] = a * res[1]
	res[2] = a * res[2]
	res[3] = a * res[3]
	return res
}

func (q *Quaternion) Nrm2Sq() float64 {
	return q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3]
}

func (q *Quaternion) Nrm2() float64 {
	return math.Sqrt(q.Nrm2Sq())
}

func (res *Quaternion) Normalize(a *Quaternion) *Quaternion {
	alpha := 1 / a.Nrm2()
	res[0] = alpha * a[0]
	res[1] = alpha * a[1]
	res[2] = alpha * a[2]
	res[3] = alpha * a[3]
	return res
}

func (res *Quaternion) Mul(a, b *Quaternion) *Quaternion {
	res[0] = a[0]*b[0] - a[1]*b[1] - a[2]*b[2] - a[3]*b[3]
	res[1] = a[0]*b[1] + a[1]*b[0] + a[2]*b[3] - a[3]*b[2]
	res[2] = a[0]*b[2] + a[2]*b[0] + a[3]*b[1] - a[1]*b[3]
	res[3] = a[0]*b[3] + a[3]*b[0] + a[1]*b[2] - a[2]*b[1]
	return res
}
