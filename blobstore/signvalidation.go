package blobstore

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"io"
)

type privateKey *rsa.PrivateKey

func createSignValidatedBlobFromReaderGenerator(
	readerGenerator func() io.Reader,
	privKey privateKey,
	dataVersion int64,
	storage BlobStorage,
) (
	// Return values
	bidRet string,
	keyRet string,
	err error,
) {

	// We're using hash of the private key to create the encryption data key
	dataKey := sha512.New().Sum(x509.MarshalPKCS1PrivateKey(privKey))

	// Version + encrypted data buffer
	verDataBuffer := bytes.Buffer{}
	serializeInt(dataVersion, &verDataBuffer)

	// Encrypt the data, note we're using the version bytes as IV
	encryptedWriter, key, err := createEncryptor(dataKey, verDataBuffer.Bytes(), &verDataBuffer)
	if err != nil {
		return
	}
	io.Copy(encryptedWriter, readerGenerator())

	// Calculate the signature of version + encrypted data blob
	signature, err := rsa.SignPKCS1v15(
		nil, privKey, crypto.SHA512,
		crypto.SHA512.New().Sum(verDataBuffer.Bytes()))
	if err != nil {
		return
	}

	// Generate the public key blob
	pubKey, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return
	}

	// Generate the BID from the public key
	bid := hex.EncodeToString(sha512.New().Sum(pubKey))

	// Open the blob for writing
	blobWriter, err := storage.NewBlobWriter(bid)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			blobWriter.Cancel()
		}
	}()

	// Write blob header (without version)
	header := bytes.Buffer{}
	header.WriteByte(validationMethodSign)
	serializeBuffer(pubKey, &header)
	serializeBuffer(signature, &header)

	if _, err = blobWriter.Write(header.Bytes()); err != nil {
		return
	}

	// Write the version and encrypted data
	if _, err = blobWriter.Write(verDataBuffer.Bytes()); err != nil {
		return
	}

	// Ok, we're done here
	return bid, key, nil
}
