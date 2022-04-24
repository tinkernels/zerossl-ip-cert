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
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type CertConf struct {
	ConfID           string `yaml:"confId"`
	ApiKey           string `yaml:"apiKey"`
	Country          string `yaml:"country"`
	Province         string `yaml:"province"`
	City             string `yaml:"city"`
	Locality         string `yaml:"locality"`
	Organization     string `yaml:"organization"`
	OrganizationUnit string `yaml:"organizationUnit"`
	CommonName       string `yaml:"commonName"`
	Days             int    `yaml:"days"`
	KeyType          string `yaml:"keyType"`
	KeyBits          int    `yaml:"keyBits"`
	KeyCurve         string `yaml:"keyCurve"`
	SigAlg           string `yaml:"sigAlg"`
	StrictDomains    int    `yaml:"strictDomains"`
	VerifyMethod     string `yaml:"verifyMethod"`
	VerifyHook       string `yaml:"verifyHook"`
	PostHook         string `yaml:"postHook"`
	CertFile         string `yaml:"certFile"`
	KeyFile          string `yaml:"keyFile"`
}

type Config struct {
	DataDir         string     `yaml:"dataDir"`
	LogFile         string     `yaml:"logFile"`
	CleanUnfinished bool       `yaml:"cleanUnfinished"`
	CertConfigs     []CertConf `yaml:"certConfigs"`
}

// ReadConfig reads the config file and returns a Config struct.
func ReadConfig(path string) (config *Config, err error) {
	var input_ []byte
	input_, err = ioutil.ReadFile(path)
	err = yaml.Unmarshal(input_, &config)
	if err != nil {
		return nil, err
	}
	return
}

type CurrentData struct {
	Certs []CurrentCertData `yaml:"certs"`
}

type CurrentCertData struct {
	CommonName string `yaml:"commonName"`
	ConfID     string `yaml:"confId"`
	CertID     string `yaml:"certId"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
}

// ReadCurrentData reads the current data file and returns a CurrentData struct.
func ReadCurrentData(path string) (data *CurrentData, err error) {
	var input_ []byte
	input_, err = ioutil.ReadFile(path)
	if err = yaml.Unmarshal(input_, &data); err != nil {
		return nil, err
	}
	return
}

// WriteCurrentData writes the current data file.
func WriteCurrentData(path string, data *CurrentData) (err error) {
	var output []byte
	output, err = yaml.Marshal(data)
	if err = ioutil.WriteFile(path, output, os.ModePerm); err != nil {
		return err
	}
	return
}
