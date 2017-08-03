package mpc

// This file defines the constant num used

// ---secp256k1_N is the order of finite field

import "math/big"

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)
var secp256k1_N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
