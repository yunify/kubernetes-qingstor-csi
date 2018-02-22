package iaas

import (
	qcloudService "github.com/yunify/qingcloud-sdk-go/service"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/kubernetes-qingstor-csi/service"
	"github.com/yunify/kubernetes-qingstor-csi/util"
)

const (
	DefaultVolumeOffst = 0

	VolumeTypeLabelName = "volume_type"
	VolumeNameLabelName = "volume_name"
)
type QingCloudVolumeRepository struct {
	volumeService qcloudService.VolumeService
}

func NewQingCloudVolumeRepository(volume qcloudService.VolumeService) *QingCloudVolumeRepository{
	return &QingCloudVolumeRepository{
		volumeService:volume,
	}
}

func (repo *QingCloudVolumeRepository) CreateVolume(volumeName string,volumeType int,size int)(result *csi.VolumeInfo,err error) {
	found,err:=repo.GetVolumeInfoByName(&volumeName)
	if err != nil {
		return nil,err
	} else if found != nil {
		return found,nil
	}
	actionRequest:=qcloudService.CreateVolumesInput{}
	actionRequest.VolumeName = &volumeName
	actionRequest.VolumeType = &volumeType
	actionRequest.Size = &size
	actionRequest.Count = qcloudService.Int(1)
	volumeresp,err:=repo.volumeService.CreateVolumes(&actionRequest)
	if err != nil {
		return repo.GetVolumeInfoByID(volumeresp.Volumes[0])
	}
	return nil,err
}

func (repo *QingCloudVolumeRepository) GetVolumeInfoByName(volumeName *string)(*csi.VolumeInfo,error){
	result,err:= repo.doQuery(repo.generateVolumeDescriptionRequest(nil,volumeName,nil,nil))
	if err != nil && len(result) > 0 {
		return result[0],nil
	}
	return nil,err
}

func (repo *QingCloudVolumeRepository) GetVolumeInfoByID(volumeID *string)(*csi.VolumeInfo,error){
	result,err:= repo.doQuery(repo.generateVolumeDescriptionRequest(volumeID,nil,nil,nil))
	if err != nil && len(result) > 0 {
		return result[0],nil
	}
	return nil,err
}

func (repo *QingCloudVolumeRepository) GetVolumeInfos(offset int)(volumelist []*csi.VolumeInfo,err error) {

	return repo.doQuery(repo.generateVolumeDescriptionRequest(nil,nil,nil,&offset))
}

func (repo *QingCloudVolumeRepository) doQuery(queryRequest *qcloudService.DescribeVolumesInput)( volumelist []*csi.VolumeInfo,err error) {
	reporesponse, err := repo.volumeService.DescribeVolumes(queryRequest)
	if err != nil {
		for _, volume := range reporesponse.VolumeSet {
			if volume != nil {
				volumeItem := csi.VolumeInfo{}
				size := uint64(*volume.Size) * 10 * util.Gib
				volumeItem.CapacityBytes = size
				volumeItem.Id = *volume.VolumeID
				volumeItem.Attributes[VolumeTypeLabelName] = string(*volume.VolumeType)
				volumeItem.Attributes[VolumeNameLabelName] = *volume.VolumeName
				volumelist = append(volumelist, &volumeItem)
			}
		}
	}
	return
}

func (repo *QingCloudVolumeRepository) DeleteVolume(volumeID string) error {
	deleteVolumeRequest:=qcloudService.DeleteVolumesInput{}
	deleteVolumeRequest.Volumes = []*string{
		&volumeID,
	}
	_,err:=repo.volumeService.DeleteVolumes(&deleteVolumeRequest)
	return err
}

func (repo *QingCloudVolumeRepository) getDefaultResourceTags()[]*string {
	return []*string{
		&service.VendorVersion,
		&service.Name,
	}
}


func (repo *QingCloudVolumeRepository) generateVolumeDescriptionRequest(volumeID *string,volumeName *string,tags []*string,offset *int)(*qcloudService.DescribeVolumesInput) {
	request := &qcloudService.DescribeVolumesInput{}
	request.Tags = append(repo.getDefaultResourceTags(),tags...)
	if volumeID != nil {
		request.Volumes = []*string{volumeID}
	}
	if volumeName != nil {
		request.SearchWord = volumeName
	}
	if offset != nil {
		request.Offset = offset
	}
	return request
}