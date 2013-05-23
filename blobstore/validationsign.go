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
type publicKey *rsa.PublicKey

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

func createReaderForSignedBlobData(reader io.Reader, bid, key string) (rawReader io.Reader, err error) {

	// Grab the public key blob
	pubkey, err := deserializeBuffer(reader, maxSanePubKeyLength)
	if err != nil {
		return
	}

	// Validate blob id agains public key
	if hex.EncodeToString(sha512.New().Sum(pubkey)) != bid {
		return nil, ErrInvalidPublicKeyBid
	}

	// Parse the public key
	pubKeyParsedRaw, err := x509.ParsePKIXPublicKey(pubkey)
	if err != nil {
		return
	}
	pubKeyParsed, ok := pubKeyParsedRaw.(publicKey)
	if !ok {
		return nil, ErrUnknownPublicKeyType
	}

	// Read the signature
	signature, err := deserializeBuffer(reader, maxSaneSignatureLength)
	if err != nil {
		return
	}

	// Read the version
	version, err := deserializeInt(reader)
	if err != nil {
		return
	}

	// TODO: Create validating reader that will check the signature when
	//       we reach EOF
	_, _ = signature, pubKeyParsed

	// Create the decryptor of the content
	verBuffer := bytes.Buffer{}
	serializeInt(version, &verBuffer)
	return createDecryptor(key, verBuffer.Bytes(), reader)
}

func createReaderForSignedBlob(bid string, key string, storage BlobStorage) (rawReader io.Reader, err error) {

	// Get the reader
	encryptedReader, err := storage.NewBlobReader(bid)
	if err != nil {
		return
	}

	// Test the validation method
	validationType, err := deserializeInt(encryptedReader)
	if err != nil {
		return
	}
	if validationType != validationMethodSign {
		return nil, ErrInvalidValidationMethod
	}

	// Get the encryptor
	return createReaderForSignedBlobData(encryptedReader, bid, key)
}
