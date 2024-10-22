//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021 BlackRock, Inc.

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

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	common "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.BusConfig":           schema_pkg_apis_eventbus_v1alpha1_BusConfig(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate":   schema_pkg_apis_eventbus_v1alpha1_ContainerTemplate(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBus":            schema_pkg_apis_eventbus_v1alpha1_EventBus(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusList":        schema_pkg_apis_eventbus_v1alpha1_EventBusList(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusSpec":        schema_pkg_apis_eventbus_v1alpha1_EventBusSpec(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusStatus":      schema_pkg_apis_eventbus_v1alpha1_EventBusStatus(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamBus":        schema_pkg_apis_eventbus_v1alpha1_JetStreamBus(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamConfig":     schema_pkg_apis_eventbus_v1alpha1_JetStreamConfig(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaBus":            schema_pkg_apis_eventbus_v1alpha1_KafkaBus(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaConsumerGroup":  schema_pkg_apis_eventbus_v1alpha1_KafkaConsumerGroup(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSBus":             schema_pkg_apis_eventbus_v1alpha1_NATSBus(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSConfig":          schema_pkg_apis_eventbus_v1alpha1_NATSConfig(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NativeStrategy":      schema_pkg_apis_eventbus_v1alpha1_NativeStrategy(ref),
		"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.PersistenceStrategy": schema_pkg_apis_eventbus_v1alpha1_PersistenceStrategy(ref),
	}
}

func schema_pkg_apis_eventbus_v1alpha1_BusConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "BusConfig has the finalized configuration for EventBus",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"nats": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSConfig"),
						},
					},
					"jetstream": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamConfig"),
						},
					},
					"kafka": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaBus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamConfig", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaBus", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSConfig"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_ContainerTemplate(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ContainerTemplate defines customized spec for a container",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"resources": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/api/core/v1.ResourceRequirements"),
						},
					},
					"imagePullPolicy": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"securityContext": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/api/core/v1.SecurityContext"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.ResourceRequirements", "k8s.io/api/core/v1.SecurityContext"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_EventBus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EventBus is the definition of a eventbus resource",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusStatus"),
						},
					},
				},
				Required: []string{"metadata", "spec"},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusSpec", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBusStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_EventBusList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EventBusList is the list of eventbus resources",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBus"),
									},
								},
							},
						},
					},
				},
				Required: []string{"metadata", "items"},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.EventBus", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_EventBusSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EventBusSpec refers to specification of eventbus resource",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"nats": {
						SchemaProps: spec.SchemaProps{
							Description: "NATS eventbus",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSBus"),
						},
					},
					"jetstream": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamBus"),
						},
					},
					"kafka": {
						SchemaProps: spec.SchemaProps{
							Description: "Kafka eventbus",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaBus"),
						},
					},
					"jetstreamExotic": {
						SchemaProps: spec.SchemaProps{
							Description: "Exotic JetStream",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamConfig"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamBus", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.JetStreamConfig", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaBus", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSBus"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_EventBusStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EventBusStatus holds the status of the eventbus resource",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "type",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions are the latest available observations of a resource's current state.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/argoproj/argo-events/pkg/apis/common.Condition"),
									},
								},
							},
						},
					},
					"config": {
						SchemaProps: spec.SchemaProps{
							Description: "Config holds the fininalized configuration of EventBus",
							Default:     map[string]interface{}{},
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.BusConfig"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/common.Condition", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.BusConfig"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_JetStreamBus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JetStreamBus holds the JetStream EventBus information",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "JetStream version, such as \"2.7.3\"",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"replicas": {
						SchemaProps: spec.SchemaProps{
							Description: "JetStream StatefulSet size",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"containerTemplate": {
						SchemaProps: spec.SchemaProps{
							Description: "ContainerTemplate contains customized spec for Nats JetStream container",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate"),
						},
					},
					"reloaderContainerTemplate": {
						SchemaProps: spec.SchemaProps{
							Description: "ReloaderContainerTemplate contains customized spec for config reloader container",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate"),
						},
					},
					"metricsContainerTemplate": {
						SchemaProps: spec.SchemaProps{
							Description: "MetricsContainerTemplate contains customized spec for metrics container",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate"),
						},
					},
					"persistence": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.PersistenceStrategy"),
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Metadata sets the pods's metadata, i.e. annotations and labels",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/common.Metadata"),
						},
					},
					"nodeSelector": {
						SchemaProps: spec.SchemaProps{
							Description: "NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"tolerations": {
						SchemaProps: spec.SchemaProps{
							Description: "If specified, the pod's tolerations.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/api/core/v1.Toleration"),
									},
								},
							},
						},
					},
					"securityContext": {
						SchemaProps: spec.SchemaProps{
							Description: "SecurityContext holds pod-level security attributes and common container settings. Optional: Defaults to empty.  See type description for default values of each field.",
							Ref:         ref("k8s.io/api/core/v1.PodSecurityContext"),
						},
					},
					"imagePullSecrets": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "name",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/api/core/v1.LocalObjectReference"),
									},
								},
							},
						},
					},
					"priorityClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "If specified, indicates the Redis pod's priority. \"system-node-critical\" and \"system-cluster-critical\" are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name. If not specified, the pod priority will be default or zero if there is no default. More info: https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"priority": {
						SchemaProps: spec.SchemaProps{
							Description: "The priority value. Various system components use this field to find the priority of the Redis pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority. More info: https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"affinity": {
						SchemaProps: spec.SchemaProps{
							Description: "The pod's scheduling constraints More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/",
							Ref:         ref("k8s.io/api/core/v1.Affinity"),
						},
					},
					"serviceAccountName": {
						SchemaProps: spec.SchemaProps{
							Description: "ServiceAccountName to apply to the StatefulSet",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"settings": {
						SchemaProps: spec.SchemaProps{
							Description: "JetStream configuration, if not specified, global settings in controller-config will be used. See https://docs.nats.io/running-a-nats-service/configuration#jetstream. Only configure \"max_memory_store\" or \"max_file_store\", do not set \"store_dir\" as it has been hardcoded.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"startArgs": {
						SchemaProps: spec.SchemaProps{
							Description: "Optional arguments to start nats-server. For example, \"-D\" to enable debugging output, \"-DV\" to enable debugging and tracing. Check https://docs.nats.io/ for all the available arguments.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"streamConfig": {
						SchemaProps: spec.SchemaProps{
							Description: "Optional configuration for the streams to be created in this JetStream service, if specified, it will be merged with the default configuration in controller-config. It accepts a YAML format configuration, available fields include, \"maxBytes\", \"maxMsgs\", \"maxAge\" (e.g. 72h), \"replicas\" (1, 3, 5), \"duplicates\" (e.g. 5m).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"maxPayload": {
						SchemaProps: spec.SchemaProps{
							Description: "Maximum number of bytes in a message payload, 0 means unlimited. Defaults to 1MB",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/common.Metadata", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.PersistenceStrategy", "k8s.io/api/core/v1.Affinity", "k8s.io/api/core/v1.LocalObjectReference", "k8s.io/api/core/v1.PodSecurityContext", "k8s.io/api/core/v1.Toleration"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_JetStreamConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"url": {
						SchemaProps: spec.SchemaProps{
							Description: "JetStream (Nats) URL",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"accessSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "Secret for auth",
							Ref:         ref("k8s.io/api/core/v1.SecretKeySelector"),
						},
					},
					"streamConfig": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.SecretKeySelector"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_KafkaBus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KafkaBus holds the KafkaBus EventBus information",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"url": {
						SchemaProps: spec.SchemaProps{
							Description: "URL to kafka cluster, multiple URLs separated by comma",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"topic": {
						SchemaProps: spec.SchemaProps{
							Description: "Topic name, defaults to {namespace_name}-{eventbus_name}",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "Kafka version, sarama defaults to the oldest supported stable version",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"tls": {
						SchemaProps: spec.SchemaProps{
							Description: "TLS configuration for the kafka client.",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/common.TLSConfig"),
						},
					},
					"sasl": {
						SchemaProps: spec.SchemaProps{
							Description: "SASL configuration for the kafka client",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/common.SASLConfig"),
						},
					},
					"consumerGroup": {
						SchemaProps: spec.SchemaProps{
							Description: "Consumer group for kafka client",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaConsumerGroup"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/common.SASLConfig", "github.com/argoproj/argo-events/pkg/apis/common.TLSConfig", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.KafkaConsumerGroup"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_KafkaConsumerGroup(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"groupName": {
						SchemaProps: spec.SchemaProps{
							Description: "Consumer group name, defaults to {namespace_name}-{sensor_name}",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"rebalanceStrategy": {
						SchemaProps: spec.SchemaProps{
							Description: "Rebalance strategy can be one of: sticky, roundrobin, range. Range is the default.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"startOldest": {
						SchemaProps: spec.SchemaProps{
							Description: "When starting up a new group do we want to start from the oldest event (true) or the newest event (false), defaults to false",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_NATSBus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NATSBus holds the NATS eventbus information",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"native": {
						SchemaProps: spec.SchemaProps{
							Description: "Native means to bring up a native NATS service",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NativeStrategy"),
						},
					},
					"exotic": {
						SchemaProps: spec.SchemaProps{
							Description: "Exotic holds an exotic NATS config",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSConfig"),
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Description: "StatefulSet metadata, we actually uses only annotation from it for now",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/common.Metadata"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/common.Metadata", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NATSConfig", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.NativeStrategy"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_NATSConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NATSConfig holds the config of NATS",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"url": {
						SchemaProps: spec.SchemaProps{
							Description: "NATS streaming url",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"clusterID": {
						SchemaProps: spec.SchemaProps{
							Description: "Cluster ID for nats streaming",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"auth": {
						SchemaProps: spec.SchemaProps{
							Description: "Auth strategy, default to AuthStrategyNone",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"accessSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "Secret for auth",
							Ref:         ref("k8s.io/api/core/v1.SecretKeySelector"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.SecretKeySelector"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_NativeStrategy(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NativeStrategy indicates to install a native NATS service",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"replicas": {
						SchemaProps: spec.SchemaProps{
							Description: "Size is the NATS StatefulSet size",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"auth": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"persistence": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.PersistenceStrategy"),
						},
					},
					"containerTemplate": {
						SchemaProps: spec.SchemaProps{
							Description: "ContainerTemplate contains customized spec for NATS container",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate"),
						},
					},
					"metricsContainerTemplate": {
						SchemaProps: spec.SchemaProps{
							Description: "MetricsContainerTemplate contains customized spec for metrics container",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate"),
						},
					},
					"nodeSelector": {
						SchemaProps: spec.SchemaProps{
							Description: "NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"tolerations": {
						SchemaProps: spec.SchemaProps{
							Description: "If specified, the pod's tolerations.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/api/core/v1.Toleration"),
									},
								},
							},
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Metadata sets the pods's metadata, i.e. annotations and labels",
							Ref:         ref("github.com/argoproj/argo-events/pkg/apis/common.Metadata"),
						},
					},
					"securityContext": {
						SchemaProps: spec.SchemaProps{
							Description: "SecurityContext holds pod-level security attributes and common container settings. Optional: Defaults to empty.  See type description for default values of each field.",
							Ref:         ref("k8s.io/api/core/v1.PodSecurityContext"),
						},
					},
					"maxAge": {
						SchemaProps: spec.SchemaProps{
							Description: "Max Age of existing messages, i.e. \"72h\", “4h35m”",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"imagePullSecrets": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "name",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/api/core/v1.LocalObjectReference"),
									},
								},
							},
						},
					},
					"serviceAccountName": {
						SchemaProps: spec.SchemaProps{
							Description: "ServiceAccountName to apply to NATS StatefulSet",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"priorityClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "If specified, indicates the EventSource pod's priority. \"system-node-critical\" and \"system-cluster-critical\" are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name. If not specified, the pod priority will be default or zero if there is no default. More info: https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"priority": {
						SchemaProps: spec.SchemaProps{
							Description: "The priority value. Various system components use this field to find the priority of the EventSource pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority. More info: https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"affinity": {
						SchemaProps: spec.SchemaProps{
							Description: "The pod's scheduling constraints More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/",
							Ref:         ref("k8s.io/api/core/v1.Affinity"),
						},
					},
					"maxMsgs": {
						SchemaProps: spec.SchemaProps{
							Description: "Maximum number of messages per channel, 0 means unlimited. Defaults to 1000000",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"maxBytes": {
						SchemaProps: spec.SchemaProps{
							Description: "Total size of messages per channel, 0 means unlimited. Defaults to 1GB",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"maxSubs": {
						SchemaProps: spec.SchemaProps{
							Description: "Maximum number of subscriptions per channel, 0 means unlimited. Defaults to 1000",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"maxPayload": {
						SchemaProps: spec.SchemaProps{
							Description: "Maximum number of bytes in a message payload, 0 means unlimited. Defaults to 1MB",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"raftHeartbeatTimeout": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the time in follower state without a leader before attempting an election, i.e. \"72h\", “4h35m”. Defaults to 2s",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"raftElectionTimeout": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the time in candidate state without a leader before attempting an election, i.e. \"72h\", “4h35m”. Defaults to 2s",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"raftLeaseTimeout": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies how long a leader waits without being able to contact a quorum of nodes before stepping down as leader, i.e. \"72h\", “4h35m”. Defaults to 1s",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"raftCommitTimeout": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the time without an Apply() operation before sending an heartbeat to ensure timely commit, i.e. \"72h\", “4h35m”. Defaults to 100ms",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/common.Metadata", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.ContainerTemplate", "github.com/argoproj/argo-events/pkg/apis/eventbus/v1alpha1.PersistenceStrategy", "k8s.io/api/core/v1.Affinity", "k8s.io/api/core/v1.LocalObjectReference", "k8s.io/api/core/v1.PodSecurityContext", "k8s.io/api/core/v1.Toleration"},
	}
}

func schema_pkg_apis_eventbus_v1alpha1_PersistenceStrategy(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PersistenceStrategy defines the strategy of persistence",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"storageClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"accessMode": {
						SchemaProps: spec.SchemaProps{
							Description: "Available access modes such as ReadWriteOnce, ReadWriteMany https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"volumeSize": {
						SchemaProps: spec.SchemaProps{
							Description: "Volume size, e.g. 10Gi",
							Ref:         ref("k8s.io/apimachinery/pkg/api/resource.Quantity"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/api/resource.Quantity"},
	}
}
