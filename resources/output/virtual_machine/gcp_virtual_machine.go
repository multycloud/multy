package virtual_machine

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/common"
	"golang.org/x/exp/slices"
	"strings"
)

type GoogleComputeInstance struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_compute_instance"`

	MachineType           string                                `hcl:"machine_type" json:"machine_type"`
	BootDisk              []GoogleBootDisk                      `hcl:"boot_disk,blocks" json:"boot_disk"`
	Zone                  string                                `hcl:"zone" hcle:"omitempty" json:"zone"`
	Tags                  []string                              `hcl:"tags" hcle:"omitempty" json:"tags"`
	MetadataStartupScript string                                `hcl:"metadata_startup_script" hcle:"omitempty" json:"metadata_startup_script"`
	NetworkInterface      []GoogleNetworkInterface              `hcl:"network_interface,blocks" json:"network_interface" json:"network_interface"`
	Metadata              map[string]string                     `hcl:"metadata" hcle:"omitempty" json:"metadata"`
	ServiceAccount        []GoogleComputeInstanceServiceAccount `hcl:"service_account,blocks" json:"service_account"`
}

type GoogleBootDisk struct {
	InitializeParams GoogleBootDiskInitializeParams `hcl:"initialize_params"`
}

type GoogleBootDiskInitializeParams struct {
	Image string `hcl:"image"`
	Size  int    `hcl:"size" hcle:"omitempty"`
}

type GoogleNetworkInterface struct {
	Subnetwork   string                               `hcl:"subnetwork,expr"`
	AccessConfig []GoogleNetworkInterfaceAccessConfig `hcl:"access_config,blocks" json:"access_config"`
}

type GoogleNetworkInterfaceAccessConfig struct {
	NetworkTier string `hcl:"network_tier"  hcle:"omitempty"` //STANDARD

	// outputs
	NatIp string `hcl:"nat_ip,expr" hcle:"omitempty" json:"nat_ip"`
}

type GoogleComputeInstanceServiceAccount struct {
	Email  string   `hcl:"email,expr" json:"email"`
	Scopes []string `hcl:"scopes" json:"scopes"`
}

func GetLatestGcpImage(ref *resourcespb.ImageReference) (string, error) {
	switch ref.Os {
	case resourcespb.ImageReference_UBUNTU:
		ubuntuVersions := []string{"18.04", "20.04", "22.04"}
		if !slices.Contains(ubuntuVersions, ref.Version) {
			return "", fmt.Errorf("ubuntu version %s is not supported in GCP", ref.Version)
		}
		version := fmt.Sprintf("%s-lts", strings.Replace(ref.Version, ".", "", 1))
		return fmt.Sprintf("ubuntu-os-cloud/ubuntu-%s", version), nil
	case resourcespb.ImageReference_DEBIAN:
		debianVersions := []string{"9", "10", "11"}
		if !slices.Contains(debianVersions, ref.Version) {
			return "", fmt.Errorf("debian version %s is not supported in GCP", ref.Version)
		}
		return fmt.Sprintf("debian-cloud/debian-%s", ref.Version), nil
	case resourcespb.ImageReference_CENT_OS:
		if ref.Version == "7" {
			return "centos-cloud/centos-7", nil
		}
		centosVersions := []string{"8", "9"}
		if !slices.Contains(centosVersions, ref.Version) {
			return "", fmt.Errorf("centos version %s is not supported in GCP", ref.Version)
		}
		return fmt.Sprintf("centos-cloud/centos-stream-%s", ref.Version), nil
	default:
		return "", fmt.Errorf("unknown operating system distibution %s", ref.Os)
	}
}
