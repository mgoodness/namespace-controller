package kube

import (
	"path/filepath"
	"testing"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestLoadNamespaces(t *testing.T) {
	file := "prd1811.yaml"
	path := filepath.Join("testdata", file)

	n, err := LoadNamespace(path)
	if err != nil {
		t.Error(err)
	}

	want := "prd1811"
	got := n.GetName()
	if got != want {
		t.Errorf("Wanted %s, got %s", want, got)
	}

	file = "test9.yaml"
	path = filepath.Join("testdata", file)
	n, err = LoadNamespace(path)
	if err == nil {
		t.Errorf("File %s should not exist", file)
	}

	file = "tiller-service.yaml"
	path = filepath.Join("testdata", file)
	n, err = LoadNamespace(path)
	if err == nil {
		t.Errorf("File %s should not be a valid Namespace manifest", file)
	}

	file = "junk.txt"
	path = filepath.Join("testdata", file)
	n, err = LoadNamespace(path)
	if err == nil {
		t.Errorf("File %s should fail to decode", file)
	}
}

func TestAddNamespaceAnnotations(t *testing.T) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{"key1": "value1"},
			Name:        "prd1234",
		},
	}

	annotations := map[string]string{"key2": "value2"}

	AddNamespaceAnnotations(namespace, annotations)
	if namespace.GetAnnotations()["key2"] != "value2" {
		t.Error("Annotation should be key: value")
	}
}

func TestUpdateCreateNamespace(t *testing.T) {
	c := testclient.NewSimpleClientset()
	k := New(c)

	want := "prd1234"
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	_, err := k.Client.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		t.Error(err)
	}

	if err := k.UpdateCreateNamespace(namespace); err != nil {
		t.Error(err)
	}

	want = "prd354"
	namespace = &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	if err := k.UpdateCreateNamespace(namespace); err != nil {
		t.Error(err)
	}
}

func TestDeleteNamespace(t *testing.T) {
	c := testclient.NewSimpleClientset()
	k := New(c)

	want := "prd1234"
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	_, err := k.Client.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		t.Error(err)
	}

	if err := k.DeleteNamespace(want); err != nil {
		t.Error(err)
	}
}

func TestGetNamespaces(t *testing.T) {
	c := testclient.NewSimpleClientset()
	k := New(c)

	want := "prd1234"
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: want,
		},
	}

	_, err := k.Client.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		t.Error(err)
	}

	namespaces, err := k.GetNamespaces()
	if err != nil {
		t.Error(err)
	}

	got := namespaces[0].GetName()
	if got != want {
		t.Errorf("Wanted %s, got %s", want, got)
	}
}
