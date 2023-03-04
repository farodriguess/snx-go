package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"snxgo/util"
)

type PwEncode struct {
	Modulus  string
	Exponent int
	Testing  bool
	Debug    bool
}

type CustomRand struct {
	testing bool
}

func (t *CustomRand) Read(p []byte) (n int, err error) {
	if t.testing {
		util.FillArrayWithValue(p, 1)
		n = len(p)
	} else {
		n, err = rand.Read(p)
	}
	//replaces all zero bytes
	for i, b := range p {
		if b == 0 {
			p[i] = 1
		}
	}
	return n, err
}

func (pwEnconde *PwEncode) EncodePWD(password string) string {
	pwEnconde.log("EncodePWD start")
	bigN := new(big.Int)
	_, ok := bigN.SetString(pwEnconde.Modulus, 16)
	if !ok {
		panic("failed to parse")
	}
	pub := rsa.PublicKey{
		N: bigN,
		E: pwEnconde.Exponent,
	}
	pwEnconde.log("Public Key: %v", pub)
	pwEnconde.log("pub.Size(): %v", pub.Size())

	pb := util.ReverseArray([]byte(password))
	l := pub.Size()
	r := []byte{}
	r = append(r, 0x0)
	r = append(r, pb...)
	r = append(r, 0x0)
	n := l - len(r) - 2

	if n > 0 {
		randoms := CustomRand{testing: pwEnconde.Testing}
		g := make([]byte, n)
		randoms.Read(g)
		pwEnconde.log("Generated Randoms: %v", g)
		r = append(r, g...)
	}
	r = append(r, 2)
	r = append(r, 0)
	r = util.ReverseArray(r)

	encrypted := new(big.Int)
	e := big.NewInt(int64(pub.E))
	payload := new(big.Int).SetBytes(r)
	encrypted.Exp(payload, e, pub.N)
	encryptedBytes := encrypted.Bytes()
	encryptedBytes = util.ReverseArray(encryptedBytes)

	pwEnconde.log("encryptedBytes: %v", encryptedBytes)
	encodedPWD := hex.EncodeToString(encryptedBytes)
	pwEnconde.log("encodedPWD: %v", encodedPWD)
	pwEnconde.log("EncodePWD end")

	return encodedPWD
}

func (pwEnconde *PwEncode) log(msg string, a ...any) {
	if pwEnconde.Debug {
		fmt.Println(fmt.Sprintf(msg, a...))
	}
}
