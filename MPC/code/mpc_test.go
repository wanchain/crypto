package mpc

import (
	"crypto/ecdsa"
	Rand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
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

func TestAll(t *testing.T) {

	m := []byte{0x22, 0x33}

	degree := int(10)

	length := (degree+1)*2 - 1

	source := rand.NewSource(int64(1))

	randReader := rand.New(source)

	x := make([]big.Int, length)
	for i := 0; i < length; i++ {
		x[i] = *big.NewInt(int64(i + 1))
	}

	f := make([]big.Int, length)
	for i := 0; i < length; i++ {
		temp, _ := Rand.Int(randReader, secp256k1_N)
		poly := RandPoly(degree, *temp)
		for j := 0; j < length; j++ {
			temp1 := EvaluatePoly(poly, &x[j])
			f[j].Add(&f[j], &temp1)
			f[j].Mod(&f[j], secp256k1_N)
		}
	}

	ff := Lagrange(f, x)
	tt := new(ecdsa.PublicKey)
	tt.X, tt.Y = crypto.S256().ScalarBaseMult(ff.Bytes())
	fmt.Println("tt.x", *tt.X)
	fmt.Println("tt.y", *tt.Y)

	r, s := Mpc_ecdsaSign(m, f, x)

	R := make([]ecdsa.PublicKey, length)
	for i := 0; i < length; i++ {
		f[i].Mod(&f[i], secp256k1_N)
		R[i].X, R[i].Y = crypto.S256().ScalarBaseMult(f[i].Bytes())
		R[i].Curve = crypto.S256()
	}

	b := evaluateB(x)

	kG := new(ecdsa.PublicKey)
	kG.X, kG.Y = crypto.S256().ScalarMult(R[0].X, R[0].Y, b[0].Bytes()) //in case the pointer is nil
	for i := 1; i < length; i++ {
		buffer1, buffer2 := crypto.S256().ScalarMult(R[i].X, R[i].Y, b[i].Bytes())
		kG.X, kG.Y = crypto.S256().Add(kG.X, kG.Y, buffer1, buffer2)
	}

	kG.Curve = crypto.S256()

	fmt.Println("kG.x", *kG.X)
	fmt.Println("kG.y", *kG.Y)

	if ecdsaMpcVerify(m, *kG, r, s) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}
