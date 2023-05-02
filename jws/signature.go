package jws

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/alexandru-ionut-balan/ing-jws/crypto"
	"github.com/alexandru-ionut-balan/ing-jws/logging"
)

func generateEncodedHeader(jwsHeader *JwsHeader) (string, error) {
	rawHeaderBytes, err := json.Marshal(jwsHeader)
	logging.Info("Encoding jws header: " + string(rawHeaderBytes))
	if err != nil {
		logging.Error("Cannot marshal JWS Header into JSON.", nil)
		return "", err
	}

	return crypto.ApplyExtraFormatting(crypto.Base64(rawHeaderBytes)), nil
}

func generateSignatureValue(encodedJwsHeader string, httpHeaders map[string]string, privateKey *rsa.PrivateKey) (string, error) {
	signatureInput := encodedJwsHeader + "."

	for name, value := range httpHeaders {
		signatureInput += (name + ": " + value + "\n")
	}

	logging.Info("Signing jws value: " + signatureInput[:len(signatureInput)-1])

	signedInput, err := crypto.Sign(signatureInput[:len(signatureInput)-1], privateKey)
	if err != nil {
		logging.Error("Cannot generate signature value.", nil)
		return "", err
	}

	return crypto.ApplyExtraFormatting(crypto.Base64(signedInput)), nil
}

func GenerateSignature(jwsHeader *JwsHeader, httpHeaders map[string]string, privateKey *rsa.PrivateKey) (string, error) {
	encodedHeader, err := generateEncodedHeader(jwsHeader)
	if err != nil {
		logging.Error("Cannot create signature!", nil)
		return "", err
	}

	signatureValue, err := generateSignatureValue(encodedHeader, httpHeaders, privateKey)
	if err != nil {
		logging.Error("Cannot create signature!", nil)
		return "", err
	}

	return encodedHeader + ".." + signatureValue, nil
}
