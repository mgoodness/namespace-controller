package kube

import (
	"path/filepath"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLoadNamespaces(t *testing.T) {
	k := &Kubernetes{}
	directory := "testdata"
	if err := k.LoadNamespaces(&directory); err != nil {
		t.Error(err)
	}

	if len(k.Namespaces) != 1 {
		t.Error("Should have one Namespace")
	}
}

func TestLoadResource(t *testing.T) {
	k := &Kubernetes{}
	desired := &schema.GroupVersionKind{
		Group:   "",
		Kind:    "Namespace",
		Version: "v1",
	}

	file := "test9.yaml"
	path := filepath.Join("testdata", file)
	if err := k.LoadResource(&path, desired); err == nil {
		t.Errorf("File %s should not exist", file)
	}

	file = "tiller-service.yaml"
	path = filepath.Join("testdata", file)
	if err := k.LoadResource(&path, desired); err == nil {
		t.Errorf("File %s should not be a valid Namespace manifest", file)
	}

	if err := k.LoadResource(&path, nil); err != nil {
		t.Error(err)
	}

	if len(k.services) != 1 {
		t.Error("Should have one service")
	}

	file = "junk.txt"
	path = filepath.Join("testdata", file)
	if err := k.LoadResource(&path, desired); err == nil {
		t.Errorf("File %s should fail to decode", file)
	}
}

func TestAddAnnotations(t *testing.T) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{"key1": "value1"},
			Name:        "prd1234",
		},
	}
	annotations := map[string]string{"key2": "value2"}

	AddAnnotations(namespace, annotations)
	if namespace.GetAnnotations()["key2"] != "value2" {
		t.Error("Annotation should be key: value")
	}
}

func TestUpdateCreateClusterNamespace(t *testing.T) {
	want := "prd1234"
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}
	clusterNamespace := &ClusterNamespace{Namespace: namespace}

	k8sClient := fake.NewSimpleClientset()
	k := &Kubernetes{Client: k8sClient}

	if _, err := k.Client.CoreV1().Namespaces().Create(namespace); err != nil {
		t.Error(err)
	}

	if err := k.UpdateCreateClusterNamespace(clusterNamespace); err != nil {
		t.Error(err)
	}

	want = "prd354"
	clusterNamespace = &ClusterNamespace{
		Namespace: &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: want,
			},
		},
	}
	if err := k.UpdateCreateClusterNamespace(clusterNamespace); err != nil {
		t.Error(err)
	}
}

func TestUpdateCreateDeployment(t *testing.T) {
	namespace := "prd1811"
	want := "test"
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	k8sClient := fake.NewSimpleClientset()
	k := &Kubernetes{Client: k8sClient}

	if _, err := k.Client.AppsV1().Deployments(namespace).Create(deployment); err != nil {
		t.Error(err)
	}

	if err := k.updateCreateDeployment(&namespace, deployment); err != nil {
		t.Error(err)
	}

	want = "test2"
	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	if err := k.updateCreateDeployment(&namespace, deployment); err != nil {
		t.Error(err)
	}
}

func TestUpdateCreateRoleBinding(t *testing.T) {
	namespace := "prd1811"
	want := "test"
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	k8sClient := fake.NewSimpleClientset()
	k := &Kubernetes{Client: k8sClient}

	if _, err := k.Client.RbacV1().RoleBindings(namespace).Create(roleBinding); err != nil {
		t.Error(err)
	}

	if err := k.updateCreateRoleBinding(&namespace, roleBinding); err != nil {
		t.Error(err)
	}

	want = "test2"
	roleBinding = &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	if err := k.updateCreateRoleBinding(&namespace, roleBinding); err != nil {
		t.Error(err)
	}
}

func TestUpdateCreateService(t *testing.T) {
	namespace := "prd1811"
	want := "test"
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	k8sClient := fake.NewSimpleClientset()
	k := &Kubernetes{Client: k8sClient}

	if _, err := k.Client.CoreV1().Services(namespace).Create(service); err != nil {
		t.Error(err)
	}

	if err := k.updateCreateService(&namespace, service); err != nil {
		t.Error(err)
	}

	want = "test2"
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	if err := k.updateCreateService(&namespace, service); err != nil {
		t.Error(err)
	}
}

func TestUpdateCreateServiceAccount(t *testing.T) {
	namespace := "prd1811"
	want := "test"
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	k8sClient := fake.NewSimpleClientset()
	k := &Kubernetes{Client: k8sClient}

	if _, err := k.Client.CoreV1().ServiceAccounts(namespace).Create(serviceAccount); err != nil {
		t.Error(err)
	}

	if err := k.updateCreateServiceAccount(&namespace, serviceAccount); err != nil {
		t.Error(err)
	}

	want = "test2"
	serviceAccount = &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	if err := k.updateCreateServiceAccount(&namespace, serviceAccount); err != nil {
		t.Error(err)
	}
}

// func TestDeleteNamespace(t *testing.T) {
// 	c := testclient.NewSimpleClientset()
// 	k := New(c)

// 	want := "prd1234"
// 	namespace := &corev1.Namespace{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: want,
// 		},
// 	}

// 	_, err := k.Client.CoreV1().Namespaces().Create(namespace)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if err := k.DeleteNamespace(want); err != nil {
// 		t.Error(err)
// 	}
// }

// func TestGetNamespaces(t *testing.T) {
// 	c := testclient.NewSimpleClientset()
// 	k := New(c)

// 	want := "prd1234"
// 	namespace := &corev1.Namespace{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: want,
// 		},
// 	}

// 	_, err := k.Client.CoreV1().Namespaces().Create(namespace)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	namespaces, err := k.GetNamespaces()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	got := namespaces[0].GetName()
// 	if got != want {
// 		t.Errorf("Wanted %s, got %s", want, got)
// 	}
// }
