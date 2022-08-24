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

// CertStatus represents the status of a certificate.
var CertStatus = struct {
	Draft             string
	PendingValidation string
	Issued            string
	Cancelled         string
	ExpiringSoon      string
	Expired           string
}{
	Draft:             "draft",
	PendingValidation: "pending_validation",
	Issued:            "issued",
	Cancelled:         "cancelled",
	ExpiringSoon:      "expiring_soon",
	Expired:           "expired",
}

type CertificateInfoModel struct {
	ID                string              `json:"id"`
	Type              string              `json:"type"`
	CommonName        string              `json:"common_name"`
	AdditionalDomains string              `json:"additional_domains"`
	Created           string              `json:"created"`
	Expires           string              `json:"expires"`
	Status            string              `json:"status"`
	ValidationType    string              `json:"validation_type"`
	ValidationEmails  string              `json:"validation_email"`
	ReplacementFor    string              `json:"replacement_for"`
	Validation        ValidationInfoModel `json:"validation,omitempty"`
}

type ValidationInfoModel struct {
	EmailValidation map[string][]string                 `json:"email_validation"`
	OtherMethods    map[string]OtherValidationInfoModel `json:"other_methods"`
}

type OtherValidationInfoModel struct {
	FileValidationUrlHttp  string   `json:"file_validation_url_http"`
	FileValidationUrlHttps string   `json:"file_validation_url_https"`
	FileValidationContent  []string `json:"file_validation_content"`
	CNameValidationP1      string   `json:"cname_validation_p1"`
	CNameValidationP2      string   `json:"cname_validation_p2"`
}
