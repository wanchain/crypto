// this file realizes the basic functions of mpc, including

// ---- MPC protocol for addition operation
// ---- MPC protocol for multiplication operation
// ---- MPC protocol for unary inverse operation

package mpc

import (
	Rand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
)

// based on mpc protocol for add operation in wanchain white paper
func Mpc_add(x []big.Int, y []big.Int) []big.Int {
	if len(x) != len(y) {
		fmt.Errorf("Input len doesn't match!-----Mpc_add")
		return nil
	}

	z := make([]big.Int, len(x))

	for i := 0; i < len(x); i++ {
		z[i].Add(&x[i], &y[i])
	}
	return z
}

func Mpc_mult(f1 []big.Int, f2 []big.Int, x []big.Int) []big.Int {
	if len(f1) != len(x) {
		fmt.Errorf("Input length doesn't match!-----Mpc_mult")
		return nil
	}
	length := len(f1)
	k := (length+1)/2 - 1
	b := evaluateB(x)
	result := make([]big.Int, length)
	temp := big.NewInt(0)
	temp1 := big.NewInt(0)

	for i := 0; i < length; i++ {
		f := make([]big.Int, length)
		poly := RandPoly(k, *temp.Mul(&f1[i], &f2[i]))
		for j := 0; j < length; j++ {
			f[j] = EvaluatePoly(poly, &x[j])
			result[j].Add(&result[j], temp1.Mul(&b[i], &f[j]))
		}
	}

	return result
}

func Mpc_inverse(f []big.Int, x []big.Int) []big.Int {

	if len(f) != len(x) {
		fmt.Errorf("Input length doesn't match!-----Mpc_inverse")
	}
	length := len(f)
	k := (length+1)/2 - 1
	source := rand.NewSource(int64(length))
	t := rand.New(source)
	r := make([]big.Int, length)

	for i := 0; i < length; i++ {
		temp, _ := Rand.Int(t, secp256k1_N)
		poly := RandPoly(k, *temp)
		for j := 0; j < length; j++ {
			temp1 := EvaluatePoly(poly, &x[j])
			r[j].Add(&r[j], &temp1)
		}
	}
	ar_share := Mpc_mult(f, r, x)
	ar := Lagrange(ar_share, x)

	ar.ModInverse(&ar, secp256k1_N)

	for i := 0; i < length; i++ {
		r[i].Mul(&r[i], &ar)
	}
	return r
}
