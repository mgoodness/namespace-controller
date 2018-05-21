package kube

import (
	"errors"
	"io/ioutil"

	apiv1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	// Required for OIDC auth support
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"k8s.io/client-go/util/retry"
)

// Kubernetes struct contains a Kubernetes client
type Kubernetes struct {
	Client kubernetes.Interface
}

// New creates a new Kubernetes struct
func New(clientset kubernetes.Interface) *Kubernetes {
	return &Kubernetes{
		Client: clientset,
	}
}

// LoadNamespace creates a Namespace resource from a manifest on disk
func LoadNamespace(filePath string) (*apiv1.Namespace, error) {
	object, err := loadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	namespace, ok := object.(*apiv1.Namespace)
	if !ok {
		return nil, errors.New("Not a valid Namespace manifest")
	}

	return namespace, nil
}

// AddNamespaceAnnotations adds annotations to a Namespace object
func AddNamespaceAnnotations(namespace *apiv1.Namespace, annotations map[string]string) {
	for key, value := range namespace.GetAnnotations() {
		annotations[key] = value
	}

	namespace.SetAnnotations(annotations)
}

// UpdateCreateNamespace updates or creates a Namespace in the cluster
func (k *Kubernetes) UpdateCreateNamespace(namespace *apiv1.Namespace) error {
	namespacesClient := k.Client.CoreV1().Namespaces()

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := namespacesClient.Get(namespace.GetName(), metav1.GetOptions{})
		if err != nil {
			switch {
			case apiErrors.IsNotFound(err):
				_, err = namespacesClient.Create(namespace)
				if err != nil {
					return err
				}
			default:
				return err
			}
		} else {
			if _, err := namespacesClient.Update(namespace); err != nil {
				return err
			}
		}

		return err
	})
	if err != nil {
		return err
	}

	return err
}

// DeleteNamespace deletes a Namespace from the cluster
func (k *Kubernetes) DeleteNamespace(name string) error {
	namespacesClient := k.Client.CoreV1().Namespaces()

	deletePolicy := metav1.DeletePropagationForeground

	return namespacesClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// GetNamespaces gets a list of current Namespaces from the cluster
func (k *Kubernetes) GetNamespaces() ([]apiv1.Namespace, error) {
	namespacesClient := k.Client.CoreV1().Namespaces()

	namespaceList, err := namespacesClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaceList.Items, nil
}

func loadFromFile(filePath string) (runtime.Object, error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	object, _, err := decode(contents, nil, nil)
	if err != nil {
		return nil, err
	}

	return object, err
}
