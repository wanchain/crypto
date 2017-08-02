package mpc

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {

	secret := int(100)

	degree := int(5)

	p := RandPoly(degree, secret)

	fmt.Println("Poly ", p)

	f := make([]uint64, degree+1)

	x := make([]int, degree+1)

	for i := 0; i < degree+1; i++ {
		x[i] = i + 1
	}

	fmt.Println("x ", x)

	for i := 0; i < degree+1; i++ {
		f[i] = EvaluatePoly(p, x[i])
	}

	fmt.Println("f ", f)

	result := Lagrange(f, x)

	fmt.Println(result)

}
