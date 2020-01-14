// This vrf package makes the VRF API in Algorand's libsodium C library available to golang.
package vrf

/*
#cgo CFLAGS: -Wall -std=c99
#cgo CFLAGS: -I./include/
#cgo LDFLAGS: -L./lib -lsodium
#include "sodium.h"
*/
import "C"
import (
	"encoding/hex"
	"errors"
	"fmt"
	"unsafe"
)

const (
	PUBLICKEYBYTES = uint32(C.crypto_vrf_PUBLICKEYBYTES)
	SECRETKEYBYTES = uint32(C.crypto_vrf_SECRETKEYBYTES)
	SEEDBYTES = uint32(C.crypto_vrf_SEEDBYTES)
	PROOFBYTES = uint32(C.crypto_vrf_PROOFBYTES)
	OUTPUTBYTES = uint32(C.crypto_vrf_OUTPUTBYTES)
	PRIMITIVE = C.crypto_vrf_PRIMITIVE
)

// Generate an Ed25519 key pair for use with VRF.
func KeyPair() (*[PUBLICKEYBYTES]byte, *[SECRETKEYBYTES]byte) {
	publicKey := [PUBLICKEYBYTES]byte{}
	privateKey := [SECRETKEYBYTES]byte{}
	publicKeyPtr := (*C.uchar)(unsafe.Pointer(&publicKey))
	privateKeyPtr := (*C.uchar)(unsafe.Pointer(&privateKey))
	C.crypto_vrf_keypair(publicKeyPtr, privateKeyPtr)
	return &publicKey, &privateKey
}

// Generate an Ed25519 key pair for use with VRF. Parameter `seed` means the cofactor in Curve25519 and EdDSA.
func KeyPairFromSeed(seed *[SEEDBYTES]byte) (*[PUBLICKEYBYTES]byte, *[SECRETKEYBYTES]byte) {
	publicKey := [PUBLICKEYBYTES]byte{}
	privateKey := [SECRETKEYBYTES]byte{}
	publicKeyPtr := (*C.uchar)(unsafe.Pointer(&publicKey))
	privateKeyPtr := (*C.uchar)(unsafe.Pointer(&privateKey))
	seedPtr := (*C.uchar)(unsafe.Pointer(seed))
	C.crypto_vrf_keypair_from_seed(publicKeyPtr, privateKeyPtr, seedPtr)
	return &publicKey, &privateKey
}

// Verifies that the specified public key is valid.
func IsValidKey(publicKey *[PUBLICKEYBYTES]byte) bool {
	publicKeyPtr := (*C.uchar)(unsafe.Pointer(publicKey))
	return C.crypto_vrf_is_valid_key(publicKeyPtr) != 0
}

// Construct a VRF proof from given secret key and message.
func Prove(privateKey *[SECRETKEYBYTES]byte, message []byte) (*[PROOFBYTES]byte, error) {
    proof := [PROOFBYTES]byte{}
    proofPtr := (*C.uchar)(unsafe.Pointer(&proof))
    privateKeyPtr := (*C.uchar)(unsafe.Pointer(privateKey))
	messagePtr := bytesToUnsignedCharPointer(message)
    messageLen := (C.ulonglong)(len(message))
    if C.crypto_vrf_prove(proofPtr, privateKeyPtr, messagePtr, messageLen) != 0 {
        return nil, errors.New(fmt.Sprintf("unable to decode the given privateKey"))
    }
    return &proof, nil
}

// Verifies that proof was legitimately generated by private key for the given public key, and stores the
// VRF hash in output. Note that VRF "verify()" means the process of generating output from public key,
// proof, and message.
// https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#section-5.3
func Verify(publicKey *[PUBLICKEYBYTES]byte, proof *[PROOFBYTES]byte, message []byte) (*[OUTPUTBYTES]byte, error) {
    output := [OUTPUTBYTES]byte{}
    outputPtr := (*C.uchar)(unsafe.Pointer(&output))
    publicKeyPtr := (*C.uchar)(unsafe.Pointer(publicKey))
    proofPtr := (*C.uchar)(unsafe.Pointer(proof))
    messagePtr := bytesToUnsignedCharPointer(message)
    messageLen := (C.ulonglong)(len(message))
    if C.crypto_vrf_verify(outputPtr, publicKeyPtr, proofPtr, messagePtr, messageLen) != 0 {
        return nil, errors.New(fmt.Sprintf(
            "given public key is invalid, or the proof isn't legitimately generated for the message:"+
                " public_key=%s, proof=%s, message=%s",
            hex.EncodeToString(publicKey[:]), hex.EncodeToString(proof[:]), hex.EncodeToString(message[:])))
    }
    return &output, nil
}

// Calculate the output (hash value) from the specified proof.
// In essence, this function returns a valid value if given proof is any point on the finite field. Otherwise,
// this will return an error.
func ProofToHash(proof *[PROOFBYTES]byte) (*[OUTPUTBYTES]byte, error) {
	output := [OUTPUTBYTES]byte{}
	outputPtr := (*C.uchar)(unsafe.Pointer(&output))
	proofPtr := (*C.uchar)(unsafe.Pointer(proof))
	if C.crypto_vrf_proof_to_hash(outputPtr, proofPtr) != 0 {
		return nil, errors.New(fmt.Sprintf(
			"given proof isn't legitimately generated: proof=%s", hex.EncodeToString(proof[:])))
	}
	return &output, nil
}

func SkToPk(privateKey *[SECRETKEYBYTES]byte) *[PUBLICKEYBYTES]byte {
	publicKey := [PUBLICKEYBYTES]byte{}
	publicKeyPtr := (*C.uchar)(unsafe.Pointer(&publicKey))
	privateKeyPtr := (*C.uchar)(unsafe.Pointer(privateKey))
	C.crypto_vrf_sk_to_pk(publicKeyPtr, privateKeyPtr) // void
	return &publicKey
}

func SkToSeed(privateKey *[SECRETKEYBYTES]byte) *[SEEDBYTES]byte {
	seed := [SEEDBYTES]byte{}
	seedPtr := (*C.uchar)(unsafe.Pointer(&seed))
	privateKeyPtr := (*C.uchar)(unsafe.Pointer(privateKey))
	C.crypto_vrf_sk_to_seed(seedPtr, privateKeyPtr) // void
	return &seed
}

func bytesToUnsignedCharPointer(msg []byte) *C.uchar {
	if len(msg) == 0 {
		return (*C.uchar)(C.NULL)
	}
	return (*C.uchar)(unsafe.Pointer(&msg[0]))
}
