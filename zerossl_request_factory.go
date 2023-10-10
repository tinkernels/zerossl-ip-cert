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
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ApiEndpoint is the zerossl api endpoint.
const ApiEndpoint = "api.zerossl.com"

// ApiReqFactory is a factory for creating API requests.
var ApiReqFactory = struct {
	// Request of creating a new certificate.
	CreateCertificate func(accessKey, certificateDomains, certificateCsr, certificateValidityDays,
		strictDomains string) (req *http.Request)
	// Request of listing all certificates.
	ListCertificates func(accessKey, certificateStatus, search, limit, page string) (req *http.Request)
	// Request of getting a certificate.
	GetCertificate func(accessKey, id string) (req *http.Request)
	// Request of verifying a certificate.
	VerifyDomains func(accessKey, certificateId, validationMethod, validationEmail string) (req *http.Request)
	// Request of verification status.
	VerificationStatus func(accessKey, id string) (req *http.Request)
	// Request of cancelation a certificate.
	CancelCertificate func(accessKey, id string) (req *http.Request)
	// Request of downloading a certificate.
	DownloadCertificateInline func(accessKey, certID, includeCrossSigned string) (req *http.Request)
}{
	CreateCertificate: func(accessKey, certificateDomains, certificateCsr, certificateValidityDays,
		strictDomains string) (req *http.Request) {
		req = &http.Request{Method: http.MethodPost}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates"}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		url_.RawQuery = q_.Encode()
		req.URL = url_
		bodyForm_ := make(url.Values)
		if certificateDomains != "" {
			bodyForm_.Add("certificate_domains", certificateDomains)
		}
		if certificateCsr != "" {
			bodyForm_.Add("certificate_csr", certificateCsr)
		}
		if certificateValidityDays != "" {
			bodyForm_.Add("certificate_validity_days", certificateValidityDays)
		}
		if strictDomains != "" {
			bodyForm_.Add("strict_domains", strictDomains)
		}
		if len(bodyForm_) > 0 {
			req.Header = make(http.Header)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Body = io.NopCloser(strings.NewReader(bodyForm_.Encode()))
		}
		return
	},
	ListCertificates: func(accessKey, certificateStatus, search, limit, page string) (req *http.Request) {
		req = &http.Request{Method: http.MethodGet}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates"}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		if certificateStatus != "" {
			q_.Add("certificate_status", certificateStatus)
		}
		if search != "" {
			q_.Add("search", search)
		}
		if limit != "" {
			q_.Add("limit", limit)
		}
		if page != "" {
			q_.Add("page", page)
		}
		url_.RawQuery = q_.Encode()
		req.URL = url_
		return
	},
	GetCertificate: func(accessKey, id string) (req *http.Request) {
		req = &http.Request{Method: http.MethodGet}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates/" + id}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		url_.RawQuery = q_.Encode()
		req.URL = url_
		return
	},
	VerifyDomains: func(accessKey, certificateId, validationMethod, validationEmail string) (req *http.Request) {
		req = &http.Request{Method: http.MethodPost}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates/" + certificateId + "/challenges"}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		url_.RawQuery = q_.Encode()
		req.URL = url_
		bodyForm_ := make(url.Values)
		if validationMethod != "" {
			bodyForm_.Add("validation_method", validationMethod)
		}
		if validationEmail != "" {
			bodyForm_.Add("validation_email", validationEmail)
		}
		if len(bodyForm_) > 0 {
			req.Header = make(http.Header)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Body = io.NopCloser(strings.NewReader(bodyForm_.Encode()))
		}
		return
	},
	VerificationStatus: func(accessKey, id string) (req *http.Request) {
		req = &http.Request{Method: http.MethodGet}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates/" + id + "/status"}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		url_.RawQuery = q_.Encode()
		req.URL = url_
		return
	},
	CancelCertificate: func(accessKey, id string) (req *http.Request) {
		req = &http.Request{Method: http.MethodPost}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates/" + id + "/cancel"}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		url_.RawQuery = q_.Encode()
		req.URL = url_
		return
	},
	DownloadCertificateInline: func(accessKey, certID, includeCrossSigned string) (req *http.Request) {
		req = &http.Request{Method: http.MethodGet}
		url_ := &url.URL{Scheme: "https", Host: ApiEndpoint, Path: "/certificates/" + certID + "/download/return"}
		q_ := make(url.Values)
		q_.Add("access_key", accessKey)
		if includeCrossSigned != "" {
			q_.Add("include_cross_signed", includeCrossSigned)
		}
		url_.RawQuery = q_.Encode()
		req.URL = url_
		return
	},
}
