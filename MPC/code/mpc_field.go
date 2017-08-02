//Copyright 2017 Wanchain crypto team

// This file realize basic mpc functions of finite field. inluding

// ---- Lagrange's polynomial interpolation algorithm
// ---- MPC protocol for addition operation
// ---- MPC protocol for multiplication operation
// ---- MPC protocol for unary inverse operation

package mpc

import (
	"math/rand"
)

type polynomial []int

// generate a random polynomial, its constant item is nominated
func RandPoly(degree int, constant int) polynomial {

	poly := make(polynomial, degree+1)

	poly[0] = constant

	for i := 1; i < degree+1; i++ {
		poly[i] = rand.Int()
	}

	return poly
}

// calculate polynomial's evaluation at some point
func EvaluatePoly(f polynomial, x int) int {

	degree := len(f) - 1

	sum := 0

	for i := 0; i < degree+1; i++ {
		sum += f[i] * pow(x, i)
	}

	return sum
}

func pow(x int, n int) int {
	if n == 0 {
		return 1
	} else {
		return x * pow(x, n-1)
	}
}

func evaluateB(x []int) []float32 {

	k := len(x)

	b := make([]float32, k)

	for i := 0; i < k; i++ {
		b[i] = evaluateb(x, i)
	}

	return b
}

func evaluateb(x []int, i int) float32 {

	k := len(x)

	sum := float32(0)

	for j := 0; j < k; j++ {
		if j != i {
			sum += (float32)(x[j]) / (float32)(x[i]-x[j])
		} else {
			continue
		}
	}

	return sum
}
