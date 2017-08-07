package mpc

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
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

	degree := int(10)

	len := degree + 1

	p1 := RandPoly(degree, *secret1)

	p2 := RandPoly(degree, *secret2)

	f1 := make([]big.Int, len)

	f2 := make([]big.Int, len)

	x := make([]big.Int, len)

	for i := 0; i < len; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	for i := 0; i < len; i++ {
		f1[i] = EvaluatePoly(p1, &x[i])
	}

	for i := 0; i < len; i++ {
		f2[i] = EvaluatePoly(p2, &x[i])
	}

	z := Mpc_add(f1, f2)

	result := Lagrange(z, x)

	fmt.Println("result ", result)

}

func TestMult(t *testing.T) {

	secret1 := big.NewInt(34)

	secret2 := big.NewInt(2)

	degree := int(10)

	len := (degree+1)*2 - 1

	p1 := RandPoly(degree, *secret1)

	p2 := RandPoly(degree, *secret2)

	f1 := make([]big.Int, len)

	f2 := make([]big.Int, len)

	x := make([]big.Int, len)

	for i := 0; i < len; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	for i := 0; i < len; i++ {
		f1[i] = EvaluatePoly(p1, &x[i])
	}

	for i := 0; i < len; i++ {
		f2[i] = EvaluatePoly(p2, &x[i])
	}

	z := Mpc_mult(f1, f2, x)

	result := Lagrange(z, x)

	fmt.Println("result ", result)

}

func TestInverse(t *testing.T) {

	secret := big.NewInt(34)

	degree := int(10)

	len := (degree+1)*2 - 1

	p := RandPoly(degree, *secret)

	f := make([]big.Int, len)

	x := make([]big.Int, len)

	for i := 0; i < len; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	for i := 0; i < len; i++ {
		f[i] = EvaluatePoly(p, &x[i])
	}

	inverse_share := Mpc_inverse(f, x)

	secret_inverse := Lagrange(inverse_share, x)

	secret_inverse.ModInverse(&secret_inverse, secp256k1_N)

	fmt.Println("secret ", secret_inverse)

}

func TestEcdsa(t *testing.T) {

	m := []byte{0x22, 0x33}

	d, _ := crypto.GenerateKey()

	r, s := ecdsaSign(m, *d)

	if ecdsaVerify(m, d.PublicKey, r, s) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}

}

func TestMpcEcdsa(t *testing.T) {

	m := []byte{0x22, 0x33}

	d, _ := crypto.GenerateKey()

	degree := int(10)

	len := (degree+1)*2 - 1

	p := RandPoly(degree, *d.D)

	f := make([]big.Int, len)

	x := make([]big.Int, len)

	for i := 0; i < len; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	for i := 0; i < len; i++ {
		f[i] = EvaluatePoly(p, &x[i])
	}

	r, s := Mpc_ecdsaSign(m, f, x)

	if ecdsaMpcVerify(m, d.PublicKey, r, s) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}
