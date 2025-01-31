// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package volume

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	apicommon "github.com/DataDog/datadog-operator/apis/datadoghq/common"
	apicommonv1 "github.com/DataDog/datadog-operator/apis/datadoghq/common/v1"
)

// GetVolumes creates a corev1.Volume and corev1.VolumeMount corresponding to a host path.
func GetVolumes(volumeName, hostPath, mountPath string, readOnly bool) (corev1.Volume, corev1.VolumeMount) {
	var volume corev1.Volume
	var volumeMount corev1.VolumeMount

	volume = corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: hostPath,
			},
		},
	}
	volumeMount = corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
		ReadOnly:  readOnly,
	}

	return volume, volumeMount
}

// GetVolumesEmptyDir creates a corev1.Volume (with an empty dir) and corev1.VolumeMount.
func GetVolumesEmptyDir(volumeName, mountPath string, readOnly bool) (corev1.Volume, corev1.VolumeMount) {
	var volume corev1.Volume
	var volumeMount corev1.VolumeMount

	volume = corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	volumeMount = corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
		ReadOnly:  readOnly,
	}

	return volume, volumeMount
}

// GetCustomConfigSpecVolumes use to generate the corev1.Volume and corev1.VolumeMount corresponding to a CustomConfig.
func GetCustomConfigSpecVolumes(customConfig *apicommonv1.CustomConfig, volumeName, defaultCMName, configFolder string) (corev1.Volume, corev1.VolumeMount) {
	var volume corev1.Volume
	var volumeMount corev1.VolumeMount
	if customConfig != nil {
		volume = GetVolumeFromCustomConfigSpec(
			customConfig,
			defaultCMName,
			volumeName,
		)
		// subpath only updated to Filekey if config uses configMap, default to ksmCoreCheckName for configData.
		volumeMount = GetVolumeMountFromCustomConfigSpec(
			customConfig,
			volumeName,
			fmt.Sprintf("%s%s/%s", apicommon.ConfigVolumePath, apicommon.ConfdVolumePath, configFolder),
			"",
		)
	} else {
		volume = corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: defaultCMName,
					},
				},
			},
		}
		volumeMount = corev1.VolumeMount{
			Name:      volumeName,
			MountPath: fmt.Sprintf("%s%s/%s", apicommon.ConfigVolumePath, apicommon.ConfdVolumePath, configFolder),
			ReadOnly:  true,
		}
	}
	return volume, volumeMount
}

// GetVolumeFromCustomConfigSpec return a corev1.Volume corresponding to a CustomConfig.
func GetVolumeFromCustomConfigSpec(cfcm *apicommonv1.CustomConfig, defaultConfigMapName, volumeName string) corev1.Volume {
	confdVolumeSource := *buildVolumeSourceFromCustomConfigSpec(cfcm, defaultConfigMapName)

	return corev1.Volume{
		Name:         volumeName,
		VolumeSource: confdVolumeSource,
	}
}

// GetVolumeMountFromCustomConfigSpec return a corev1.Volume corresponding to a CustomConfig.
func GetVolumeMountFromCustomConfigSpec(cfcm *apicommonv1.CustomConfig, volumeName, volumePath, defaultSubPath string) corev1.VolumeMount {
	subPath := defaultSubPath
	if cfcm.ConfigMap != nil && len(cfcm.ConfigMap.Items) > 0 {
		subPath = cfcm.ConfigMap.Items[0].Path
	}

	return corev1.VolumeMount{
		Name:      volumeName,
		MountPath: volumePath,
		SubPath:   subPath,
		ReadOnly:  true,
	}
}

func buildVolumeSourceFromCustomConfigSpec(configDir *apicommonv1.CustomConfig, defaultConfigMapName string) *corev1.VolumeSource {
	if configDir == nil {
		return nil
	}

	return buildVolumeSourceFromConfigMapConfig(configDir.ConfigMap, defaultConfigMapName)
}

// GetConfigMapVolumes use to generate the corev1.Volume and corev1.VolumeMount corresponding to a ConfigMapConfig.
func GetConfigMapVolumes(configMap *apicommonv1.ConfigMapConfig, defaultCMName, volumeName, volumePath string) (corev1.Volume, corev1.VolumeMount) {
	var volume corev1.Volume
	var volumeMount corev1.VolumeMount
	volume = GetVolumeFromConfigMapConfig(
		configMap,
		defaultCMName,
		volumeName,
	)

	volumeMount = GetVolumeMountFromConfigMapConfig(
		configMap,
		volumeName,
		volumePath,
		"",
	)
	return volume, volumeMount
}

// GetVolumeFromConfigMapConfig return a corev1.Volume corresponding to a ConfigMapConfig.
func GetVolumeFromConfigMapConfig(configMap *apicommonv1.ConfigMapConfig, defaultConfigMapName, volumeName string) corev1.Volume {
	confdVolumeSource := *buildVolumeSourceFromConfigMapConfig(configMap, defaultConfigMapName)

	return corev1.Volume{
		Name:         volumeName,
		VolumeSource: confdVolumeSource,
	}
}

// GetVolumeMountFromConfigMapConfig return a corev1.Volume corresponding to a ConfigMapConfig.
func GetVolumeMountFromConfigMapConfig(configMap *apicommonv1.ConfigMapConfig, volumeName, volumePath, defaultSubPath string) corev1.VolumeMount {
	subPath := defaultSubPath
	if configMap != nil && len(configMap.Items) > 0 {
		subPath = configMap.Items[0].Path
	}

	return corev1.VolumeMount{
		Name:      volumeName,
		MountPath: volumePath,
		SubPath:   subPath,
		ReadOnly:  true,
	}
}

func buildVolumeSourceFromConfigMapConfig(configMap *apicommonv1.ConfigMapConfig, defaultConfigMapName string) *corev1.VolumeSource {
	cmName := defaultConfigMapName
	if configMap != nil && len(configMap.Name) > 0 {
		cmName = configMap.Name
	}

	cmSource := &corev1.ConfigMapVolumeSource{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: cmName,
		},
	}

	return &corev1.VolumeSource{
		ConfigMap: cmSource,
	}
}
