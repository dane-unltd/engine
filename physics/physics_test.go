package physics

import (
	"fmt"
	. "github.com/dane-unltd/linalg/matrix"
	"testing"
)

func Test_RigidBody(t *testing.T) {
	Aa := VecD{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		0, 0, 0,
	}
	Ba := VecD{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		1, 1, 1,
	}

	A := FromArrayD(Aa, true, 3, 4)
	B := FromArrayD(Ba, true, 3, 4)

	rb1 := RigidBody{}
	rb1.Rot = Eye(3)
	rb1.LinOpt = LinOptPoly(A)
	rb1.Pos = NewVecD(3)

	rb2 := RigidBody{}
	rb2.Rot = Eye(3)
	rb2.LinOpt = LinOptPoly(B)
	rb2.Pos = VecD{-1.2, -1.2, -1}

	c := NewContact()
	c.A = &rb1
	c.B = &rb2
	c.Update()

	fmt.Println(c.Normal, c.Dist)

	fmt.Println("done")

}

func Benchmark_RigidBody(b *testing.B) {
	for i := 0; i < b.N; i++ {
		A := RandN(3, 10)
		B := RandN(3, 10)

		rb1 := RigidBody{}
		rb1.Rot = Eye(3)
		rb1.LinOpt = LinOptPoly(A)
		rb1.Pos = NewVecD(3)

		rb2 := RigidBody{}
		rb2.Rot = Eye(3)
		rb2.LinOpt = LinOptPoly(B)
		rb2.Pos = VecD{10, 10, 10}

		c := NewContact()
		c.A = &rb1
		c.B = &rb2
		c.Update()
	}
}
