//Copyright 2017 Wanchain crypto team

// This file realize basic mpc functions of finite field. inluding

// ---- Lagrange's polynomial interpolation algorithm
// ---- MPC protocol for addition operation
// ---- MPC protocol for multiplication operation
// ---- MPC protocol for unary inverse operation

package mpc

import (
	"fmt"
	"math/rand"
)

type polynomial []int

// generate a random polynomial, its constant item is nominated
func RandPoly(degree int, constant int) polynomial {

	poly := make(polynomial, degree+1)

	poly[0] = constant

	for i := 1; i < degree+1; i++ {
		rand.Seed(int64(i + 1))
		poly[i] = rand.Intn(100) + 1
	}

	return poly
}

// calculate polynomial's evaluation at some point
func EvaluatePoly(f polynomial, x int) uint64 {

	degree := len(f) - 1

	sum := uint64(0)

	for i := 0; i < degree+1; i++ {
		sum += uint64(f[i]) * pow(x, i)
	}

	return sum
}

func pow(x int, n int) uint64 {
	if n == 0 {
		return 1
	} else {
		return uint64(x) * pow(x, n-1)
	}
}

// calculate the b coefficient in Lagrange's polynomial interpolation algorithm
func evaluateB(x []int) []float64 {

	k := len(x)

	b := make([]float64, k)

	for i := 0; i < k; i++ {
		b[i] = evaluateb(x, i)
	}

	return b
}

// sub-function for evaluateB
func evaluateb(x []int, i int) float64 {

	k := len(x)

	sum := float64(1)

	for j := 0; j < k; j++ {
		if j != i {
			sum *= float64(x[j]) / (float64(x[j]) - float64(x[i]))
		} else {
			continue
		}
	}

	return sum
}

// Lagrange's polynomial interpolation algorithm
func Lagrange(f []uint64, x []int) int {

	degree := len(x) - 1

	b := evaluateB(x)

	fmt.Println("b", b)

	s := float64(0)

	for i := 0; i < degree+1; i++ {

		s += float64(f[i]) * b[i]

	}

	return int(s)
}
