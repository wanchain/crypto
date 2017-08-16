package Anonymous

import (
	"crypto/cipher"
	"crypto/ecdsa"
	Rand "crypto/rand"
	"crypto/rsa"
	"io"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

//Shi: generate a random num in [1,N]. N is the order of group [G].
func randFieldElement(rand io.Reader) (k *big.Int, err error) {
	params := crypto.S256().Params()
	b := make([]byte, params.BitSize/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return
	}
	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)

	return
}

//Shi:calc keyimage in ringsignature I=[x]Hash(P)
func xScalarHashP(x []byte, pub *ecdsa.PublicKey) (I *ecdsa.PublicKey) {
	HashP := new(ecdsa.PublicKey)
	I = new(ecdsa.PublicKey)
	//calc Hash(P) to get a random point on curve
	HashP.X, HashP.Y = crypto.S256().ScalarBaseMult(crypto.Keccak256(crypto.FromECDSAPub(pub))) 
	I.X, I.Y = crypto.S256().ScalarMult(HashP.X, HashP.Y, x) //I=[x]Hash(P)
	I.Curve = crypto.S256()
	return
}

//Shi: RingSignature with message M, privatekey x, publickeyset publickeys as input and publickeys, keyimage and two random bigint array.
//     The real publickey belonging to signer is PublicKeys[0] and will be exchanged to PublicKeys[s].
func RingSign(M []byte, x *big.Int, PublicKeys []*ecdsa.PublicKey) ([]*ecdsa.PublicKey, *ecdsa.PublicKey, []*big.Int, []*big.Int) {

	n := len(PublicKeys)
	I := xScalarHashP(x.Bytes(), PublicKeys[0]) //calc keyimage
	s := rand.Intn(n)                           //random num to determin the location of the real publickey
	if s > 0 {
		PublicKeys[0], PublicKeys[s] = PublicKeys[s], PublicKeys[0] //exchange the location
	}

	var (
		q = make([]*big.Int, n)
		w = make([]*big.Int, n)
	)
	SumC := new(big.Int).SetInt64(0)
	Lpub := new(ecdsa.PublicKey)  
	Rpub := new(ecdsa.PublicKey)  //Li, Ri in RingSignature Algorithm
	d := sha3.NewKeccak256()

	//calc Hash(M,Li,Ri)
	//firstly Hash(M)
	d.Write(M)  
	//then Hash(M,Li)
	for i := 0; i < n; i++ {
		q[i], _ = randFieldElement(Rand.Reader)
		w[i], _ = randFieldElement(Rand.Reader)
		//if i = s, Li=[qi]G
		Lpub.X, Lpub.Y = crypto.S256().ScalarBaseMult(q[i].Bytes())
		//if i != s, Li=[qi]G + [wi]Pi
		if i != s {
			Ppub := new(ecdsa.PublicKey)
			Ppub.X, Ppub.Y = crypto.S256().ScalarMult(PublicKeys[i].X, PublicKeys[i].Y, w[i].Bytes()) //[wi]Pi
			Lpub.X, Lpub.Y = crypto.S256().Add(Lpub.X, Lpub.Y, Ppub.X, Ppub.Y)                        //[qi]G+[wi]Pi

			SumC.Add(SumC, w[i])
			SumC.Mod(SumC, secp256k1_N)
		}

		d.Write(crypto.FromECDSAPub(Lpub))
	}
	//then Hash(M,Li,Ri)
	for i := 0; i < n; i++ {
		//if i = s, Ri=[qi]Hash(Pi)
		Rpub = xScalarHashP(q[i].Bytes(), PublicKeys[i]) //[qi]Hash(Pi)
		//if i != s, Ri=[qi]Hash(Pi)+[wi]I
		if i != s {
			Ppub := new(ecdsa.PublicKey)
			Ppub.X, Ppub.Y = crypto.S256().ScalarMult(I.X, I.Y, w[i].Bytes())  //[wi]I
			Rpub.X, Rpub.Y = crypto.S256().Add(Rpub.X, Rpub.Y, Ppub.X, Ppub.Y) //[qi]Hash(Pi)+[wi]I
		}

		d.Write(crypto.FromECDSAPub(Rpub))
	}
	Cs := new(big.Int).SetBytes(d.Sum(nil)) //hash(m,Li,Ri)

	Cs.Sub(Cs, SumC)
	Cs.Mod(Cs, secp256k1_N)

	tmp := new(big.Int).Mul(Cs, x)
	Rs := new(big.Int).Sub(q[s], tmp)
	Rs.Mod(Rs, secp256k1_N)
	w[s] = Cs
	q[s] = Rs

	return PublicKeys, I, w, q
}

//Shi: Verify the RingSignature
func VerifyRingSign(M []byte, PublicKeys []*ecdsa.PublicKey, I *ecdsa.PublicKey, c []*big.Int, r []*big.Int) bool {
	ret := false
	n := len(PublicKeys)
	SumC := new(big.Int).SetInt64(0)
	Lpub := new(ecdsa.PublicKey)
	d := sha3.NewKeccak256()
	d.Write(M)
	//hash(M,Li,Ri)
	for i := 0; i < n; i++ {
		Lpub.X, Lpub.Y = crypto.S256().ScalarBaseMult(r[i].Bytes()) //[ri]G

		Ppub := new(ecdsa.PublicKey)
		Ppub.X, Ppub.Y = crypto.S256().ScalarMult(PublicKeys[i].X, PublicKeys[i].Y, c[i].Bytes()) //[ci]Pi
		Lpub.X, Lpub.Y = crypto.S256().Add(Lpub.X, Lpub.Y, Ppub.X, Ppub.Y)                        //[ri]G+[ci]Pi
		SumC.Add(SumC, c[i])
		SumC.Mod(SumC, secp256k1_N)
		d.Write(crypto.FromECDSAPub(Lpub))
	}
	Rpub := new(ecdsa.PublicKey)
	for i := 0; i < n; i++ {
		Rpub = xScalarHashP(r[i].Bytes(), PublicKeys[i]) //[qi]Hash(Pi)
		Ppub := new(ecdsa.PublicKey)
		Ppub.X, Ppub.Y = crypto.S256().ScalarMult(I.X, I.Y, c[i].Bytes())  //[wi]I
		Rpub.X, Rpub.Y = crypto.S256().Add(Rpub.X, Rpub.Y, Ppub.X, Ppub.Y) //[qi]Hash(Pi)+[wi]I

		d.Write(crypto.FromECDSAPub(Rpub))
	}
	hash := new(big.Int).SetBytes(d.Sum(nil)) //hash(m,Li,Ri)
	hash.Mod(hash, secp256k1_N)
	if hash.Cmp(SumC) == 0 {
		ret = true
	}
	return ret
}




