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
	"testing"
)

func TestReadConfig(t *testing.T) {
	conf_, err := ReadConfig("config.yaml")
	if err != nil {
		t.Errorf("ReadConfig failed: %s", err)
	}
	t.Logf("config: %+v", conf_)
}

func TestReadCurrentData(t *testing.T) {
	data_, err := ReadCurrentData("current.yaml")
	if err != nil {
		t.Errorf("ReadCurrentData failed: %s", err)
	}
	t.Logf("data: %+v", data_)
}

func TestWriteCurrentData(t *testing.T) {
	data_, err := ReadCurrentData("current.yaml")
	if err != nil {
		t.Errorf("ReadCurrentData failed: %s", err)
	}
	for k := range data_.Certs {
		data_.Certs[k].CertID += "1"
	}
	err = WriteCurrentData("current.yaml", data_)
	if err != nil {
		t.Errorf("WriteCurrentData failed: %s", err)
	}
}
