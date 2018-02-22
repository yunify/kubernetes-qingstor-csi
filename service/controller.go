package service

import (
	"golang.org/x/net/context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	qs "github.com/yunify/qingstor-sdk-go/service"
	"github.com/yunify/qingstor-sdk-go/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {
	cfg,_:= config.NewDefault();
	accessKey:=req.UserCredentials["QY_ACCESS_KEY"]
	secretKey:=req.UserCredentials["QY_SECRET_KEY"]
	zone:= req.UserCredentials["QY_ZONE"]
	volumeName:= req.Name
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretKey
	qsService,_:= qs.Init(cfg)
	bucket,err:=qsService.Bucket(volumeName,zone)
	if err !=nil {
		return nil, err
	}
	return &csi.CreateVolumeResponse{
		VolumeInfo: &csi.VolumeInfo{
			CapacityBytes: 0,
			Id: *bucket.Properties.BucketName,
		},
	}, nil
}

func (s *service) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {
	cfg,_:= config.NewDefault();
	accessKey:=req.UserCredentials["QY_ACCESS_KEY"]
	secretKey:=req.UserCredentials["QY_SECRET_KEY"]
	zone:= req.UserCredentials["QY_ZONE"]
	volumeName:= req.VolumeId
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretKey
	qsService,_:= qs.Init(cfg)
	bucket,err:=qsService.Bucket(volumeName,zone)
	if err !=nil {
		return nil, err
	}
	bOutput, err := bucket.ListObjects(nil)
	if err != nil {
		return nil,err
	}
	for _, item := range bOutput.Keys{
		bucket.DeleteObject(*item.Key)
	}
	bucket.Delete()
	return nil, nil
}

func (s *service) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {
	cfg,_:= config.NewDefault();
	accessKey:=req.UserCredentials["QY_ACCESS_KEY"]
	secretKey:=req.UserCredentials["QY_SECRET_KEY"]
	zone:= req.UserCredentials["QY_ZONE"]
	volumeName:= req.VolumeId
	cfg.AccessKeyID = accessKey
	cfg.SecretAccessKey = secretKey
	qsService,_:=qs.Init(cfg)
	_,err:=qsService.Bucket(volumeName,zone)
	if err != nil {
		return nil,err
	}
	return &csi.ControllerPublishVolumeResponse{
		PublishVolumeInfo: map[string]string{

		},
	}, nil
}

func (s *service) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (s *service) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	return nil, nil
}

func (s *service) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	return nil, status.Error(codes.Unimplemented, "")
}

func (s *service) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {

	return nil, status.Error(codes.Unimplemented, "")
}

func (s *service) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	return  &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
					},
				},
			},
		},
	}, nil
}

func (s *service) ControllerProbe(
	ctx context.Context,
	req *csi.ControllerProbeRequest) (
	*csi.ControllerProbeResponse, error) {

	return &csi.ControllerProbeResponse{}, nil
}
