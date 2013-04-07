package physics

import (
	"fmt"
	_ "github.com/dane-unltd/linalg/blasinit"
	_ "github.com/dane-unltd/linalg/lapackinit"
	. "github.com/dane-unltd/linalg/matrix"
	"testing"
)

func Test_RigidBody(t *testing.T) {
	A := RandN(3, 110)
	B := RandN(3, 110)

	rb1 := RigidBody{}
	rb1.Rot = Eye(3)
	rb1.Points = B
	rb1.LinOpt = LinOptPoly(B)
	rb1.Pos = ZeroVec(3)

	rb2 := RigidBody{}
	rb2.Rot = Eye(3)
	rb2.LinOpt = LinOptPoly(A)
	rb2.Points = A
	rb2.Pos = OnesVec(3)
	rb2.Pos.Mul(2, rb2.Pos)

	c := CreateContact(&rb1, &rb2)
	fmt.Println(c)

	fmt.Println("done")

}

func Benchmark_RigidBody(b *testing.B) {
	for i := 0; i < b.N; i++ {
		A := RandN(3, 10)
		B := RandN(3, 10)

		rb1 := RigidBody{}
		rb1.Rot = Eye(3)
		rb1.LinOpt = LinOptPoly(A)
		rb1.Points = A
		rb1.Pos = ZeroVec(3)

		rb2 := RigidBody{}
		rb2.Rot = Eye(3)
		rb2.LinOpt = LinOptPoly(B)
		rb2.Points = B
		rb2.Pos = OnesVec(3)
		rb2.Pos.Mul(10, rb2.Pos)

		CreateContact(&rb1, &rb2)
	}
}
