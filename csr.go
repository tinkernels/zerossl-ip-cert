/*
 * Copyright [2022] [tinkernels (github.com/tinkernels)]
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zerosslIPCert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var SignatureAlgorithms = map[string]x509.SignatureAlgorithm{
	"SHA256-RSA":   x509.SHA256WithRSA,
	"SHA384-RSA":   x509.SHA384WithRSA,
	"ECDSA-SHA256": x509.ECDSAWithSHA256,
	"ECDSA-SHA384": x509.ECDSAWithSHA384,
}

var EcdsaCurves = map[string]elliptic.Curve{
	"P-256": elliptic.P256(),
	"P-384": elliptic.P384(),
}

// KeyGeneratorWrapper is a wrapper for generating keys.
func KeyGeneratorWrapper(keyType string, keyBits int, keyCurve string) (key interface{}) {
	keyType_ := strings.ToUpper(keyType)
	keyCurve_ := strings.ToUpper(keyCurve)
	switch keyType_ {
	case "RSA":
		{
			key = GenRsaKey(keyBits)
		}
	case "ECDSA":
		{
			keyCurve, ok := EcdsaCurves[keyCurve_]
			if !ok {
				panic("invalid ecdsa curve")
			}
			key = GenEccKey(keyCurve)
		}
	}
	return
}

// WritePrivKeyWrapper is a wrapper for writing private keys.
func WritePrivKeyWrapper(keyType string, key interface{}, keyFile string) (err error) {
	keyType_ := strings.ToUpper(keyType)
	file_, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(file_ *os.File) {
		err := file_.Close()
		if err != nil {
			log.Println(err)
		}
	}(file_)
	switch keyType_ {
	case "RSA":
		{
			err = WriteRsaPrivKeyPem(file_, key.(*rsa.PrivateKey))
			if err != nil {
				log.Println(err)
				return
			}
		}
	case "ECDSA":
		{
			err = WriteEccPrivKeyPem(file_, key.(*ecdsa.PrivateKey))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	return
}

// CSRGeneratorWrapper is a wrapper for generating CSR.
func CSRGeneratorWrapper(keyType string, subj pkix.Name, key interface{}, sigAlgStr string) (csr []byte, err error) {
	keyType_ := strings.ToUpper(keyType)
	sigAlgStr_ := strings.ToUpper(sigAlgStr)
	sigAlg_, ok := SignatureAlgorithms[sigAlgStr_]
	if !ok {
		err = fmt.Errorf("invalid signature algorithm")
		return
	}
	switch keyType_ {
	case "RSA":
		{
			csr, err = GenRsaCSR(subj, key.(*rsa.PrivateKey), sigAlg_)
		}
	case "ECDSA":
		{
			csr, err = GenEccCSR(subj, key.(*ecdsa.PrivateKey), sigAlg_)
		}
	}
	return
}

// GenRsaKey generates a new RSA private key.
func GenRsaKey(bits int) *rsa.PrivateKey {
	privKey_, _ := rsa.GenerateKey(rand.Reader, bits)
	return privKey_
}

// GenRsaCSR generates a new RSA CSR.
func GenRsaCSR(subj pkix.Name, key *rsa.PrivateKey, sigAlg x509.SignatureAlgorithm) (csr []byte, err error) {
	template_ := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: sigAlg,
	}
	csr, err = x509.CreateCertificateRequest(rand.Reader, &template_, key)
	return
}

// GenEccKey generates a new ECC private key.
func GenEccKey(curve elliptic.Curve) (key *ecdsa.PrivateKey) {
	key, _ = ecdsa.GenerateKey(curve, rand.Reader)
	return
}

// GenEccCSR generates a new ECC CSR.
func GenEccCSR(subj pkix.Name, key *ecdsa.PrivateKey, sigAlg x509.SignatureAlgorithm) (csr []byte, err error) {
	template_ := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: sigAlg,
	}
	csr, err = x509.CreateCertificateRequest(rand.Reader, &template_, key)
	return
}

// WriteRsaPrivKeyPem writes an RSA private key to a PEM file.
func WriteRsaPrivKeyPem(out io.Writer, key *rsa.PrivateKey) (err error) {
	err = pem.Encode(out, &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return
}

// WriteEccPrivKeyPem writes an ECC private key to a PEM file.
func WriteEccPrivKeyPem(out io.Writer, key *ecdsa.PrivateKey) (err error) {
	privKBytes_, _ := x509.MarshalECPrivateKey(key)
	err = pem.Encode(out, &pem.Block{Type: "PRIVATE KEY", Bytes: privKBytes_})
	return
}

// WriteCSRPem writes a CSR to a PEM file.
func WriteCSRPem(out io.Writer, csr []byte) (err error) {
	err = pem.Encode(out, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})
	return
}

// GetCSRString returns a CSR as a string.
func GetCSRString(csr []byte) (csrStr string) {
	buf_ := bytes.NewBuffer([]byte{})
	err := WriteCSRPem(buf_, csr)
	if err != nil {
		return ""
	}
	return buf_.String()
}
