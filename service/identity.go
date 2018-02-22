package service

import (
	"golang.org/x/net/context"
	"github.com/container-storage-interface/spec/lib/go/csi"
)

func (s *service) GetSupportedVersions(
	ctx context.Context,
	req *csi.GetSupportedVersionsRequest) (
	*csi.GetSupportedVersionsResponse, error) {
	return &csi.GetSupportedVersionsResponse{
		SupportedVersions: []*csi.Version{
			&csi.Version{Major: 0, Minor: 1, Patch: 0},
		},
	}, nil
}

func (s *service) GetPluginInfo(
	ctx context.Context,
	req *csi.GetPluginInfoRequest) (
	*csi.GetPluginInfoResponse, error) {

	return &csi.GetPluginInfoResponse{
		Name: Name,
		VendorVersion: VendorVersion,
	}, nil
}
