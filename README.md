![Lisensya](https://user-images.githubusercontent.com/58651329/79626813-ae0d2b00-8165-11ea-9c64-0419b7b91ece.png)
The **lisensya** package is a simple license key generator that manages your license key requirements for your Go's application.

# Installation
```
go get -u github.com/itrepablik/lisensya
```

# Usage
For this example, this is how you can use the main function to check your license key which is most likely in your main.go file.
```
package main

import (
	"fmt"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/lisensya"
)

// Keep this secret
const (
	secretKey              = "abc&1*~#^2^#s0^=)^^7%b34"
	appName                = "NiceApp"
	licenseExpiryDelimiter = "--expiry:"
)

// IsLicenseKeyValid is declared globally at your main.go check either the license key is expired or not.
var IsLicenseKeyValid bool = false

func main() {
	// This is the main validation for the license key.
	IsLicenseKeyValid, err := lisensya.IsLicenseKeyValid(appName, secretKey, licenseExpiryDelimiter)
	if err != nil {
		msg := `error getting license key: `
		fmt.Println(msg, err)
		itrlog.Fatal(msg, err)
	}

	// Check if the license key is valid or not.
	if !IsLicenseKeyValid {
		msg := `oops, invalid license key`
		itrlog.Error(msg)
		fmt.Println(msg)
		return
	}
}
```

This is to generate a new license key:
```
package main

import (
	"fmt"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/lisensya"
)

// Keep this secret
const (
	secretKey              = "abc&1*~#^2^#s0^=)^^7%b34"
	appName                = "NiceApp"
	licenseExpiryDelimiter = "--expiry:"
	APInewlicense          = "" // e.g. https://your_site_name.com/api/your_api_url
	expiredInDays          = 30
)

// IsLicenseKeyValid is declared globally at your main.go check either the license key is expired or not.
var IsLicenseKeyValid bool = false
var hostName string = ""

func main() {
	// Get the hostname
	hostName, err := lisensya.GetHostName()
	if err != nil {
		itrlog.Fatal("error getting hostname: ", err)
	}

	// Starts writing a new license key
	userName := "username"
	isNewLicenseOK, err := lisensya.GenerateLicenseKey(APInewlicense, appName,
		secretKey, licenseExpiryDelimiter, expiredInDays, hostName, userName, secretKey)

	if err != nil {
		itrlog.Fatal("error generating new license key: ", err)
	}

	// Inform user about the new license key successfully generated.
	if isNewLicenseOK {
		fmt.Println("You've successfully generated your gokopy's license key, you can now use this software.")
	}
}
```

This is to revoke existing license key:
```
package main

import (
	"fmt"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/lisensya"
)

// Keep this secret
const (
	secretKey              = "abc&1*~#^2^#s0^=)^^7%b34"
	appName                = "NiceApp"
	licenseExpiryDelimiter = "--expiry:"
	APInewlicense          = "" // e.g. https://your_site_name.com/api/your_api_url
	expiredInDays          = 30
)

// IsLicenseKeyValid is declared globally at your main.go check either the license key is expired or not.
var IsLicenseKeyValid bool = false
var hostName string = ""

func main() {
	userName := "username"

	// Get the existing license key from a file.
	licenseKey, err := lisensya.ReadLicenseKey(appName)
	if err != nil {
		itrlog.Fatal("error getting current license key: ", err)
	}

	// Update the license key at the backend
	rLicKey, err := lisensya.RevokeLicenseKey(APInewlicense, appName, hostName, licenseKey, userName, secretKey)
	if err != nil {
		itrlog.Fatal("error revoking your existing license key: ", err)
	}

	if rLicKey {
		fmt.Println("Successfully revoked your current gokopy's license key.")
	} else {
		fmt.Println("Oops!, so far probably you don't have any existing license key to be revoked.")
	}
}
```

# License
Code is distributed under MIT license, feel free to use it in your proprietary projects as well.
