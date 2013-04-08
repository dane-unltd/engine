package physics

import (
	. "github.com/dane-unltd/linalg/matrix"
)

func CheckTri(y0, y1, y2 VecD) (v, lambda VecD) {
	lambda = ZeroVec(3)

	m := len(y0)

	v = ZeroVec(m)
	temp := ZeroVec(m)

	y01 := ZeroVec(m)
	y01.Sub(y0, y1)
	y02 := ZeroVec(m)
	y02.Sub(y0, y2)
	y12 := ZeroVec(m)
	y12.Sub(y1, y2)

	normal := ZeroVec(m)
	normal.Cross(y01, y02)

	normal01 := ZeroVec(m)
	normal01.Cross(y01, normal)
	normal02 := ZeroVec(m)
	normal02.Cross(normal, y02)
	normal12 := ZeroVec(m)
	normal12.Cross(y12, normal)

	inside := true
	d0d01 := y0.Dot(y01)
	d1d01 := y1.Dot(y01)
	if normal01.Dot(y0) <= 0 {
		inside = false
		if d0d01 <= 0 {
			if d1d01 >= 0 {
				//edge 01
				lambda[0] = -d0d01
				lambda[1] = d1d01
				lambda.Mul(1/lambda.Norm1(), lambda)
				v.Mul(lambda[0], y0)
				temp.Mul(lambda[1], y1)
				v.Add(v, temp)
				return
			}
		}
	}
	d0d02 := y0.Dot(y02)
	d2d02 := y2.Dot(y02)
	if normal02.Dot(y0) <= 0 {
		inside = false
		if d0d02 <= 0 {
			if d2d02 >= 0 {
				//edge 02
				lambda[0] = -d0d02
				lambda[2] = d2d02
				lambda.Mul(1/lambda.Norm1(), lambda)
				v.Mul(lambda[0], y0)
				temp.Mul(lambda[2], y2)
				v.Add(v, temp)
				return
			}
		}
	}
	d1d12 := y1.Dot(y12)
	d2d12 := y2.Dot(y12)
	if normal12.Dot(y1) <= 0 {
		inside = false
		if d1d12 <= 0 {
			if d2d12 >= 0 {
				//edge 12
				lambda[1] = -d1d12
				lambda[2] = d2d12
				lambda.Mul(1/lambda.Norm1(), lambda)
				v.Mul(lambda[1], y1)
				temp.Mul(lambda[2], y2)
				v.Add(v, temp)
				return
			}
		}
	}
	if inside {
		//in triangle
		lambda = OnesVec(m)
		normal.Normalize(normal)
		v.Mul(y0.Dot(normal), normal)
		return
	}

	if (d0d01 <= 0) && (d0d02 <= 0) {
		lambda[0] = 1
		v = y0.Copy()
		return
	}
	if (d1d01 >= 0) && (d1d12 <= 0) {
		lambda[1] = 1
		v = y1.Copy()
		return
	}
	if (d2d02 >= 0) && (d2d12 >= 0) {
		lambda[2] = 1
		v = y2.Copy()
		return
	}
	return
}
