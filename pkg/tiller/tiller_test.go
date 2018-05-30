package tiller

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildResources(t *testing.T) {
	config := &Config{
		Annotation:     "ticketmaster.com/tiller",
		DefaultVersion: "v2.9.1",
	}

	tiller := New(config)

	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"ticketmaster.com/tiller": "",
			},
			Name: "prd1811",
		},
	}

	resources := tiller.BuildResources(namespace)
	if len(resources) != 4 {
		t.Error("Should be 4 elements")
	}

	for _, resource := range resources {
		deployment, ok := resource.(*appsv1.Deployment)
		if ok {
			want := "gcr.io/kubernetes-helm/tiller:v2.9.1"
			got := deployment.Spec.Template.Spec.Containers[0].Image
			if got != want {
				t.Errorf("Wanted %s, got %s", want, got)
			}
		} else {
			continue
		}
	}

	namespace = &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"ticketmaster.com/tiller": "v2.8.1",
			},
			Name: "prd1811",
		},
	}

	resources = tiller.BuildResources(namespace)
	for _, resource := range resources {
		deployment, ok := resource.(*appsv1.Deployment)
		if ok {
			want := "gcr.io/kubernetes-helm/tiller:v2.8.1"
			got := deployment.Spec.Template.Spec.Containers[0].Image
			if got != want {
				t.Errorf("Wanted %s, got %s", want, got)
			}
		} else {
			continue
		}
	}

	namespace = &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "prd1811",
		},
	}

	resources = tiller.BuildResources(namespace)
	if len(resources) != 0 {
		t.Error("Should be empty")
	}
}
