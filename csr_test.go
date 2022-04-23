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
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"os"
	"testing"
)

func TestGenRsaCSR(t *testing.T) {
	subj_ := pkix.Name{
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"California"},
		Organization:       []string{"ZeroSSL"},
		OrganizationalUnit: []string{"ZeroSSL"},
		CommonName:         "example.com",
	}
	rsaPrivKey_ := GenRsaKey(2048)
	csr_, err := GenRsaCSR(subj_, rsaPrivKey_, x509.SHA256WithRSA)
	if err != nil {
		t.Error(err)
	}
	err = WriteRsaPrivKeyPem(os.Stdout, rsaPrivKey_)
	if err != nil {
		t.Error(err)
	}
	err = WriteCSRPem(os.Stdout, csr_)
	if err != nil {
		t.Error(err)
	}
}

func TestGetCSRString(t *testing.T) {
	subj_ := pkix.Name{
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"California"},
		Organization:       []string{"ZeroSSL"},
		OrganizationalUnit: []string{"ZeroSSL"},
		CommonName:         "example.com",
	}
	rsaPrivKey_ := GenRsaKey(2048)
	csr_, err := GenRsaCSR(subj_, rsaPrivKey_, x509.SHA256WithRSA)
	if err != nil {
		t.Error(err)
	}
	t.Log(GetCSRString(csr_))
}

func TestGenECCCSR(t *testing.T) {
	subj_ := pkix.Name{
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"California"},
		Organization:       []string{"ZeroSSL"},
		OrganizationalUnit: []string{"ZeroSSL"},
		CommonName:         "example.com",
	}
	privKey_ := GenEccKey(elliptic.P256())
	csr_, err := GenEccCSR(subj_, privKey_, x509.ECDSAWithSHA256)
	if err != nil {
		t.Error(err)
	}
	_ = WriteEccPrivKeyPem(os.Stdout, privKey_)
	t.Log(GetCSRString(csr_))
}
