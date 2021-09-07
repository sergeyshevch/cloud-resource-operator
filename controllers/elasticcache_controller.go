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

package controllers

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	awsv1alpha1 "github.com/sergeyshevch/cloud-resource-operator/api/v1alpha1"
)

var awsResource = schema.GroupResource{Group: "aws.sergeyshevch.dev", Resource: "AwsResource"}
var elasticCacheFinalizer = "aws.serveyshevch.dev/finalizer"
var lastAppliedSpecAnnotation = "aws.sergeyshevch.dev/last-applied"

// ElasticCacheReconciler reconciles a ElasticCache object
type ElasticCacheReconciler struct {
	client.Client
	AwsConfig aws.Config
	Scheme    *runtime.Scheme
}

//+kubebuilder:rbac:groups=aws.sergeyshevch.dev,resources=elasticcaches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.sergeyshevch.dev,resources=elasticcaches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.sergeyshevch.dev,resources=elasticcaches/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ElasticCache object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *ElasticCacheReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	instance := &awsv1alpha1.ElasticCache{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	awsClient := elasticache.NewFromConfig(r.AwsConfig)

	// Process elasticCache cluster
	cacheCluster, err := r.getElasticCacheCluster(awsClient, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			cacheCluster, err = r.createElasticCacheCluster(awsClient, instance)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Update cluster status
			err = r.updateClusterStatus(cacheCluster, instance)
			if err != nil {
				return ctrl.Result{}, err
			}

			err = r.SetLastAppliedAnnotation(instance)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Cluster setup time
			return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
		}
		return ctrl.Result{}, err
	} else {
		needPatch, err := isPatchNeeded(instance)
		if err != nil {
			return ctrl.Result{}, nil
		}
		if needPatch {
			cacheCluster, err = r.patchElasticCacheCluster(awsClient, instance)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Update cluster status
			err = r.updateClusterStatus(cacheCluster, instance)
			if err != nil {
				return ctrl.Result{}, err
			}

			err = r.SetLastAppliedAnnotation(instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Update cluster status
	err = r.updateClusterStatus(cacheCluster, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	isElasticCacheMarkedToDeletion := instance.GetDeletionTimestamp() != nil
	if isElasticCacheMarkedToDeletion {
		if controllerutil.ContainsFinalizer(instance, elasticCacheFinalizer) {
			err = r.deleteElasticCacheCluster(awsClient, instance)
			if err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(instance, elasticCacheFinalizer)
			err = r.Update(ctx, instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if !controllerutil.ContainsFinalizer(instance, elasticCacheFinalizer) {
		controllerutil.AddFinalizer(instance, elasticCacheFinalizer)
		err = r.Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

func (r *ElasticCacheReconciler) updateClusterStatus(cluster *types.CacheCluster, instance *awsv1alpha1.ElasticCache) error {
	if instance.Status.CacheClusterStatus != cluster.CacheClusterStatus {
		instance.Status.CacheClusterStatus = cluster.CacheClusterStatus
		err := r.Status().Update(context.TODO(), instance)
		if err != nil {
			return err
		}
	}
	return nil
}

func isPatchNeeded(cr *awsv1alpha1.ElasticCache) (bool, error) {
	marshaled, err  := json.Marshal(cr.Spec)
	if err != nil {
		return false, err
	}

	original := cr.GetAnnotations()[lastAppliedSpecAnnotation]
	current := base64.StdEncoding.EncodeToString(marshaled)

	return original != current, nil
}

func (r *ElasticCacheReconciler) patchElasticCacheCluster(awsClient *elasticache.Client, cr *awsv1alpha1.ElasticCache) (*types.CacheCluster, error) {
	params := &elasticache.ModifyCacheClusterInput{
		CacheClusterId:             &cr.Name,
		AZMode:                     cr.Spec.AWSConfig.AZMode,
		ApplyImmediately:           true,
		AuthToken:                  cr.Spec.AWSConfig.AuthToken,
		AuthTokenUpdateStrategy:    cr.Spec.AWSConfig.AuthTokenUpdateStrategy,
		CacheNodeType:              cr.Spec.AWSConfig.CacheNodeType,
		CacheParameterGroupName:    cr.Spec.AWSConfig.CacheParameterGroupName,
		CacheSecurityGroupNames:    cr.Spec.AWSConfig.CacheSecurityGroupNames,
		EngineVersion:              cr.Spec.AWSConfig.EngineVersion,
		//LogDeliveryConfigurations:  cr.Spec.AWSConfig.LogDeliveryConfigurations,
		NotificationTopicArn:       cr.Spec.AWSConfig.NotificationTopicArn,
		NumCacheNodes:              cr.Spec.AWSConfig.NumCacheNodes,
		PreferredMaintenanceWindow: cr.Spec.AWSConfig.PreferredMaintenanceWindow,
		SecurityGroupIds:           cr.Spec.AWSConfig.SecurityGroupIds,
		SnapshotRetentionLimit:     cr.Spec.AWSConfig.SnapshotRetentionLimit,
		SnapshotWindow:             cr.Spec.AWSConfig.SnapshotWindow,
	}

	output, err := awsClient.ModifyCacheCluster(context.TODO(), params)
	if err != nil {
		return &types.CacheCluster{}, err
	}
	return output.CacheCluster, nil
}

func (r *ElasticCacheReconciler) createElasticCacheCluster(awsClient *elasticache.Client, cr *awsv1alpha1.ElasticCache) (*types.CacheCluster, error) {
	params := &elasticache.CreateCacheClusterInput{
		CacheClusterId:             &cr.Name,
		AZMode:                     cr.Spec.AWSConfig.AZMode,
		AuthToken:                  cr.Spec.AWSConfig.AuthToken,
		CacheNodeType:              cr.Spec.AWSConfig.CacheNodeType,
		CacheParameterGroupName:    cr.Spec.AWSConfig.CacheParameterGroupName,
		CacheSecurityGroupNames:    cr.Spec.AWSConfig.CacheSecurityGroupNames,
		CacheSubnetGroupName:       cr.Spec.AWSConfig.CacheSubnetGroupName,
		Engine:                     cr.Spec.AWSConfig.Engine,
		EngineVersion:              cr.Spec.AWSConfig.EngineVersion,
		//LogDeliveryConfigurations:  cr.Spec.AWSConfig.LogDeliveryConfigurations,
		NotificationTopicArn:       cr.Spec.AWSConfig.NotificationTopicArn,
		NumCacheNodes:              cr.Spec.AWSConfig.NumCacheNodes,
		OutpostMode:                cr.Spec.AWSConfig.OutpostMode,
		Port:                       cr.Spec.AWSConfig.Port,
		PreferredAvailabilityZone:  cr.Spec.AWSConfig.PreferredAvailabilityZone,
		PreferredAvailabilityZones: cr.Spec.AWSConfig.PreferredAvailabilityZones,
		PreferredMaintenanceWindow: cr.Spec.AWSConfig.PreferredMaintenanceWindow,
		PreferredOutpostArn:        cr.Spec.AWSConfig.PreferredOutpostArn,
		PreferredOutpostArns:       cr.Spec.AWSConfig.PreferredOutpostArns,
		ReplicationGroupId:         cr.Spec.AWSConfig.ReplicationGroupId,
		SecurityGroupIds:           cr.Spec.AWSConfig.SecurityGroupIds,
		SnapshotArns:               cr.Spec.AWSConfig.SnapshotArns,
		SnapshotName:               cr.Spec.AWSConfig.SnapshotName,
		SnapshotRetentionLimit:     cr.Spec.AWSConfig.SnapshotRetentionLimit,
		SnapshotWindow:             cr.Spec.AWSConfig.SnapshotWindow,
		Tags:                       cr.Spec.AWSConfig.Tags,
	}

	output, err := awsClient.CreateCacheCluster(context.TODO(), params)
	if err != nil {
		return &types.CacheCluster{}, err
	}
	return output.CacheCluster, nil

}

func (r *ElasticCacheReconciler) getElasticCacheCluster(awsClient *elasticache.Client, cr *awsv1alpha1.ElasticCache) (*types.CacheCluster, error) {
	params := &elasticache.DescribeCacheClustersInput{
		CacheClusterId: &cr.Name,
	}

	output, err := awsClient.DescribeCacheClusters(context.TODO(), params)
	if err != nil {
		return &types.CacheCluster{}, err
	}

	clusters := output.CacheClusters

	if len(clusters) == 1 {
		return &clusters[0], nil
	} else {
		return &types.CacheCluster{}, errors.NewNotFound(awsResource, "ElasticCacheCluster")
	}
}

func (r *ElasticCacheReconciler) deleteElasticCacheCluster(awsClient *elasticache.Client, cr *awsv1alpha1.ElasticCache) error {
	params := &elasticache.DeleteCacheClusterInput{
		CacheClusterId: &cr.Name,
	}

	_, err := awsClient.DeleteCacheCluster(context.TODO(), params)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ElasticCacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsv1alpha1.ElasticCache{}).
		Complete(r)
}

func (r *ElasticCacheReconciler) SetLastAppliedAnnotation(instance *awsv1alpha1.ElasticCache) error {
	marshaled, err  := json.Marshal(instance.Spec)
	if err != nil {
		return err
	}
	instance.SetAnnotations(map[string]string{lastAppliedSpecAnnotation: base64.StdEncoding.EncodeToString(marshaled)})
	err = r.Update(context.TODO(), instance)
	return err
}
