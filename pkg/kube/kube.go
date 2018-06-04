package kube

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	// Required for OIDC auth support
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type ClusterNamespace struct {
	Namespace       *corev1.Namespace
	Deployments     []*appsv1.Deployment
	RoleBindings    []*rbacv1.RoleBinding
	ServiceAccounts []*corev1.ServiceAccount
	Services        []*corev1.Service
}

type Kubernetes struct {
	Client          kubernetes.Interface
	Namespaces      []*corev1.Namespace
	deployments     []*appsv1.Deployment
	roleBindings    []*rbacv1.RoleBinding
	serviceAccounts []*corev1.ServiceAccount
	services        []*corev1.Service
}

var kubeLogger zerolog.Logger

func init() {
	kubeLogger = log.With().Str("component", "kube").Logger()
}

func (k *Kubernetes) addResource(resource runtime.Object) error {
	switch r := resource.(type) {
	case *appsv1.Deployment:
		k.deployments = append(k.deployments, r)
	case *corev1.Namespace:
		k.Namespaces = append(k.Namespaces, r)
	case *rbacv1.RoleBinding:
		k.roleBindings = append(k.roleBindings, r)
	case *corev1.ServiceAccount:
		k.serviceAccounts = append(k.serviceAccounts, r)
	case *corev1.Service:
		k.services = append(k.services, r)
	default:
		return fmt.Errorf("Kind %s is not supported", r.GetObjectKind().GroupVersionKind().Kind)
	}

	return nil
}

// LoadResource loads a Kubernetes resource from a manifest on disk
func (k *Kubernetes) LoadResource(filePath *string, desired *schema.GroupVersionKind) error {
	contents, err := ioutil.ReadFile(*filePath)
	if err != nil {
		return err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	resource, gvk, err := decode(contents, desired, nil)
	if err != nil {
		return err
	}

	if desired != nil && *gvk != *desired {
		return &KindError{kind: desired.Kind}
	}

	k.addResource(resource)
	return nil
}

// LoadNamespaces loads Kubernetes Namespace resources from a directory of manifests on disk
func (k *Kubernetes) LoadNamespaces(directory *string) (err error) {
	absPath, _ := filepath.Abs(*directory)
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return
	}

	gvk := &schema.GroupVersionKind{
		Group:   "",
		Kind:    "Namespace",
		Version: "v1",
	}

	for _, file := range files {
		f := filepath.Join(absPath, file.Name())
		if filepath.Ext(f) != ".yaml" {
			kubeLogger.Info().Msgf("Skipping file %s", f)
			continue
		}

		if err = k.LoadResource(&f, gvk); err != nil {
			if _, ok := err.(*KindError); !ok {
				return
			}

			kubeLogger.Info().Err(err).Msgf("Skipping file %s", f)
			continue
		}
	}

	return nil
}

// AddAnnotations adds annotations to a Kubernetes resource from a map
func AddAnnotations(resource metav1.Object, annotations map[string]string) {
	for key, value := range resource.GetAnnotations() {
		annotations[key] = value
	}

	resource.SetAnnotations(annotations)
}

func (k *Kubernetes) UpdateCreateClusterNamespace(clusterNamespace *ClusterNamespace) error {
	namespaceName := clusterNamespace.Namespace.GetName()

	if err := k.updateCreateNamespace(clusterNamespace.Namespace); err != nil {
		return err
	}

	for _, deployment := range clusterNamespace.Deployments {
		if err := k.updateCreateDeployment(&namespaceName, deployment); err != nil {
			kubeLogger.Error().Err(err).Msgf("Unable to create/update Deployment %s", deployment.GetName())
		}
	}

	for _, roleBinding := range clusterNamespace.RoleBindings {
		if err := k.updateCreateRoleBinding(&namespaceName, roleBinding); err != nil {
			kubeLogger.Error().Err(err).Msgf("Unable to create/update RoleBinding %s", roleBinding.GetName())
		}
	}

	for _, service := range clusterNamespace.Services {
		if err := k.updateCreateService(&namespaceName, service); err != nil {
			kubeLogger.Error().Err(err).Msgf("Unable to create/update Service %s", service.GetName())
		}
	}

	for _, serviceAccount := range clusterNamespace.ServiceAccounts {
		if err := k.updateCreateServiceAccount(&namespaceName, serviceAccount); err != nil {
			kubeLogger.Error().Err(err).Msgf("Unable to create/update ServiceAccount %s", serviceAccount.GetName())
		}
	}

	return nil
}

// DeleteNamespace deletes a Namespace from the cluster
// func (k *Kubernetes) DeleteNamespace(name string) error {
// 	namespacesClient := k.client.CoreV1().Namespaces()

// 	deletePolicy := metav1.DeletePropagationForeground

// 	return namespacesClient.Delete(name, &metav1.DeleteOptions{
// 		PropagationPolicy: &deletePolicy,
// 	})
// }

// GetNamespaces gets a list of current Namespaces from the cluster
// func (k *Kubernetes) GetNamespaces() ([]corev1.Namespace, error) {
// 	namespacesClient := k.client.CoreV1().Namespaces()

// 	namespaceList, err := namespacesClient.List(metav1.ListOptions{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return namespaceList.Items, nil
// }
