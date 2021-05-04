module github.com/tmax-cloud/hypersds-operator

go 1.13

require (
	github.com/go-logr/logr v0.3.0
	github.com/golang/mock v1.4.1
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/tmax-cloud/hypersds-operator => ../hypersds-operator
