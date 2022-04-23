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
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"github.com/tinkernels/zerossl-ip-cert"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Version is the version of this application.
const Version = "0.2.1"

var (
	renewFlag  = flag.Bool("renew", false, "Renew existing certs only")
	configFlag = flag.String("config", "", "Config file")
)

var usingConfig *Config
var currentData *CurrentData
var currentDataFilePath string

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "\nVersion: %v\n\nUsage: %v [ -renew ] -config CONFIG_FILE\n\n",
			Version, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if !PathExists(*configFlag) {
		flag.Usage()
		panic("Config file not found")
	}
	usingConfig_, err := ReadConfig(*configFlag)
	if err != nil {
		flag.Usage()
		panic(err)
	}
	usingConfig = usingConfig_
	// Enable line numbers in logging.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile_, err := os.OpenFile(usingConfig.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("log file create failed")
		panic(err)
	}
	// Write log to both console and file.
	multiLogWr_ := io.MultiWriter(os.Stdout, logFile_)
	log.SetOutput(multiLogWr_)

	log.Printf("Using config file: %v", *configFlag)

	err = CreateDirIfNotExists(usingConfig.DataDir, os.ModePerm)
	if err != nil {
		flag.Usage()
		panic(err)
	}
	currentDataFilePath = filepath.Join(usingConfig.DataDir, "/current.yaml")
	if PathExists(currentDataFilePath) {
		currentData_, err := ReadCurrentData(currentDataFilePath)
		if err != nil {
			flag.Usage()
			panic(err)
		}
		if currentData_ == nil {
			currentData = &CurrentData{}
		} else {
			currentData = currentData_
		}
	} else {
		currentData = &CurrentData{}
	}
	if *renewFlag {
		renew()
	} else {
		issueCerts()
	}
}

// issueCerts issues certs referenced in the config file.
func issueCerts() {
	log.Printf("Issuing certs")
	for _, c := range usingConfig.CertConfigs {
		log.Printf("Issuing cert for domain: %v", c.CommonName)
		err := issueCert(&c)
		if err != nil {
			log.Printf("Failed to issue cert for domain %v: %v\n", c.CommonName, err)
		}
	}
}

// issueCert issues a cert for the given domain config.
func issueCert(conf *CertConf) (err error) {
	for _, cert := range currentData.Certs {
		// Use ConfID to match.
		if cert.ConfID == conf.ConfID {
			log.Printf("Cert for domain %v already exists, try renew.\n", conf.CommonName)
			err = renewCert(cert.CertID, conf)
			return
		}
	}
	log.Printf("Cert for domain %v does not exist, try issue.\n", conf.CommonName)
	certId_, err := issueCertImpl(conf)
	if err == nil {
		log.Printf("Cert for domain %v issued successfully.\n", conf.CommonName)
		currentData.Certs = append(currentData.Certs, CurrentCertData{
			CommonName: conf.CommonName,
			CertID:     certId_,
			CertFile:   conf.CertFile,
			KeyFile:    conf.KeyFile,
			ConfID:     conf.ConfID,
		})
		if err = WriteCurrentData(currentDataFilePath, currentData); err != nil {
			log.Printf("Failed to write current data: %v\n", err)
		}
	}
	return
}

func issueCertImpl(conf *CertConf) (certID string, err error) {
	tempDir_ := filepath.Join(usingConfig.DataDir, "/temp")
	tempPrivKeyPath_ := filepath.Join(tempDir_, "/privkey.pem")
	log.Printf("tempPrivKeyPath: %v\n", tempPrivKeyPath_)
	tempCertPath_ := filepath.Join(tempDir_, "/cert-fullchain.pem")
	log.Printf("tempCertPath: %v\n", tempCertPath_)
	log.Printf("Cleaning temp dir: %v\n", tempDir_)
	if err = os.RemoveAll(tempDir_); err != nil {
		return
	}
	log.Printf("Creating temp dir: %v\n", tempDir_)
	if err = CreateDirIfNotExists(tempDir_, os.ModePerm); err != nil {
		return
	}
	client_ := &zerosslIPCert.Client{ApiKey: conf.ApiKey}
	// Generate PrivateKey.
	log.Printf("Generating private key for %v\n", conf.CommonName)
	privKey_ := zerosslIPCert.KeyGeneratorWrapper(conf.KeyType, conf.KeyBits, conf.KeyCurve)
	subj_ := pkix.Name{
		Country:            []string{conf.Country},
		Province:           []string{conf.Province},
		Locality:           []string{conf.Locality},
		Organization:       []string{conf.Organization},
		OrganizationalUnit: []string{conf.OrganizationUnit},
		CommonName:         conf.CommonName,
	}
	// Generate CSR.
	log.Printf("Generating CSR for %v\n", conf.CommonName)
	csr_, err := zerosslIPCert.CSRGeneratorWrapper(conf.KeyType, subj_, privKey_, conf.SigAlg)
	if err != nil {
		log.Println(err)
		return
	}
	csrStr_ := zerosslIPCert.GetCSRString(csr_)
	log.Printf("CSR for %v: %v\n", conf.CommonName, csrStr_)
	if csrStr_ == "" {
		log.Println("failed to get csr string")
		return
	}
	// Write PrivateKey to file.
	log.Printf("Writing private key to file %v\n", tempPrivKeyPath_)
	if err = zerosslIPCert.WritePrivKeyWrapper(conf.KeyType, privKey_, tempPrivKeyPath_); err != nil {
		log.Println(err)
		return
	}
	// Create Cert.
	log.Printf("Creating cert for %v\n", conf.CommonName)
	certInfo_, err := client_.CreateCert(conf.CommonName, csrStr_, strconv.Itoa(conf.Days),
		strconv.Itoa(conf.StrictDomains))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("cert info: %+v\n", certInfo_)
	if err = runVerifyHook(conf.VerifyHook, &certInfo_); err != nil {
		log.Println(err)
		return
	}
	// Verify Domains.
	verifyRsp_, err := client_.VerifyDomains(certInfo_.ID, zerosslIPCert.VerifyDomainsMethod.HttpCsrHash, "")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("domains verification result: %+v\n", verifyRsp_)
	// Wait seconds before checking cert status.
	time.Sleep(time.Second * 5)
	certInfoTmp_, err := client_.GetCert(certInfo_.ID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("cert info: %+v", certInfoTmp_)
	if certInfoTmp_.Status == zerosslIPCert.CertStatus.Draft {
		err = fmt.Errorf("cert verification failed, still in draft status")
		log.Println(err)
		return "", err
	}
	// Wait for cert to be ready.
	if err = waitCert2BReady(client_, &certInfo_); err != nil {
		log.Println(err)
		return
	}
	// Download cert.
	cert_, err := client_.DownloadCertInline(certInfo_.ID, "1")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("cert + ca: %+v\n", cert_)
	fullChainPem_ := fmt.Sprintf("%s\n%s\n", strings.TrimSpace(cert_.Certificate), strings.TrimSpace(cert_.CaBundle))
	// Write cert to file.
	file_, err := os.Create(tempCertPath_)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = file_.WriteString(fullChainPem_)
	if err != nil {
		return
	}
	// Copy cert files to dest.
	if err = CopyFile(tempCertPath_, conf.CertFile, os.ModePerm); err != nil {
		log.Println(err)
		return
	}
	if err = CopyFile(tempPrivKeyPath_, conf.KeyFile, os.ModePerm); err != nil {
		log.Println(err)
		return
	}
	// Run post hook.
	if err = runPostHook(conf); err != nil {
		log.Println(err)
		return
	}
	// Clean temp files.
	log.Printf("Cleaning temp files\n")
	_ = os.RemoveAll(tempDir_)
	certID = certInfo_.ID
	return
}

// runVerifyHook runs verify hook.
func runVerifyHook(executable string, cerInfo *zerosslIPCert.CertificateInfoModel) (err error) {
	if !PathExists(executable) {
		return fmt.Errorf("verify hook executable %v not exists", executable)
	}
	log.Println("try make verify hook file executable")
	err = ChmodPlusX(executable)
	if err != nil {
		log.Printf("chmod verify hook file permission failed: %v\n", err)
	}
	for k, v := range cerInfo.Validation.OtherMethods {
		if k == cerInfo.CommonName {
			validateHttpUrl_, err := url.Parse(v.FileValidationUrlHttp)
			if err != nil {
				log.Println(err)
				return err
			}
			host_ := validateHttpUrl_.Host
			path_ := validateHttpUrl_.Path
			port_ := validateHttpUrl_.Port()
			if port_ == "" {
				port_ = "80"
			}
			content_ := strings.Join(v.FileValidationContent, "\n")
			// Prepare hook exec env.
			cmdEnv_ := os.Environ()
			cmdEnv_ = append(cmdEnv_, fmt.Sprintf("%v=%v", "ZEROSSL_HTTP_FV_HOST", host_))
			cmdEnv_ = append(cmdEnv_, fmt.Sprintf("%v=%v", "ZEROSSL_HTTP_FV_PATH", path_))
			cmdEnv_ = append(cmdEnv_, fmt.Sprintf("%v=%v", "ZEROSSL_HTTP_FV_PORT", port_))
			cmdEnv_ = append(cmdEnv_, fmt.Sprintf("%v=%v", "ZEROSSL_HTTP_FV_CONTENT", content_))
			cmd_ := exec.Command(executable)
			cmd_.Env = cmdEnv_
			cmd_.Stdout = os.Stdout
			cmd_.Stderr = os.Stdout
			if err = cmd_.Run(); err != nil {
				return err
			}
			return err
		}
	}
	return
}

// waitCert2BReady waits for the cert to be ready.
func waitCert2BReady(client *zerosslIPCert.Client, certInfo *zerosslIPCert.CertificateInfoModel) (err error) {
	for {
		// loop every other seconds until cert is ready.
		certInfo_, err := client.GetCert(certInfo.ID)
		if err != nil {
			log.Println(err)
			return err
		}
		if certInfo_.Status == zerosslIPCert.CertStatus.Issued {
			log.Printf("cert is ready: %+v\n", certInfo_)
			return nil
		}
		time.Sleep(time.Second * 15)
	}
}

func runPostHook(certConf *CertConf) (err error) {
	if !PathExists(certConf.PostHook) {
		return fmt.Errorf("post hook executable %v not exists", certConf.PostHook)
	}
	log.Println("try make post hook file executable")
	err = ChmodPlusX(certConf.PostHook)
	if err != nil {
		log.Printf("chmod +x post hook file failed: %v\n", err)
	}
	// Prepare hook exec env.
	cmdEnv_ := os.Environ()
	cmdEnv_ = append(cmdEnv_, fmt.Sprintf("%v=%v", "ZEROSSL_CERT_FPATH", certConf.CertFile))
	cmdEnv_ = append(cmdEnv_, fmt.Sprintf("%v=%v", "ZEROSSL_KEY_FPATH", certConf.KeyFile))
	cmd_ := exec.Command(certConf.PostHook)
	cmd_.Env = cmdEnv_
	cmd_.Stdout = os.Stdout
	cmd_.Stderr = os.Stdout
	if err = cmd_.Run(); err != nil {
		return err
	}
	return
}

// renew current certs.
func renew() {
	log.Println("will renew current certs")
loopRenew:
	for _, cert := range currentData.Certs {
		log.Printf("try renew cert: %v\n", cert.CommonName)
		for _, c := range usingConfig.CertConfigs {
			// ConfID to match cert config.
			if c.ConfID == cert.ConfID {
				err := renewCert(cert.CertID, &c)
				if err != nil {
					log.Printf("Failed to renew cert for domain %v: %v\n", c.CommonName, err)
				}
				continue loopRenew
			}
		}
		log.Printf("no config for renewing cert: %v\n", cert.CommonName)
	}
}

func renewCert(id string, conf *CertConf) (err error) {
	log.Printf("Renewing cert %v with config: %v\n", conf.CommonName, conf.ConfID)
	client_ := &zerosslIPCert.Client{ApiKey: conf.ApiKey}
	certInfo_, err := client_.GetCert(id)
	if err != nil {
		log.Printf("Failed to get cert info: %v\n", err)
		return err
	}
	if certInfo_.Status != zerosslIPCert.CertStatus.ExpiringSoon {
		log.Printf("Cert %v is not due for renewal, skip renewing.\n", conf.CommonName)
		return nil
	}
	certId_, err := issueCertImpl(conf)
	if err == nil {
		log.Printf("Cert for domain %v issued successfully.\n", conf.CommonName)
		for i, c := range currentData.Certs {
			// Use original cert ID to match cert.
			if c.CertID == id {
				currentData.Certs[i].ConfID = conf.ConfID
				currentData.Certs[i].CommonName = conf.CommonName
				currentData.Certs[i].CertID = certId_
				currentData.Certs[i].CertFile = conf.CertFile
				currentData.Certs[i].KeyFile = conf.KeyFile
				break
			}
		}
		if err = WriteCurrentData(currentDataFilePath, currentData); err != nil {
			log.Printf("Failed to write current data: %v\n", err)
		}
	}
	return
}
