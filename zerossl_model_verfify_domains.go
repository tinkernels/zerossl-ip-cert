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

// VerifyDomainsMethod represents the method of verifying domains.
var VerifyDomainsMethod = struct {
	Email        string // EMAIL
	CNameCsrHash string //CNAME_CSR_HASH
	HttpCsrHash  string //HTTP_CSR_HASH
	HttpsCsrHash string //HTTPS_CSR_HASH
}{
	Email:        "EMAIL",
	CNameCsrHash: "CNAME_CSR_HASH",
	HttpCsrHash:  "HTTP_CSR_HASH",
	HttpsCsrHash: "HTTPS_CSR_HASH",
}

type VerifyDomainsModel struct {
	Success bool                    `json:"success"`
	Error   VerifyDomainsErrorModel `json:"error"`
}

type VerifyDomainsErrorModel struct {
	Code    int                           `json:"code"`
	Type    string                        `json:"type"`
	Details VerifyDomainsErrorDetailModel `json:"details"`
}

type VerifyDomainsErrorDetailModel struct {
	CNameFound    int    `json:"cname_found"`
	RecordCorrect int    `json:"record_correct"`
	TargetHost    string `json:"target_host"`
	TargetRecord  string `json:"target_record"`
	ActualRecord  string `json:"actual_record"`
}
