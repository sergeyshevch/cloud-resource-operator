/*
Copyright 2021 Sergey Shevchenko <sergeyshevchdevelop@gmail.com>.

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

package v1alpha1

import (
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	smithydocument "github.com/aws/smithy-go/document"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type noSmithyDocumentSerde = smithydocument.NoSerde

// Tag A tag that can be added to an ElastiCache cluster or replication group. Tags are
// composed of a Key/Value pair. You can use tags to categorize and track all your
// ElastiCache resources, with the exception of global replication group. When you
// add or remove tags on replication groups, those actions will be replicated to
// all nodes in the replication group. A tag with a null Value is permitted.
type Tag struct {

	// The key for the tag. May not be null.
	Key *string `json:"key"`

	// The tag's value. May be null.
	Value *string `json:"value"`

	noSmithyDocumentSerde
}

type ElasticCacheAwsConfig struct {

	// Specifies whether the nodes in this Memcached cluster are created in a single
	// Availability Zone or created across multiple Availability Zones in the cluster's
	// region. This parameter is only supported for Memcached clusters. If the AZMode
	// and PreferredAvailabilityZones are not specified, ElastiCache assumes single-az
	// mode.
	AZMode types.AZMode `json:"azMode,omitempty"`

	// Reserved parameter. The password used to access a password protected server.
	// Password constraints:
	//
	// * Must be only printable ASCII characters.
	//
	// * Must be at
	// least 16 characters and no more than 128 characters in length.
	//
	// * The only
	// permitted printable special characters are !, &, #, $, ^, <, >, and -. Other
	// printable special characters cannot be used in the AUTH token.
	//
	// For more
	// information, see AUTH password (http://redis.io/commands/AUTH) at
	// http://redis.io/commands/AUTH.
	AuthToken *string `json:"authToken,omitempty"`

	// Specifies the strategy to use to update the AUTH token. This parameter must be
	// specified with the auth-token parameter. Possible values:
	//
	// * Rotate
	//
	// * Set
	//
	// For
	// more information, see Authenticating Users with Redis AUTH
	// (http://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/auth.html)
	AuthTokenUpdateStrategy types.AuthTokenUpdateStrategyType `json:"authTokenUpdateStrategy,omitempty"`

	// The compute and memory capacity of the nodes in the node group (shard). The
	// following node types are supported by ElastiCache. Generally speaking, the
	// current generation types provide more memory and computational power at lower
	// cost when compared to their equivalent previous generation counterparts.
	//
	// *
	// General purpose:
	//
	// * Current generation: M6g node types (available only for Redis
	// engine version 5.0.6 onward and for Memcached engine version 1.5.16 onward).
	// cache.m6g.large, cache.m6g.xlarge, cache.m6g.2xlarge, cache.m6g.4xlarge,
	// cache.m6g.8xlarge, cache.m6g.12xlarge, cache.m6g.16xlarge For region
	// availability, see Supported Node Types
	// (https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/CacheNodes.SupportedTypes.html#CacheNodes.SupportedTypesByRegion)
	// M5 node types: cache.m5.large, cache.m5.xlarge, cache.m5.2xlarge,
	// cache.m5.4xlarge, cache.m5.12xlarge, cache.m5.24xlarge M4 node types:
	// cache.m4.large, cache.m4.xlarge, cache.m4.2xlarge, cache.m4.4xlarge,
	// cache.m4.10xlarge T3 node types: cache.t3.micro, cache.t3.small, cache.t3.medium
	// T2 node types: cache.t2.micro, cache.t2.small, cache.t2.medium
	//
	// * Previous
	// generation: (not recommended) T1 node types: cache.t1.micro M1 node types:
	// cache.m1.small, cache.m1.medium, cache.m1.large, cache.m1.xlarge M3 node types:
	// cache.m3.medium, cache.m3.large, cache.m3.xlarge, cache.m3.2xlarge
	//
	// * Compute
	// optimized:
	//
	// * Previous generation: (not recommended) C1 node types:
	// cache.c1.xlarge
	//
	// * Memory optimized:
	//
	// * Current generation: R6g node types
	// (available only for Redis engine version 5.0.6 onward and for Memcached engine
	// version 1.5.16 onward). cache.r6g.large, cache.r6g.xlarge, cache.r6g.2xlarge,
	// cache.r6g.4xlarge, cache.r6g.8xlarge, cache.r6g.12xlarge, cache.r6g.16xlarge For
	// region availability, see Supported Node Types
	// (https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/CacheNodes.SupportedTypes.html#CacheNodes.SupportedTypesByRegion)
	// R5 node types: cache.r5.large, cache.r5.xlarge, cache.r5.2xlarge,
	// cache.r5.4xlarge, cache.r5.12xlarge, cache.r5.24xlarge R4 node types:
	// cache.r4.large, cache.r4.xlarge, cache.r4.2xlarge, cache.r4.4xlarge,
	// cache.r4.8xlarge, cache.r4.16xlarge
	//
	// * Previous generation: (not recommended) M2
	// node types: cache.m2.xlarge, cache.m2.2xlarge, cache.m2.4xlarge R3 node types:
	// cache.r3.large, cache.r3.xlarge, cache.r3.2xlarge,
	//
	// cache.r3.4xlarge,
	// cache.r3.8xlarge
	//
	// Additional node type info
	//
	// * All current generation instance
	// types are created in Amazon VPC by default.
	//
	// * Redis append-only files (AOF) are
	// not supported for T1 or T2 instances.
	//
	// * Redis Multi-AZ with automatic failover
	// is not supported on T1 instances.
	//
	// * Redis configuration variables appendonly
	// and appendfsync are not supported on Redis version 2.8.22 and later.
	CacheNodeType *string `json:"cacheNodeType"`

	// The name of the parameter group to associate with this cluster. If this argument
	// is omitted, the default parameter group for the specified engine is used. You
	// cannot use any parameter group which has cluster-enabled='yes' when creating a
	// cluster.
	CacheParameterGroupName *string `json:"cacheParameterGroupName,omitempty"`

	// A list of security group names to associate with this cluster. Use this
	// parameter only when you are creating a cluster outside of an Amazon Virtual
	// Private Cloud (Amazon VPC).
	CacheSecurityGroupNames []string `json:"cacheSecurityGroupNames,omitempty"`

	// The name of the subnet group to be used for the cluster. Use this parameter only
	// when you are creating a cluster in an Amazon Virtual Private Cloud (Amazon VPC).
	// If you're going to launch your cluster in an Amazon VPC, you need to create a
	// subnet group before you start creating a cluster. For more information, see
	// Subnets and Subnet Groups
	// (https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/SubnetGroups.html).
	CacheSubnetGroupName *string `json:"cacheSubnetGroupName,omitempty"`

	// The name of the cache engine to be used for this cluster. Valid values for this
	// parameter are: memcached | redis
	Engine *string `json:"engine"`

	// The version number of the cache engine to be used for this cluster. To view the
	// supported cache engine versions, use the DescribeCacheEngineVersions operation.
	// Important: You can upgrade to a newer engine version (see Selecting a Cache
	// Engine and Version
	// (https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/SelectEngine.html#VersionManagement)),
	// but you cannot downgrade to an earlier engine version. If you want to use an
	// earlier engine version, you must delete the existing cluster or replication
	// group and create it anew with the earlier engine version.
	EngineVersion *string `json:"engineVersion"`

	// Specifies the destination, format and type of the logs.
	// TODO: Enable LogDeliveryConfigurations
	// LogDeliveryConfigurations []types.LogDeliveryConfigurationRequest `json:"logDeliveryConfigurations,omitempty"`

	// The Amazon Resource Name (ARN) of the Amazon Simple Notification Service (SNS)
	// topic to which notifications are sent. The Amazon SNS topic owner must be the
	// same as the cluster owner.
	NotificationTopicArn *string `json:"notificationTopicArn,omitempty"`

	// The initial number of cache nodes that the cluster has. For clusters running
	// Redis, this value must be 1. For clusters running Memcached, this value must be
	// between 1 and 40. If you need more than 40 nodes for your Memcached cluster,
	// please fill out the ElastiCache Limit Increase Request form at
	// http://aws.amazon.com/contact-us/elasticache-node-limit-request/
	// (http://aws.amazon.com/contact-us/elasticache-node-limit-request/).
	NumCacheNodes *int32 `json:"numCacheNodes"`

	// Specifies whether the nodes in the cluster are created in a single outpost or
	// across multiple outposts.
	OutpostMode types.OutpostMode `json:"outpostMode,omitempty"`

	// The port number on which each of the cache nodes accepts connections.
	Port *int32 `json:"port,omitempty"`

	// The EC2 Availability Zone in which the cluster is created. All nodes belonging
	// to this cluster are placed in the preferred Availability Zone. If you want to
	// create your nodes across multiple Availability Zones, use
	// PreferredAvailabilityZones. Default: System chosen Availability Zone.
	PreferredAvailabilityZone *string `json:"preferredAvailabilityZone,omitempty"`

	// A list of the Availability Zones in which cache nodes are created. The order of
	// the zones in the list is not important. This option is only supported on
	// Memcached. If you are creating your cluster in an Amazon VPC (recommended) you
	// can only locate nodes in Availability Zones that are associated with the subnets
	// in the selected subnet group. The number of Availability Zones listed must equal
	// the value of NumCacheNodes. If you want all the nodes in the same Availability
	// Zone, use PreferredAvailabilityZone instead, or repeat the Availability Zone
	// multiple times in the list. Default: System chosen Availability Zones.
	PreferredAvailabilityZones []string `json:"preferredAvailabilityZones,omitempty"`

	// Specifies the weekly time range during which maintenance on the cluster is
	// performed. It is specified as a range in the format ddd:hh24:mi-ddd:hh24:mi (24H
	// Clock UTC). The minimum maintenance window is a 60 minute period. Valid values
	// for ddd are:
	PreferredMaintenanceWindow *string `json:"preferredMaintenanceWindow,omitempty"`

	// The outpost ARN in which the cache cluster is created.
	PreferredOutpostArn *string `json:"preferredOutpostArn,omitempty"`

	// The outpost ARNs in which the cache cluster is created.
	PreferredOutpostArns []string `json:"preferredOutpostArns,omitempty"`

	// The ID of the replication group to which this cluster should belong. If this
	// parameter is specified, the cluster is added to the specified replication group
	// as a read replica; otherwise, the cluster is a standalone primary that is not
	// part of any replication group. If the specified replication group is Multi-AZ
	// enabled and the Availability Zone is not specified, the cluster is created in
	// Availability Zones that provide the best spread of read replicas across
	// Availability Zones. This parameter is only valid if the Engine parameter is
	// redis.
	ReplicationGroupId *string `json:"replicationGroupId,omitempty"`

	// One or more VPC security groups associated with the cluster. Use this parameter
	// only when you are creating a cluster in an Amazon Virtual Private Cloud (Amazon
	// VPC).
	SecurityGroupIds []string `json:"securityGroupIds,omitempty"`

	// A single-element string list containing an Amazon Resource Name (ARN) that
	// uniquely identifies a Redis RDB snapshot file stored in Amazon S3. The snapshot
	// file is used to populate the node group (shard). The Amazon S3 object name in
	// the ARN cannot contain any commas. This parameter is only valid if the Engine
	// parameter is redis. Example of an Amazon S3 ARN:
	// arn:aws:s3:::my_bucket/snapshot1.rdb
	SnapshotArns []string `json:"snapshotArns,omitempty"`

	// The name of a Redis snapshot from which to restore data into the new node group
	// (shard). The snapshot status changes to restoring while the new node group
	// (shard) is being created. This parameter is only valid if the Engine parameter
	// is redis.
	SnapshotName *string `json:"snapshotName,omitempty"`

	// The number of days for which ElastiCache retains automatic snapshots before
	// deleting them. For example, if you set SnapshotRetentionLimit to 5, a snapshot
	// taken today is retained for 5 days before being deleted. This parameter is only
	// valid if the Engine parameter is redis. Default: 0 (i.e., automatic backups are
	// disabled for this cache cluster).
	SnapshotRetentionLimit *int32 `json:"snapshotRetentionLimit,omitempty"`

	// The daily time range (in UTC) during which ElastiCache begins taking a daily
	// snapshot of your node group (shard). Example: 05:00-09:00 If you do not specify
	// this parameter, ElastiCache automatically chooses an appropriate time range.
	// This parameter is only valid if the Engine parameter is redis.
	SnapshotWindow *string `json:"snapshotWindow,omitempty"`

	// A list of tags to be added to this resource.
	Tags []Tag `json:"tags,omitempty"`
}

// ElasticCacheSpec defines the desired state of ElasticCache
type ElasticCacheSpec struct {
	AWSConfig *ElasticCacheAwsConfig `json:"awsConfig"`
}

// ElasticCacheStatus defines the observed state of ElasticCache
type ElasticCacheStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CacheClusterStatus *string `json:"cacheClusterStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ElasticCache is the Schema for the elasticcaches API
type ElasticCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticCacheSpec   `json:"spec,omitempty"`
	Status ElasticCacheStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ElasticCacheList contains a list of ElasticCache
type ElasticCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticCache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticCache{}, &ElasticCacheList{})
}
