package mpc

import (
	"crypto/ecdsa"
	Rand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/crypto"
)

// this file realized ecdsa signature protocol of mpc version

// it's the basis of the whole wanchain locked account scheme

//original ecdsa signature protocol: sign
func ecdsaSign(m []byte, d ecdsa.PrivateKey) (big.Int, big.Int) {

	curve := crypto.S256()

	source := rand.NewSource(int64(1))

	randReader := rand.New(source)

	k, _ := Rand.Int(randReader, secp256k1_N)

	x1, _ := curve.ScalarBaseMult(k.Bytes()) // kG=(x1,y1)

	r := x1.Mod(x1, secp256k1_N)

	k.ModInverse(k, secp256k1_N)

	e := new(big.Int).SetBytes(crypto.Keccak256(m))

	s := new(big.Int)

	s.Mul(d.D, r)

	s.Add(s, e)

	s.Mul(s, k)

	s.Mod(s, secp256k1_N)

	return *r, *s
}

//original ecdsa signature protocol: verify
func ecdsaVerify(m []byte, D ecdsa.PublicKey, r big.Int, s big.Int) bool {

	curve := crypto.S256()

	e := new(big.Int).SetBytes(crypto.Keccak256(m))

	w := new(big.Int)

	u1 := new(big.Int)

	u2 := new(big.Int)

	w.ModInverse(&s, secp256k1_N)

	u1.Mul(e, w)

	u1.Mod(u1, secp256k1_N)

	u2.Mul(&r, w)

	u2.Mod(u2, secp256k1_N)

	A := new(ecdsa.PublicKey)

	B := new(ecdsa.PublicKey)

	C := new(ecdsa.PublicKey)

	A.X, A.Y = curve.ScalarBaseMult(u1.Bytes())

	B.X, B.Y = curve.ScalarMult(D.X, D.Y, u2.Bytes())

	C.X, C.Y = curve.Add(A.X, A.Y, B.X, B.Y)

	//to do: check whether C is infinite point of secp256k1

	C.X.Mod(C.X, secp256k1_N)

	if r.Cmp(C.X) == 0 {
		return true
	}

	return false
}

// ecdsa signature protocol in mpc version. Key function of locked account scheme
func Mpc_ecdsaSign(m []byte, f []big.Int, x []big.Int) (big.Int, big.Int) {

	if len(f) != len(x) {
		fmt.Errorf("Input length doesn't match!-----Mpc_inverse")
	}

	length := len(f)
	degree := (length+1)/2 - 1

	source := rand.NewSource(int64(1))

	randReader := rand.New(source)

	k := make([]big.Int, length)
	for i := 0; i < length; i++ {
		temp, _ := Rand.Int(randReader, secp256k1_N)
		poly := RandPoly(degree, *temp)
		for j := 0; j < length; j++ {
			temp1 := EvaluatePoly(poly, &x[j])
			k[j].Add(&k[j], &temp1)
			k[j].Mod(&k[j], secp256k1_N)
		}
	}

	R := make([]ecdsa.PublicKey, length)
	for i := 0; i < length; i++ {
		R[i].X, R[i].Y = crypto.S256().ScalarBaseMult(k[i].Bytes())
		R[i].Curve = crypto.S256()
	}

	b := evaluateB(x)

	kG := new(ecdsa.PublicKey)
	kG.X, kG.Y = crypto.S256().ScalarMult(R[0].X, R[0].Y, b[0].Bytes()) //in case the pointer is nil
	for i := 1; i < length; i++ {
		buffer1, buffer2 := crypto.S256().ScalarMult(R[i].X, R[i].Y, b[i].Bytes())
		kG.X, kG.Y = crypto.S256().Add(kG.X, kG.Y, buffer1, buffer2)
	}

	r := *(kG.X).Mod(kG.X, secp256k1_N) //to do: need to check whether r=0

	k_inverse := Mpc_inverse(k, x)

	e := new(big.Int).SetBytes(crypto.Keccak256(m))

	s_share := make([]big.Int, length)

	for i := 0; i < length; i++ { //dr
		f[i].Mul(&f[i], &r)
		f[i].Mod(&f[i], secp256k1_N)
	}
	buffer1 := Mpc_mult(k_inverse, f, x) // k_inverse * dr

	for i := 0; i < length; i++ { //k_inverse * e
		k_inverse[i].Mul(&k_inverse[i], e)
		k_inverse[i].Mod(&k_inverse[i], secp256k1_N)
	}

	for i := 0; i < length; i++ {
		s_share[i].Add(&k_inverse[i], &buffer1[i])
	}

	s := Lagrange(s_share, x)

	return r, s
}

//Mpc version, the same with original ecdsa signature protocol: verify
func ecdsaMpcVerify(m []byte, D ecdsa.PublicKey, r big.Int, s big.Int) bool {

	curve := crypto.S256()

	e := new(big.Int).SetBytes(crypto.Keccak256(m))

	w := new(big.Int)

	u1 := new(big.Int)

	u2 := new(big.Int)

	w.ModInverse(&s, secp256k1_N)

	u1.Mul(e, w)

	u1.Mod(u1, secp256k1_N)

	u2.Mul(&r, w)

	u2.Mod(u2, secp256k1_N)

	A := new(ecdsa.PublicKey)

	B := new(ecdsa.PublicKey)

	C := new(ecdsa.PublicKey)

	A.X, A.Y = curve.ScalarBaseMult(u1.Bytes())

	B.X, B.Y = curve.ScalarMult(D.X, D.Y, u2.Bytes())

	C.X, C.Y = curve.Add(A.X, A.Y, B.X, B.Y)

	//to do: check whether C is infinite point of secp256k1

	C.X.Mod(C.X, secp256k1_N)

	if r.Cmp(C.X) == 0 {
		return true
	}

	return false
}
