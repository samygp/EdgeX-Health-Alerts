package version

import (
	"fmt"
)

// ID of application.
var ID string

// Name of application.
var Name string

// Version indicates which version of the binary is running.
var Version string

func init() {
	ID = fmt.Sprintf("localhost:%s@%s", Name, Version)
}
