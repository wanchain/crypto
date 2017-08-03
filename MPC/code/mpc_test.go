package mpc

import (
	"fmt"
	"math/big"
	"testing"
)

func TestLagrange(t *testing.T) {

	secret := big.NewInt(10)

	degree := int(100)

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

func TestAdd(t *testing.T) {

	secret1 := big.NewInt(100)

	secret2 := big.NewInt(200)

	degree := int(1000)

	p1 := RandPoly(degree, *secret1)

	p2 := RandPoly(degree, *secret2)

	f1 := make([]big.Int, degree+1)

	f2 := make([]big.Int, degree+1)

	x := make([]big.Int, degree+1)

	for i := 0; i < degree+1; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	for i := 0; i < degree+1; i++ {
		f1[i] = EvaluatePoly(p1, &x[i])
	}

	for i := 0; i < degree+1; i++ {
		f2[i] = EvaluatePoly(p2, &x[i])
	}

	z := Mpc_add(f1, f2)

	result := Lagrange(z, x)

	fmt.Println("result ", result)

}
