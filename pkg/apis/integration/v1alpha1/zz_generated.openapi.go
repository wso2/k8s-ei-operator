/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.Integration":       schema_pkg_apis_integration_v1alpha1_Integration(ref),
		"github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.IntegrationSpec":   schema_pkg_apis_integration_v1alpha1_IntegrationSpec(ref),
		"github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.IntegrationStatus": schema_pkg_apis_integration_v1alpha1_IntegrationStatus(ref),
	}
}

func schema_pkg_apis_integration_v1alpha1_Integration(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Integration is the Schema for the integrations API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.IntegrationSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.IntegrationStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.IntegrationSpec", "github.com/wso2/k8s-ei-operator/pkg/apis/integration/v1alpha1.IntegrationStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_integration_v1alpha1_IntegrationSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "IntegrationSpec defines the desired state of Integration",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_integration_v1alpha1_IntegrationStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "IntegrationStatus defines the observed state of Integration",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
