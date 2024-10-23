package driver

import (
	"context"
	"fmt"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	kubeAPI "k8s.io/api/core/v1"
	"k8s.io/mount-utils"

	csi_util "github.com/oracle/oci-cloud-controller-manager/pkg/csi-util"
	"github.com/oracle/oci-cloud-controller-manager/pkg/util/disk"
)

const (
	SetupLnet        = "setupLnet"
	LustreSubnetCidr = "lustreSubnetCidr"
)

func (d LustreNodeDriver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {

	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "Volume ID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Staging path must be provided")
	}

	accessType := req.VolumeCapability.GetMount()
	if accessType == nil || accessType.FsType != "lustre" {
		return nil, status.Error(codes.InvalidArgument, "Invalid fsType provided. Only \"lustre\" fsType is supported on this driver.")
	}

	isValidVolumeId, lnetLabel := csi_util.ValidateLustreVolumeId(req.VolumeId)
	if !isValidVolumeId {
		return nil, status.Error(codes.InvalidArgument, "Invalid Volume Handle provided.")
	}

	logger := d.logger.With("volumeID", req.VolumeId)

	logger.Debugf("volume context: %v", req.VolumeContext)

	if acquired := d.volumeLocks.TryAcquire(req.VolumeId); !acquired {
		logger.Error("Could not acquire lock for NodeStageVolume.")
		return nil, status.Errorf(codes.Aborted, volumeOperationAlreadyExistsFmt, req.VolumeId)
	}
	defer d.volumeLocks.Release(req.VolumeId)

	//Lnet Setup
	if setupLnet, ok := req.GetVolumeContext()[SetupLnet]; ok && setupLnet == "true" {

		lustreSubnetCIDR, ok :=  req.GetVolumeContext()[LustreSubnetCidr]

		if !ok {
			lustreSubnetCIDR = fmt.Sprintf("%s/32", d.nodeID)
		}

		err := csi_util.SetupLnet(logger, lustreSubnetCIDR, lnetLabel)
		if err != nil {
			logger.With(zap.Error(err)).Error("Failed to setup lnet.")
			return nil, status.Errorf(codes.Internal, "Failed to setup lnet with error : %v", err.Error())
		}

	} else {
		logger.Info("Lnet setup skipped as it is disabled in PV. ")
	}

	var fsType = accessType.FsType

	var options []string
	if accessType.MountFlags != nil {
		options = accessType.MountFlags
	}

	mounter := mount.New(mountPath)



	targetPath := req.StagingTargetPath
	mountPoint, err := isMountPoint(mounter, targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.With("StagingTargetPath", targetPath).Infof("Mount point does not pre-exist, creating now.")
			// k8s v1.20+ will not create the TargetPath directory
			// https://github.com/kubernetes/kubernetes/pull/88759
			// if the path exists already (<v1.20) this is a no op
			// https://golang.org/pkg/os/#MkdirAll
			if err := os.MkdirAll(targetPath, 0750); err != nil {
				logger.With(zap.Error(err)).Error("Failed to create StagingTargetPath directory")
				return nil, status.Error(codes.Internal, "Failed to create StagingTargetPath directory")
			}
			mountPoint = false
		} else {
			logger.With(zap.Error(err)).Error("Invalid Mount Point ", targetPath)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if mountPoint {
		logger.Infof("Volume is already mounted to: %v", targetPath)
		return &csi.NodeStageVolumeResponse{}, nil
	}

	source := req.VolumeId

	err = mounter.Mount(source, targetPath, fsType, options)

	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to mount volume to staging target path.")
		return nil, status.Errorf(codes.Internal, "Failed to mount volume to staging target path with error : %v", err.Error())
	}
	logger.With("Source", source, "StagingTargetPath", targetPath).
		Info("Mounting the volume to staging target path is completed.")

	return &csi.NodeStageVolumeResponse{}, nil
}
func (d LustreNodeDriver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "Volume ID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Staging path must be provided")
	}

	isValidVolumeId, lnetLabel := csi_util.ValidateLustreVolumeId(req.VolumeId)
	if !isValidVolumeId {
		return nil, status.Error(codes.InvalidArgument, "Invalid Volume Handle provided.")
	}

	logger := d.logger.With("volumeID", req.VolumeId, "stagingPath", req.StagingTargetPath)

	if acquired := d.volumeLocks.TryAcquire(req.VolumeId); !acquired {
		logger.Error("Could not acquire lock for NodeUnstageVolume.")
		return nil, status.Errorf(codes.Aborted, volumeOperationAlreadyExistsFmt, req.VolumeId)
	}

	defer d.volumeLocks.Release(req.VolumeId)

	targetPath := req.GetStagingTargetPath()

	mounter := mount.New(mountPath)

	logger.Info("Unstage started")

	if !csi_util.IsLnetActive(logger, lnetLabel) {
		//When lnet is not active force unmount is required as regular unmounts get stuck forever.
		logger.Info("Performing force unmount as no active lnet configuration found.")
		if err := disk.UnmountWithForce(targetPath); err != nil {
			logger.With(zap.Error(err)).Error("Failed to unmount target path.")
			return nil, status.Errorf(codes.Internal, "Failed to unmount target path with error : %v", err.Error())
		}
		logger.With("StagingTargetPath", targetPath).Info("Unstage volume completed")
		return &csi.NodeUnstageVolumeResponse{}, nil
	}

	// Use mount.IsMountPoint because mounter.IsLikelyNotMountPoint can't detect bind mounts
	isMountPoint, err := mounter.IsMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.With("StagingTargetPath", targetPath).Infof("mount point does not exist")
			return &csi.NodeUnstageVolumeResponse{}, nil
		}
		return  nil, status.Error(codes.Internal, err.Error())
	}

	if !isMountPoint {
		logger.With("StagingTargetPath", targetPath).Infof("Not a mount point, removing path.")
		err = os.RemoveAll(targetPath)
		if err != nil {
			logger.With(zap.Error(err)).Error("Remove target path failed with error")
			return  nil, status.Error(codes.Internal, "Failed to remove target path")
		}
		return &csi.NodeUnstageVolumeResponse{}, nil
	}

	err = mounter.Unmount(targetPath)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to unmount  staging target path.")
		return nil, status.Errorf(codes.Internal, "Failed to unmount staging target path with error : %v", err.Error())
	}

	logger.With("StagingTargetPath", targetPath).Info("Unmounting volume completed")
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (d LustreNodeDriver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "Volume ID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Staging path must be provided")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Target Path must be provided")
	}


	logger := d.logger.With("volumeID", req.VolumeId)
	logger.Debugf("volume context: %v", req.VolumeContext)

	_, lnetLabel := csi_util.ValidateLustreVolumeId(req.VolumeId)

	//Lnet Setup
	if setupLnet, ok := req.GetVolumeContext()[SetupLnet]; ok && setupLnet == "true" {

		lustreSubnetCIDR, ok :=  req.GetVolumeContext()[LustreSubnetCidr]

		if !ok {
			lustreSubnetCIDR = fmt.Sprintf("%s/32", d.nodeID)
		}

		err := csi_util.SetupLnet(logger, lustreSubnetCIDR, lnetLabel)
		if err != nil {
			logger.With(zap.Error(err)).Error("Failed to setup lnet.")
			return nil, status.Errorf(codes.Internal, "Failed to setup lnet with error : %v", err.Error())
		}

	} else {
		logger.Info("Lnet setup skipped as it is disabled in PV. ")
	}

	var fsType = ""

	mounter := mount.New(mountPath)

	targetPath := req.GetTargetPath()
	readOnly := req.GetReadonly()
	// Use mount.IsNotMountPoint because mounter.IsLikelyNotMountPoint can't detect bind mounts
	isNotMountPoint, err := mount.IsNotMountPoint(mounter, targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.With("TargetPath", targetPath).Infof("mount point does not exist")
			// k8s v1.20+ will not create the TargetPath directory
			// https://github.com/kubernetes/kubernetes/pull/88759
			// if the path exists already (<v1.20) this is a no op
			// https://golang.org/pkg/os/#MkdirAll
			if err := os.MkdirAll(targetPath, 0750); err != nil {
				logger.With(zap.Error(err)).Error("Failed to create TargetPath directory")
				return nil, status.Error(codes.Internal, "Failed to create TargetPath directory")
			}
			isNotMountPoint = true
		} else {
			logger.With(zap.Error(err)).Error("Invalid Mount Point ", targetPath)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if !isNotMountPoint {
		logger.Infof("Volume is already mounted to: %v", targetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	options := []string{"bind"}
	if readOnly {
		options = append(options, "ro")
	}
	source := req.GetStagingTargetPath()
	err = mounter.Mount(source, targetPath, fsType, options)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to bind mount volume to target path.")
		return nil, status.Error(codes.Internal, err.Error())
	}
	logger.With("staging target path", source, "TargetPath", targetPath).
		Info("Bind mounting the volume to target path is completed.")

	return &csi.NodePublishVolumeResponse{}, nil
}

func (d LustreNodeDriver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeUnpublishVolume: Volume ID must be provided")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeUnpublishVolume: Target Path must be provided")
	}

	isValidVolumeId, lnetLabel := csi_util.ValidateLustreVolumeId(req.VolumeId)
	if !isValidVolumeId {
		return nil, status.Error(codes.InvalidArgument, "Invalid Volume Handle provided.")
	}

	logger := d.logger.With("volumeID", req.VolumeId, "targetPath", req.TargetPath)

	mounter := mount.New(mountPath)
	targetPath := req.GetTargetPath()

	logger.Info("Unmount started")

	if !csi_util.IsLnetActive(logger, lnetLabel) {
		//When lnet is not active force unmount is required as regular unmounts get stuck forever
		logger.Info("Performing force unmount as no active lnet configuration found.")
		if err := disk.UnmountWithForce(targetPath); err != nil {
			logger.With(zap.Error(err)).Error("Failed to unmount target path.")
			return nil, status.Errorf(codes.Internal, "Failed to unmount target path with error : %v", err.Error())
		}
		logger.With("TargetPath", targetPath).Info("Unmounting volume completed")
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}
	// Use mount.IsNotMountPoint because mounter.IsLikelyNotMountPoint can't detect bind mounts
	isMountPoint, err := mounter.IsMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.With("TargetPath", targetPath).Infof("mount point does not exist")
			return &csi.NodeUnpublishVolumeResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !isMountPoint {
		err = os.RemoveAll(targetPath)
		if err != nil {
			logger.With(zap.Error(err)).Error("Remove target path failed with error")
			return nil, status.Error(codes.Internal, "Failed to remove target path")
		}
		logger.With("TargetPath", targetPath).Infof("Not a mount point, removing path")
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}

	if err := mounter.Unmount(targetPath); err != nil {
		logger.With(zap.Error(err)).Error("Failed to unmount target path.")
		return nil, status.Errorf(codes.Internal, "Failed to unmount target path with error : %v", err.Error())
	}
	logger.With("TargetPath", targetPath).Info("Unmounting volume completed")
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats return the stats of the volume
func (d LustreNodeDriver) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeGetVolumeStats is not supported.")
}

// NodeExpandVolume returns the expand of the volume
func (d LustreNodeDriver) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeExpandVolume is not supported.")
}

func (d LustreNodeDriver) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	nscap := &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
			},
		},
	}

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			nscap,
		},
	}, nil
}

func (d LustreNodeDriver) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	ad, err := d.util.LookupNodeAvailableDomain(d.KubeClient, d.nodeID)
	if err != nil {
		d.logger.With(zap.Error(err)).With("nodeId", d.nodeID, "availableDomain", ad).Error("Available domain of node missing.")
	}

	d.logger.With("nodeId", d.nodeID, "availableDomain", ad).Info("Available domain of node identified.")
	return &csi.NodeGetInfoResponse{
		NodeId: d.nodeID,
		// make sure that the driver works on this particular AD only
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				kubeAPI.LabelZoneFailureDomain: ad,
			},
		},
	}, nil
}

