package tiller

import (
	"fmt"
	"testing"
)

func TestGetVersion(t *testing.T) {
	annotation := "ticketmaster.com/tiller"
	defaultVersion := "v2.9.1"
	tiller := &Tiller{
		Annotation:     &annotation,
		DefaultVersion: &defaultVersion,
	}

	annotations := map[string]string{annotation: ""}
	if version := tiller.GetVersion(annotations); version != &defaultVersion {
		t.Errorf("Got %s, want %s", *version, defaultVersion)
	}

	want := "v2.8.1"
	annotations = map[string]string{annotation: want}
	if version := tiller.GetVersion(annotations); *version != want {
		t.Errorf("Got %s, want %s", *version, want)
	}

	annotations = make(map[string]string)
	if version := tiller.GetVersion(annotations); version != nil {
		t.Errorf("Got %s, want nil", *version)
	}
}

func TestBuildDeployment(t *testing.T) {
	tiller := &Tiller{}

	namespace := "prd354"
	version := "v2.9.1"
	deployment := tiller.BuildDeployment(&namespace, &version)

	got := deployment.GetNamespace()
	if got != namespace {
		t.Errorf("Got %s, want %s", got, namespace)
	}

	got = deployment.Spec.Template.Spec.Containers[0].Image
	want := fmt.Sprintf("gcr.io/kubernetes-helm/tiller:%s", version)
	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}
}

func TestBuildRoleBinding(t *testing.T) {
	tiller := &Tiller{}

	namespace := "prd354"
	roleBinding := tiller.BuildRoleBinding(&namespace)

	got := roleBinding.GetNamespace()
	if got != namespace {
		t.Errorf("Got %s, want %s", got, namespace)
	}
}

func TestBuildService(t *testing.T) {
	tiller := &Tiller{}

	namespace := "prd354"
	service := tiller.BuildService(&namespace)

	got := service.GetNamespace()
	if got != namespace {
		t.Errorf("Got %s, want %s", got, namespace)
	}
}

func TestBuildServiceAccount(t *testing.T) {
	tiller := &Tiller{}

	namespace := "prd354"
	serviceAccount := tiller.BuildServiceAccount(&namespace)

	got := serviceAccount.GetNamespace()
	if got != namespace {
		t.Errorf("Got %s, want %s", got, namespace)
	}
}
