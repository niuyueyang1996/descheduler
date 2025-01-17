/*
Copyright 2017 The Kubernetes Authors.

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

package client

import (
	"fmt"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"

	clientset "k8s.io/client-go/kubernetes"
	componentbaseconfig "k8s.io/component-base/config"

	// Ensure to load all auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateMetricsClient(clientConnection componentbaseconfig.ClientConnectionConfiguration) (*metricsclientset.Clientset, error) {
	var cfg *rest.Config
	if len(clientConnection.Kubeconfig) != 0 {
		master, err := GetMasterFromKubeconfig(clientConnection.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig file: %v ", err)
		}

		cfg, err = clientcmd.BuildConfigFromFlags(master, clientConnection.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("unable to build config: %v", err)
		}

	} else {
		var err error
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to build in cluster config: %v", err)
		}
	}

	cfg.Burst = int(clientConnection.Burst)
	cfg.QPS = clientConnection.QPS
	return metricsclientset.NewForConfig(cfg)
}
func CreateClient(clientConnection componentbaseconfig.ClientConnectionConfiguration, userAgt string) (clientset.Interface, error) {
	var cfg *rest.Config
	if len(clientConnection.Kubeconfig) != 0 {
		master, err := GetMasterFromKubeconfig(clientConnection.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig file: %v ", err)
		}

		cfg, err = clientcmd.BuildConfigFromFlags(master, clientConnection.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("unable to build config: %v", err)
		}

	} else {
		var err error
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to build in cluster config: %v", err)
		}
	}

	cfg.Burst = int(clientConnection.Burst)
	cfg.QPS = clientConnection.QPS

	if len(userAgt) != 0 {
		cfg = rest.AddUserAgent(cfg, userAgt)
	}

	return clientset.NewForConfig(cfg)
}

func GetMasterFromKubeconfig(filename string) (string, error) {
	config, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		return "", err
	}

	context, ok := config.Contexts[config.CurrentContext]
	if !ok {
		return "", fmt.Errorf("failed to get master address from kubeconfig")
	}

	if val, ok := config.Clusters[context.Cluster]; ok {
		return val.Server, nil
	}
	return "", fmt.Errorf("failed to get master address from kubeconfig")
}
