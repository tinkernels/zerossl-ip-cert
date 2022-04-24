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

type ListCertsModel struct {
	TotalCount  int `json:"total_count"`
	ResultCount int `json:"result_count"`
	// Don't use the page field, because page in response is dynamic typed.
	//Page        string                 `json:"page"`
	Limit   int                    `json:"limit"`
	Results []CertificateInfoModel `json:"results"`
}
