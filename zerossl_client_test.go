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
	"fmt"
	"strings"
	"testing"
)

func TestZeroSSLClient_GetCert(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	cert_, err := c_.GetCert("x")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("cert: %#v", cert_)
}

func TestClient_CreateCert(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	privKey_ := GenEccKey(elliptic.P256())
	subj_ := pkix.Name{
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"California"},
		Organization:       []string{"ZeroSSL"},
		OrganizationalUnit: []string{"ZeroSSL"},
		CommonName:         "example.com",
	}
	csr_, err := GenEccCSR(subj_, privKey_, x509.ECDSAWithSHA256)
	if err != nil {
		t.Error(err)
		return
	}
	csrStr_ := GetCSRString(csr_)
	if csrStr_ == "" {
		t.Error("failed to get csr string")
		return
	}
	cert_, err := c_.CreateCert("example.com", csrStr_, "90", "1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("cert: %#v", cert_)
}

func TestClient_DeleteCert(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	err := c_.DeleteCert("x")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("cert deleted")
}

func TestClient_VerifyDomains(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	rspModel_, err := c_.VerifyDomains("x", VerifyDomainsMethod.HttpCsrHash, "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("domains verification result: %#v", rspModel_)
}

func TestClient_VerificationStatus(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	rspModel_, err := c_.VerificationStatus("x")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("verification status: %#v", rspModel_)
}

func TestClient_DownloadCertInline(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	cert_, err := c_.DownloadCertInline("x", "1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("cert: %#v", cert_)
	fullChainPem_ := fmt.Sprintf("%s\n%s\n", strings.TrimSpace(cert_.Certificate), strings.TrimSpace(cert_.CaBundle))
	t.Logf("cert full chain: %s", fullChainPem_)
}

func TestClient_ListCerts(t *testing.T) {
	c_ := &Client{ApiKey: "x"}
	rspModel_, err := c_.ListCerts("", "example.com", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("certs: %#v", rspModel_)
}
