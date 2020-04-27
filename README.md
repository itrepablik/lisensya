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

# License
Code is distributed under MIT license, feel free to use it in your proprietary projects as well.
