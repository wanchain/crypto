// this file realizes the basic functions of mpc, including

// ---- MPC protocol for addition operation
// ---- MPC protocol for multiplication operation
// ---- MPC protocol for unary inverse operation

package mpc

import (
	"fmt"
	"math/big"
)

// based on mpc protocol for add operation in wanchain white paper
func Mpc_add(x []big.Int, y []big.Int) []big.Int {
	if len(x) != len(y) {
		fmt.Errorf("Input length doesn't match!-----Mpc_add")
		return nil
	}

	z := make([]big.Int, len(x))

	for i := 0; i < len(x); i++ {
		z[i].Add(&x[i], &y[i])
	}
	return z
}
