package common

import (
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
)

var VMSIZE = map[string]map[CloudProvider]string{
	MICRO: {
		AWS:   "t2.nano",
		AZURE: "Standard_B1ls",
	},
	MEDIUM: {
		AWS:   "t2.medium",
		AZURE: "Standard_A2_v2",
	},
}

const (
	// MICRO - 1 core and 0.5 gb ram
	MICRO = "nano"
	// MEDIUM - 2 cores and 4 gb ram
	MEDIUM = "medium"
)

var DBSIZE = map[string]map[CloudProvider]string{
	MICRO: {
		AWS:   "db.t2.micro",
		AZURE: "GP_Gen5_2",
	},
}

func CheckIfSizeIsValid(size string) error {
	if !slices.Contains(maps.Keys(VMSIZE), size) {
		return fmt.Errorf("%s is not a valid vm size, supported sizes are: %s", size, strings.Join(maps.Keys(VMSIZE), ", "))
	}
	return nil
}
