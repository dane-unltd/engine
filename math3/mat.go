package math3

type Mat [9]float64

func (res *Mat) Mul(A, B *Mat) {
	res[0] = A[0]*B[0] + A[3]*B[1] + A[6]*B[2]
	res[1] = A[1]*B[0] + A[4]*B[1] + A[7]*B[2]
	res[2] = A[2]*B[0] + A[5]*B[1] + A[8]*B[2]

	res[3] = A[0]*B[3] + A[3]*B[4] + A[6]*B[5]
	res[4] = A[1]*B[3] + A[4]*B[4] + A[7]*B[5]
	res[5] = A[2]*B[3] + A[5]*B[4] + A[8]*B[5]

	res[6] = A[0]*B[6] + A[3]*B[7] + A[6]*B[8]
	res[7] = A[1]*B[6] + A[4]*B[7] + A[7]*B[8]
	res[8] = A[2]*B[6] + A[5]*B[7] + A[8]*B[8]
}

func (res *Vec) Apply(A *Mat, b *Vec) {
	res[0] = A[0]*b[0] + A[3]*b[1] + A[6]*b[2]
	res[1] = A[1]*b[0] + A[4]*b[1] + A[7]*b[2]
	res[2] = A[2]*b[0] + A[5]*b[1] + A[8]*b[2]
}

func (A *Mat) T() {
	A[1], A[3] = A[3], A[1]
	A[2], A[6] = A[6], A[2]
	A[5], A[7] = A[7], A[5]
}
