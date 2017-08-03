//Copyright 2017 Wanchain crypto team

// This file realizes Lagrange's polynomial interpolation algorithm over finite field

package mpc

import (
	Rand "crypto/rand"
	"math/big"
	"math/rand"
	//"wanchain/MPC/secp256k1"
	//"wanchain/MPC/math"
)

// generate a random polynomial, its constant item is nominated
func RandPoly(degree int, constant big.Int) polynomial {

	poly := make(polynomial, degree+1)

	poly[0] = constant

	temp := new(big.Int)

	for i := 1; i < degree+1; i++ {
		source := rand.NewSource(int64(i))
		r := rand.New(source)
		temp, _ = Rand.Int(r, secp256k1_N)
		// in case of polynomial degenerating
		poly[i] = *temp.Add(temp, bigOne)
	}

	return poly
}

// calculate polynomial's evaluation at some point
func EvaluatePoly(f polynomial, x *big.Int) big.Int {

	degree := len(f) - 1

	sum := big.NewInt(0)

	temp1 := big.NewInt(1)
	temp2 := big.NewInt(1)

	for i := 0; i < degree+1; i++ {
		temp1.Exp(x, big.NewInt(int64(i)), secp256k1_N)
		temp2.Mul(&f[i], temp1)
		sum.Add(sum, temp2)
		sum.Mod(sum, secp256k1_N)
	}

	return *sum
}

// calculate the b coefficient in Lagrange's polynomial interpolation algorithm
func evaluateB(x []big.Int) []big.Int {

	k := len(x)

	b := make([]big.Int, k)

	for i := 0; i < k; i++ {
		b[i] = evaluateb(x, i)
	}

	return b
}

// sub-function for evaluateB
func evaluateb(x []big.Int, i int) big.Int {

	k := len(x)

	sum := big.NewInt(1)

	temp1 := big.NewInt(1)
	temp2 := big.NewInt(1)

	for j := 0; j < k; j++ {
		if j != i {
			temp1.Sub(&x[j], &x[i])
			temp1.ModInverse(temp1, secp256k1_N)
			temp2.Mul(&x[j], temp1)
			sum.Mul(sum, temp2)
			sum.Mod(sum, secp256k1_N)
		} else {
			continue
		}
	}

	return *sum
}

// Lagrange's polynomial interpolation algorithm
func Lagrange(f []big.Int, x []big.Int) big.Int {

	degree := len(x) - 1

	b := evaluateB(x)

	//fmt.Println("b", b)

	s := big.NewInt(0)

	temp1 := big.NewInt(1)

	for i := 0; i < degree+1; i++ {
		temp1.Mul(&f[i], &b[i])
		s.Add(s, temp1)
		s.Mod(s, secp256k1_N)
	}

	return *s
}

// calculate the inverse of a element over finite field
func modInverse(a, n *big.Int) (*big.Int, bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	// a n not coprime
	if g.Cmp(bigOne) != 0 {
		return bigOne, false
	}
	// when x is negative
	if x.Cmp(bigOne) < 0 {
		x.Add(x, n)
	}
	return x, true
}
