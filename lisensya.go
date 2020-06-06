// Package lisensya is the simple license key generator tool for your Go's applications.
package lisensya

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/itrepablik/itrdsn"
	"github.com/itrepablik/tago"
)

const (
	_fileExt                = ".license"
	_defaultExpiryDelimiter = ";expiry:"
)

// GenerateLicenseKey writes the new license key to a custom file and stores in the root directory of your
// app directory, e.g appname.license
func GenerateLicenseKey(API, appName, secretKey, expiryDelimeter string, expiredInDays int, payLoad ...interface{}) (bool, error) {
	// Get the disk serial number
	diskSerialNo, err := itrdsn.GetDiskSerialNo()
	if err != nil {
		return false, errors.New("error getting disk serial number: " + err.Error())
	}

	// Create a license file if not exist with the '.license' custom file format.
	keyFile := strings.ToLower(appName) + _fileExt
	f, err := os.OpenFile(keyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Add license expiry date using the '+days' from the current system datetime.
	expiredDays := time.Now().AddDate(0, 0, expiredInDays).Unix()
	strExpiredDate := fmt.Sprintf("%v", expiredDays)
	newLicenseKey := ""

	// Set expiry delimeter
	if len(strings.TrimSpace(expiryDelimeter)) == 0 {
		expiryDelimeter = _defaultExpiryDelimiter
	}

	// Set extra expiry date info to be embeded for each license key.
	if expiredInDays > 0 {
		newLicenseKey = diskSerialNo + expiryDelimeter + strExpiredDate
	} else {
		newLicenseKey = diskSerialNo + expiryDelimeter + "none"
	}

	// Write a new license key to your 'appname.license' custom file.
	newLicenseKey, err = tago.Encrypt(newLicenseKey, secretKey)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(f.Name(), []byte(strings.TrimSpace(newLicenseKey)), 0644)
	if err != nil {
		return false, err
	}

	// Check if the API endpoint has been activated or not
	if len(strings.TrimSpace(API)) > 0 {
		payLoad = append(payLoad, newLicenseKey) // Add additional value to our payLoad variadic interface parameter.
		api, err := APISubmitNewLicenseKey(API, payLoad)
		if err != nil {
			return false, err
		}
		if api {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

// APISubmitNewLicenseKey to submit new license key information to any backend API endpoint.
func APISubmitNewLicenseKey(API string, payLoad []interface{}) (bool, error) {
	// Check if the API endpoint has been activated or not
	if len(strings.TrimSpace(API)) > 0 {
		api, err := APIEndPoint(API, payLoad)
		if err != nil {
			return false, err
		}
		if api {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

// RevokeLicenseKey revokes existing gokopy license key.
func RevokeLicenseKey(API, appName string, payLoad ...interface{}) (bool, error) {
	// Check if the API endpoint has been activated or not
	if len(strings.TrimSpace(API)) > 0 {
		api, err := APIEndPoint(API, payLoad)
		if err != nil {
			return false, err
		}
		if api {
			// Clear the license key file as well
			if err := ClearLicenseKeyFile(appName); err != nil {
				return false, err
			}
			return true, nil
		}
	} else {
		// Clear the license key file as well
		if err := ClearLicenseKeyFile(appName); err != nil {
			return false, err
		}
	}
	return true, nil
}

// APIEndPoint is the backend endpoint for your application when required to be triggered.
func APIEndPoint(API string, payLoad []interface{}) (bool, error) {
	// Compose the JSON post payload to the API endpoint.
	var extractPayLoad []interface{} = payLoad
	message := map[string]interface{}{}

	for c, v := range extractPayLoad {
		message[fmt.Sprintf("%v", c)] = v
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(API, "application/json", bytes.NewBuffer(bytesRepresentation))
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

// ReadLicenseKey reads the license key if found from a custom file, otherwise, throws an error.
func ReadLicenseKey(appName string) (string, error) {
	// Open the custom license file
	// Create a license file if not exist with the '.license' custom file format.
	licFile := strings.ToLower(appName) + _fileExt
	f, err := os.OpenFile(licFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	file, err := os.Open(strings.ToLower(appName) + _fileExt)
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

// IsLicenseKeyExpired ensure that the license key expiration date is still valid or not.
func IsLicenseKeyExpired(licenseKey, expiryDelimiter string) bool {
	var curTime int64 = time.Now().Unix()
	data := strings.Split(strings.TrimSpace(fmt.Sprint(licenseKey)), expiryDelimiter)
	for n, d := range data {
		if n == 1 {
			if d != "none" {
				if intUnixDate, err := strconv.Atoi(d); err == nil {
					if int64(intUnixDate) <= curTime {
						return true
					}
				}
			}
		}
	}
	return false
}

// ExtractLicenseKey extract the license key without the expiry date.
func ExtractLicenseKey(licenseKey, expiryDelimiter string) string {
	data := strings.Split(strings.TrimSpace(fmt.Sprint(licenseKey)), expiryDelimiter)
	extractedLicenseKey := ""
	for n, d := range data {
		if n == 0 {
			extractedLicenseKey = d
		}
	}
	return extractedLicenseKey
}

// IsLicenseKeyValid is the main license key's validation.
func IsLicenseKeyValid(appName, secretKey, expiryDelimiter string) (bool, error) {
	// Get the user's primary hard disk serial number
	diskSerialNo, err := itrdsn.GetDiskSerialNo()
	if err != nil {
		return false, err
	}

	// Get the license key from a file
	fileDiskSerialNo := ""
	licenseKey, err := ReadLicenseKey(appName)
	if err != nil {
		return false, err
	}

	fileDiskSerialNo, err = tago.Decrypt(licenseKey, secretKey)
	if err != nil {
		return false, errors.New("error decrypting license key: " + err.Error())
	}

	// Check if the license key is expired or not
	if IsLicenseKeyExpired(fileDiskSerialNo, expiryDelimiter) {
		return false, errors.New("license key has been expired")
	}

	// Extract the license key without expiry date information
	extractedLicenseKey := ExtractLicenseKey(fileDiskSerialNo, expiryDelimiter)

	// Check both disk serial number must match
	if strings.TrimSpace(extractedLicenseKey) != strings.TrimSpace(diskSerialNo) {
		return false, errors.New("either license key file is empty or invalid")
	}
	return true, nil
}

// ClearLicenseKeyFile remove license key from a file.
func ClearLicenseKeyFile(appName string) error {
	// Create a license file if not exist with the '.license' custom file format.
	keyFile := strings.ToLower(appName) + _fileExt
	f, err := os.OpenFile(keyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Now, empty the license key from a file
	err = ioutil.WriteFile(f.Name(), []byte(strings.TrimSpace("")), 0644)
	if err != nil {
		return err
	}
	return nil
}
