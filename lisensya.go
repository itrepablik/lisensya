/*
Copyright © 2020 ITRepablik <support@itrepablik.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lisensya

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt is the encryptor of any classified text.
func Encrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext), nil
}

// Decrypt is the decryptor of the some classified text.
func Decrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}

// RevokeLicenseKey revokes existing gokopy license key.
func RevokeLicenseKey(licenseKey, modifiedBy, APIEndPoint, secretKey string) (bool, error) {
	// Get the hostname for the source machine.
	hostName, err := GetHostName()
	if err != nil {
		return false, err
	}

	// Compose the JSON post payload to the API endpoint.
	message := map[string]interface{}{
		"hostName":   hostName,
		"licenseKey": licenseKey,
		"modifiedBy": modifiedBy,
		"secretKey":  secretKey,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(APIEndPoint, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Get the specific return value from your API endpoint.
	i := result["isSuccess"] // You must return with JSON format value of either 'true' or 'false' only.
	mStatus := fmt.Sprint(i)
	isSuccess, _ := strconv.ParseBool(mStatus)
	return isSuccess, nil
}

// GenerateLicenseKey writes the new license key to a custom file and stores in the root directory of your
// app directory, e.g appname.license
func GenerateLicenseKey(licenseKey, appName, secretKey string) (string, error) {
	// Create a license file if not exist with the '.license' custom file format.
	keyFile := strings.ToLower(appName) + ".license"
	f, err := os.OpenFile(keyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Write a new license key to your 'appname.license' custom file.
	newLicenseKey, err := Encrypt(licenseKey, secretKey)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(f.Name(), []byte(strings.TrimSpace(newLicenseKey)), 0644)
	if err != nil {
		return "", err
	}
	return newLicenseKey, nil
}

// ReadLicenseKey reads the license key if found from a custom file, otherwise, throws an error.
func ReadLicenseKey(appName string) (string, error) {
	// Open the custom license file
	file, err := os.Open(strings.ToLower(appName) + ".license")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the license key
	rKey, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(rKey), nil
}

// GetHostName gets the source machine's hostname.
func GetHostName() (string, error) {
	PCName, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return PCName, nil
}
