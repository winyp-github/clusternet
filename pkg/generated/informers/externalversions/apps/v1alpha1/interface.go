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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/clusternet/clusternet/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Bases returns a BaseInformer.
	Bases() BaseInformer
	// Descriptions returns a DescriptionInformer.
	Descriptions() DescriptionInformer
	// Globalizations returns a GlobalizationInformer.
	Globalizations() GlobalizationInformer
	// HelmCharts returns a HelmChartInformer.
	HelmCharts() HelmChartInformer
	// HelmReleases returns a HelmReleaseInformer.
	HelmReleases() HelmReleaseInformer
	// Localizations returns a LocalizationInformer.
	Localizations() LocalizationInformer
	// Manifests returns a ManifestInformer.
	Manifests() ManifestInformer
	// Subscriptions returns a SubscriptionInformer.
	Subscriptions() SubscriptionInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Bases returns a BaseInformer.
func (v *version) Bases() BaseInformer {
	return &baseInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Descriptions returns a DescriptionInformer.
func (v *version) Descriptions() DescriptionInformer {
	return &descriptionInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Globalizations returns a GlobalizationInformer.
func (v *version) Globalizations() GlobalizationInformer {
	return &globalizationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// HelmCharts returns a HelmChartInformer.
func (v *version) HelmCharts() HelmChartInformer {
	return &helmChartInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// HelmReleases returns a HelmReleaseInformer.
func (v *version) HelmReleases() HelmReleaseInformer {
	return &helmReleaseInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Localizations returns a LocalizationInformer.
func (v *version) Localizations() LocalizationInformer {
	return &localizationInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Manifests returns a ManifestInformer.
func (v *version) Manifests() ManifestInformer {
	return &manifestInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Subscriptions returns a SubscriptionInformer.
func (v *version) Subscriptions() SubscriptionInformer {
	return &subscriptionInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}