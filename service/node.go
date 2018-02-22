package service

import (
	"golang.org/x/net/context"
	"io/ioutil"
	"github.com/container-storage-interface/spec/lib/go/csi"
	goofys "github.com/kahing/goofys/api"
	"github.com/jacobsa/fuse"
	"fmt"
	"os"
)

const QingCloudInstanceIDFilePath = "/etc/qingcloud/instance-id"

//mount on user namespace
func (s *service) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	//mount on user namespace
	accessKey:=req.UserCredentials["QY_ACCESS_KEY"]
	secretKey:=req.UserCredentials["QY_SECRET_KEY"]
	zone:= req.UserCredentials["QY_ZONE"]
	volumeName:= req.VolumeId
	endpoint:= fmt.Sprintf("https://%s.s3.%s.qingstor.com/mybucket/mykey",volumeName,zone)
	os.Setenv("AWS_ACCESS_KEY_ID",accessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY",secretKey)
	config := goofys.Config{
		MountPoint: req.TargetPath,
		DirMode:    0755,
		FileMode:   0644,
		Endpoint: endpoint,
		Region: zone,
	}

	_, mp, err := goofys.Mount(context.Background(), volumeName, &config)
	if err != nil {
		panic(fmt.Sprintf("Unable to mount %v: %v", config.MountPoint, err))
	} else {
		mp.Join(context.Background())
	}
	return nil, nil
}

func (s *service) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {
	fuse.Unmount(req.TargetPath)
	return nil, nil
}


//GetNodeID get node id in qingcloud zone.
// node id is stored under QingCloudInstanceIDFilePath when vm is started
func (s *service) GetNodeID(
	ctx context.Context,
	req *csi.GetNodeIDRequest) (
	*csi.GetNodeIDResponse, error) {
	idFile,err:=ioutil.ReadFile(QingCloudInstanceIDFilePath)
	if err != nil {
		return nil,err
	} else {
		return &csi.GetNodeIDResponse{
			NodeId: string(idFile),
		},nil
	}
}

func (s *service) NodeProbe(
	ctx context.Context,
	req *csi.NodeProbeRequest) (
	*csi.NodeProbeResponse, error) {

	return &csi.NodeProbeResponse{}, nil
}

func (s *service) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			&csi.NodeServiceCapability{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_UNKNOWN,
					},
				},
			},
		},
	}, nil}
