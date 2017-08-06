package mpc

import (
	"crypto/ecdsa"
	Rand "crypto/rand"
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
