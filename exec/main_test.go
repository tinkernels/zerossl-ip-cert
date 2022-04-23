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

package main

import (
	"github.com/tinkernels/zerossl-ip-cert"
	"testing"
)

func Test_verifyHook(t *testing.T) {
	certInfoTest_ := zerosslIPCert.CertificateInfoModel{
		CommonName: "1.1.1.1",
		Validation: zerosslIPCert.ValidationInfoModel{
			OtherMethods: map[string]zerosslIPCert.OtherValidationInfoModel{
				"1.1.1.1": {
					FileValidationUrlHttp: "http://1.1.1.1/.well-known/pki-validation/715EE529C6FF317C938B79C7655710AC.txt",
					FileValidationContent: []string{
						"ABCDEF1234567890",
						" comodoca.com",
						"abcdef1234567890",
					},
				},
			},
		},
	}
	err := runVerifyHook("/Users/donjohnny/forge/sources/zerossl-ip-cert/exec/sample-nginx-verify-hook.sh", &certInfoTest_)
	if err != nil {
		t.Error(err)
		return
	}
}
