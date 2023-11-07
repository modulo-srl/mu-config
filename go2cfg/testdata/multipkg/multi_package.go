package multipkg

import (
	"github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network"
	alias "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/stats"
)

//go:generate go2cfg -type MultiPackage -out multi_package.jsonc
//go:generate go2cfg -type MultiPackage -doc-types all -out multi_package_all_fields.jsonc
//go:generate go2cfg -type MultiPackage -doc-types basic -out multi_package_basic_fields.jsonc

// MultiPackage tests the multi-package and import aliasing case.
type MultiPackage struct {
	NetStatus  network.Status // Network status.
	alias.Info                // Statistics info.
}

func MultiPackageDefaults() *MultiPackage {
	return &MultiPackage{
		NetStatus: network.Status{
			Connected: true,
			State:     network.StateDisconnected,
		},
		Info: alias.Info{
			PacketLoss:    32 * 2,
			RoundTripTime: 123,
		},
	}
}
