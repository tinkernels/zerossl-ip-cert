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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CreateDirIfNotExists creates a directory if it does not exist.
func CreateDirIfNotExists(dir string, perm os.FileMode) error {
	if PathExists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}
	return nil
}

// PathExists checks if a path exists.
func PathExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CopyFile copies content of a file from src to dst.
func CopyFile(srcFile, dstFile string, perm os.FileMode) error {
	dstDir_ := filepath.Dir(dstFile)
	err := CreateDirIfNotExists(dstDir_, perm)
	if err != nil {
		return err
	}
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Printf("failed to close file: '%s', error: '%s'\n", dstFile, err.Error())
		}
	}(out)
	in, err := os.Open(srcFile)
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			fmt.Printf("failed to close file: '%s', error: '%s'\n", srcFile, err.Error())
		}
	}(in)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

func ChmodPlusX(file string) (err error) {
	if runtime.GOOS != "windows" {
		err = exec.Command("/usr/bin/env", "chmod", "+x", file).Run()
	}
	return
}
