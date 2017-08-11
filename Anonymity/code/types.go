package Anonymous

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"go-ethereum-master/accounts"
	"go-ethereum-master/common"
	"go-ethereum-master/crypto"

	"uuid"
)

///////////////////////////////////////////////////////Key私钥结构//////////////////////////////////////////////////////////

//TeemoGuo: Key struct is both for normal Key and one-time-key, the difference is that for one-time-key PrivateKey2.D=0
type Key struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext

	//TeemoGuo:in our potocol, Address above is derived from  PrivateKey, not PrivateKey2
	PrivateKey *ecdsa.PrivateKey

	PrivateKey2 *ecdsa.PrivateKey
}

//TeemoGuo:to be continue
type plainKeyJSON struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privatekey"`
	PrivateKey2 string `json:"privatekey2"`
	Id          string `json:"id"`
	Version     int    `json:"version"`
}

//TeemoGuo:to be continue
func (k *Key) MarshalJSON() (j []byte, err error) {
	jStruct := plainKeyJSON{
		hex.EncodeToString(k.Address[:]),
		hex.EncodeToString(crypto.FromECDSA(k.PrivateKey)),
		hex.EncodeToString(crypto.FromECDSA(k.PrivateKey2)),
		k.Id.String(),
		version,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *Key) UnmarshalJSON(j []byte) (err error) {
	keyJSON := new(plainKeyJSON)
	err = json.Unmarshal(j, &keyJSON)
	if err != nil {
		return err
	}

	u := new(uuid.UUID)
	*u = uuid.Parse(keyJSON.Id)
	k.Id = *u
	addr, err := hex.DecodeString(keyJSON.Address)
	if err != nil {
		return err
	}

	privkey, err := hex.DecodeString(keyJSON.PrivateKey)
	if err != nil {
		return err
	}
	privkey2, err := hex.DecodeString(keyJSON.PrivateKey2)
	if err != nil {
		return err
	}

	k.Address = common.BytesToAddress(addr)
	k.PrivateKey, _ = crypto.ToECDSA(privkey)
	k.PrivateKey2, _ = crypto.ToECDSA(privkey2)

	return nil
}

///////////////////////////////////////////////////////Account//////////////////////////////////////////////////////////
type Account struct {
	Address common.Address   `json:"address"` // Ethereum account address derived from the key
	URL     accounts.URL     `json:"url"`     // Optional resource locator within a backend
	A       *ecdsa.PublicKey //TeemoGuo add
	B       *ecdsa.PublicKey //TeemoGuo add
}
