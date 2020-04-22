// Package lisensya is the simple license key generator tool for your Go's applications.
package lisensya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/itrepablik/tago"
)

const (
	_fileExt                = ".license"
	_defaultExpiryDelimiter = ";expiry:"
)

// GenerateLicenseKey writes the new license key to a custom file and stores in the root directory of your
// app directory, e.g appname.license
func GenerateLicenseKey(licenseKey, appName, secretKey, expiryDelimeter string, expiredInDays int) (string, error) {
	// Create a license file if not exist with the '.license' custom file format.
	keyFile := strings.ToLower(appName) + _fileExt
	f, err := os.OpenFile(keyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
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
		newLicenseKey = licenseKey + expiryDelimeter + strExpiredDate
	} else {
		newLicenseKey = licenseKey + expiryDelimeter + "none"
	}

	// Write a new license key to your 'appname.license' custom file.
	newLicenseKey, err = tago.Encrypt(newLicenseKey, secretKey)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(f.Name(), []byte(strings.TrimSpace(newLicenseKey)), 0644)
	if err != nil {
		return "", err
	}
	return newLicenseKey, nil
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
func IsLicenseKeyExpired(licenseKey, expiryDelimeter string) bool {
	var curTime int64 = time.Now().Unix()
	data := strings.Split(strings.TrimSpace(fmt.Sprint(licenseKey)), expiryDelimeter)
	for n, d := range data {
		if n == 1 {
			if d != "none" {
				if intUnixDate, err := strconv.Atoi(d); err == nil {
					if curTime >= int64(intUnixDate) {
						return true
					}
				}
			}
		}
	}
	return false
}
