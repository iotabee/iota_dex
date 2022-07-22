package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestEVMSign(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	data := []byte("1655714635")
	hash := crypto.Keccak256Hash(data)
	signature, _ := crypto.Sign(hash.Bytes(), privateKey)
	t.Log(hexutil.Encode(signature))
}

func TestIOTASign(t *testing.T) {
	privateKey, _ := hex.DecodeString("4f4b376e64ac07fab72e76d79bfe8b958541f366887d3a595dcbe971680f0ad2e30c1f106286bd8f2258d326a91ea3b54c8360f1bc99cbfab512538a88bbd17d")
	data := []byte("1655714635")
	sig := ed25519.Sign(privateKey, data)
	t.Log(hex.EncodeToString(sig))
	if !ed25519.Verify(privateKey[32:], data, sig) {
		t.Error("verify error")
	}
}
