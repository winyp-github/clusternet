/*
Copyright 2021 The Clusternet Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generic

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	appsapi "github.com/clusternet/clusternet/pkg/apis/apps/v1alpha1"
	"github.com/clusternet/clusternet/pkg/controllers/apps/description"
	clusternetclientset "github.com/clusternet/clusternet/pkg/generated/clientset/versioned"
	clusternetinformers "github.com/clusternet/clusternet/pkg/generated/informers/externalversions"
	clusterlisters "github.com/clusternet/clusternet/pkg/generated/listers/clusters/v1beta1"
	"github.com/clusternet/clusternet/pkg/known"
	"github.com/clusternet/clusternet/pkg/utils"
)

type Deployer struct {
	ctx context.Context

	clusterLister clusterlisters.ManagedClusterLister
	clusterSynced cache.InformerSynced
	secretLister  corev1lister.SecretLister
	secretSynced  cache.InformerSynced

	clusternetClient *clusternetclientset.Clientset

	descController *description.Controller

	recorder record.EventRecorder
}

func NewDeployer(ctx context.Context, clusternetClient *clusternetclientset.Clientset,
	clusternetInformerFactory clusternetinformers.SharedInformerFactory, kubeInformerFactory kubeinformers.SharedInformerFactory,
	recorder record.EventRecorder) (*Deployer, error) {

	deployer := &Deployer{
		ctx:              ctx,
		clusterLister:    clusternetInformerFactory.Clusters().V1beta1().ManagedClusters().Lister(),
		clusterSynced:    clusternetInformerFactory.Clusters().V1beta1().ManagedClusters().Informer().HasSynced,
		secretLister:     kubeInformerFactory.Core().V1().Secrets().Lister(),
		secretSynced:     kubeInformerFactory.Core().V1().Secrets().Informer().HasSynced,
		clusternetClient: clusternetClient,
		recorder:         recorder,
	}

	descController, err := description.NewController(ctx,
		clusternetClient,
		clusternetInformerFactory.Apps().V1alpha1().Descriptions(),
		clusternetInformerFactory.Apps().V1alpha1().HelmReleases(),
		deployer.recorder,
		deployer.handleDescription)
	if err != nil {
		return nil, err
	}
	deployer.descController = descController

	return deployer, nil
}

func (deployer *Deployer) Run(workers int) {
	klog.Info("starting generic deployer...")
	defer klog.Info("shutting generic deployer")

	// Wait for the caches to be synced before starting workers
	klog.V(5).Info("waiting for informer caches to sync")
	if !cache.WaitForCacheSync(deployer.ctx.Done(),
		deployer.clusterSynced,
		deployer.secretSynced) {
		return
	}

	go deployer.descController.Run(workers, deployer.ctx.Done())

	<-deployer.ctx.Done()
}

func (deployer *Deployer) handleDescription(desc *appsapi.Description) error {
	klog.V(5).Infof("handle Description %s", klog.KObj(desc))
	if desc.Spec.Deployer != appsapi.DescriptionGenericDeployer {
		return nil
	}

	deployable, err := utils.DeployableByHub(deployer.clusterLister, desc.Labels[known.ClusterIDLabel], desc.Namespace)
	if err != nil {
		klog.ErrorDepth(4, err)
		deployer.recorder.Event(desc, corev1.EventTypeWarning, "ManagedClusterNotFound", err.Error())
		return err
	}
	if !deployable {
		klog.V(5).Infof("Description %s is not deployable by hub, skipping syncing", klog.KObj(desc))
		return nil
	}

	dynamicClient, discoveryRESTMapper, err := deployer.getDynamicClient(desc)
	if err != nil {
		return err
	}

	if desc.DeletionTimestamp != nil {
		return utils.OffloadDescription(deployer.ctx, deployer.clusternetClient, dynamicClient,
			discoveryRESTMapper, desc, deployer.recorder)
	}

	return utils.ApplyDescription(deployer.ctx, deployer.clusternetClient, dynamicClient,
		discoveryRESTMapper, desc, deployer.recorder)
}

func (deployer *Deployer) getDynamicClient(desc *appsapi.Description) (dynamic.Interface, meta.RESTMapper, error) {
	config, err := utils.GetChildClusterConfig(deployer.secretLister, deployer.clusterLister, desc.Namespace, desc.Labels[known.ClusterIDLabel])
	if err != nil {
		return nil, nil, err
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	restConfig.QPS = 5
	restConfig.Burst = 10

	kubeclient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	discoveryClient := cacheddiscovery.NewMemCacheClient(kubeclient.Discovery())
	discoveryRESTMapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}

	return dynamicClient, discoveryRESTMapper, nil

}