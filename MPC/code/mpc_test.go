package mpc

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAll(t *testing.T) {

	secret := big.NewInt(10)

	degree := int(1000)

	p := RandPoly(degree, *secret)

	//fmt.Println("Poly ", p)

	f := make([]big.Int, degree+1)

	x := make([]big.Int, degree+1)

	for i := 0; i < degree+1; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	//fmt.Println("x ", x)

	for i := 0; i < degree+1; i++ {
		f[i] = EvaluatePoly(p, &x[i])
	}

	//fmt.Println("f ", f)

	result := Lagrange(f, x)

	fmt.Println("result ", result)

}
