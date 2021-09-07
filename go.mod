module github.com/sergeyshevch/cloud-resource-operator

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.9.0
	github.com/aws/aws-sdk-go-v2/config v1.7.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.10.0
	github.com/aws/smithy-go v1.8.0 // indirect
	github.com/banzaicloud/k8s-objectmatcher v1.5.2
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
)
