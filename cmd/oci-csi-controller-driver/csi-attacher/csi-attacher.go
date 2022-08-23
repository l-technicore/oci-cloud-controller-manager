// Copyright 2019 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csiattacher

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/leaderelection"
	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	"github.com/kubernetes-csi/csi-lib-utils/rpc"
	"github.com/kubernetes-csi/external-attacher/pkg/attacher"
	"github.com/kubernetes-csi/external-attacher/pkg/controller"
	"github.com/oracle/oci-cloud-controller-manager/cmd/oci-csi-controller-driver/csioptions"
	"google.golang.org/grpc"
	csitranslationlib "k8s.io/csi-translation-lib"
)

const (

	// Default timeout of short CSI calls like GetPluginInfo
	csiTimeout = time.Second

	leaderElectionTypeLeases     = "leases"
	leaderElectionTypeConfigMaps = "configmaps"
)

var (
	version = "unknown"
)

type leaderElection interface {
	Run() error
	WithNamespace(namespace string)
}

//StartCSIAttacher main function to start CSI attacher
func StartCSIAttacher(csioptions csioptions.CSIOptions) {

	if csioptions.ShowVersion {
		fmt.Println(os.Args[0], version)
		return
	}
	klog.Infof("Version: %s", version)

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(csioptions.Kubeconfig)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	if csioptions.WorkerThreads == 0 {
		klog.Error("option -worker-threads must be greater than zero")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	factory := informers.NewSharedInformerFactory(clientset, csioptions.Resync)
	var handler controller.Handler

	metricsManager := metrics.NewCSIMetricsManager("" /* driverName */)

	// Connect to CSI.
	csiConn, err := connection.Connect(csioptions.CsiAddress, metricsManager, connection.OnConnectionLoss(connection.ExitOnConnectionLoss()))
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	err = rpc.ProbeForever(csiConn, csioptions.Timeout)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	// Find driver name.
	ctx, cancel := context.WithTimeout(context.Background(), csiTimeout)
	defer cancel()
	csiAttacher, err := rpc.GetDriverName(ctx, csiConn)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}
	klog.V(2).Infof("CSI driver name: %q", csiAttacher)

	mux := http.NewServeMux()
	addr := csioptions.MetricsAddress
	// Start HTTP server for metrics + leader election healthz
	if addr != "" {
		metricsManager.RegisterToServer(mux, csioptions.MetricsPath)
		metricsManager.SetDriverName(csiAttacher)
		go func() {
			klog.Infof("ServeMux listening at %q", addr)
			err := http.ListenAndServe(addr, mux)
			if err != nil {
				klog.Fatalf("Failed to start HTTP server at specified address (%q) and metrics path (%q): %s", addr, csioptions.MetricsPath, err)
			}
		}()
	}
	//metricsManager.SetDriverName(csiAttacher)
	//metricsManager.StartMetricsEndpoint(csioptions.MetricsAddress, csioptions.MetricsPath)

	supportsService, err := supportsPluginControllerService(ctx, csiConn)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}
	if !supportsService {
		handler = controller.NewTrivialHandler(clientset)
		klog.V(2).Infof("CSI driver does not support Plugin Controller Service, using trivial handler")
	} else {
		// Find out if the driver supports attach/detach.
		supportsAttach, supportsReadOnly, err := supportsControllerPublish(ctx, csiConn)
		if err != nil {
			klog.Error(err.Error())
			os.Exit(1)
		}
		if supportsAttach {
			pvLister := factory.Core().V1().PersistentVolumes().Lister()
			vaLister := factory.Storage().V1().VolumeAttachments().Lister()
			csiNodeLister := factory.Storage().V1().CSINodes().Lister()
			volAttacher := attacher.NewAttacher(csiConn)
			CSIVolumeLister := attacher.NewVolumeLister(csiConn)
			handler = controller.NewCSIHandler(clientset, csiAttacher, volAttacher, CSIVolumeLister, pvLister, csiNodeLister, vaLister, &csioptions.Timeout, supportsReadOnly, false, csitranslationlib.New(), csioptions.DefaultFSType)
			klog.V(2).Infof("CSI driver supports ControllerPublishUnpublish, using real CSI handler")
		} else {
			handler = controller.NewTrivialHandler(clientset)
			klog.V(2).Infof("CSI driver does not support ControllerPublishUnpublish, using trivial handler")
		}
	}

	slvpn, err := supportsListVolumesPublishedNodes(ctx, csiConn)
	if err != nil {
		klog.Errorf("Failed to check if driver supports ListVolumesPublishedNodes, assuming it does not: %v", err)
	}

	if slvpn {
		klog.V(2).Infof("CSI driver supports list volumes published nodes. Using capability to reconcile volume attachment objects with actual backend state")
	}

	ctrl := controller.NewCSIAttachController(
		clientset,
		csiAttacher,
		handler,
		factory.Storage().V1().VolumeAttachments(),
		factory.Core().V1().PersistentVolumes(),
		workqueue.NewItemExponentialFailureRateLimiter(csioptions.RetryIntervalStart, csioptions.RetryIntervalMax),
		workqueue.NewItemExponentialFailureRateLimiter(csioptions.RetryIntervalStart, csioptions.RetryIntervalMax),
		slvpn,
		csioptions.ReconcileSync,
	)

	run := func(ctx context.Context) {
		stopCh := ctx.Done()
		factory.Start(stopCh)
		ctrl.Run(int(csioptions.WorkerThreads), stopCh)
	}

	if !csioptions.EnableLeaderElection {
		run(context.TODO())
	} else {
		// Name of config map with leader election lock
		lockName := "external-attacher-leader-" + csiAttacher
		le := leaderelection.NewLeaderElection(clientset, lockName, run)

		if csioptions.LeaderElectionNamespace != "" {
			le.WithNamespace(csioptions.LeaderElectionNamespace)
		}

		if err := le.Run(); err != nil {
			klog.Fatalf("failed to initialize leader election: %v", err)
		}
	}
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func supportsControllerPublish(ctx context.Context, csiConn *grpc.ClientConn) (supportsControllerPublish bool, supportsPublishReadOnly bool, err error) {
	caps, err := rpc.GetControllerCapabilities(ctx, csiConn)
	if err != nil {
		return false, false, err
	}

	supportsControllerPublish = caps[csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME]
	supportsPublishReadOnly = caps[csi.ControllerServiceCapability_RPC_PUBLISH_READONLY]
	return supportsControllerPublish, supportsPublishReadOnly, nil
}

func supportsListVolumesPublishedNodes(ctx context.Context, csiConn *grpc.ClientConn) (bool, error) {
	caps, err := rpc.GetControllerCapabilities(ctx, csiConn)
	if err != nil {
		return false, fmt.Errorf("failed to get controller capabilities: %v", err)
	}

	return caps[csi.ControllerServiceCapability_RPC_LIST_VOLUMES] && caps[csi.ControllerServiceCapability_RPC_LIST_VOLUMES_PUBLISHED_NODES], nil
}

func supportsPluginControllerService(ctx context.Context, csiConn *grpc.ClientConn) (bool, error) {
	caps, err := rpc.GetPluginCapabilities(ctx, csiConn)
	if err != nil {
		return false, err
	}

	return caps[csi.PluginCapability_Service_CONTROLLER_SERVICE], nil
}
