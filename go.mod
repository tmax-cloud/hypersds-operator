module github.com/tmax-cloud/hypersds-operator

go 1.13

require (
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/golang/mock v1.3.1
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.4
)
