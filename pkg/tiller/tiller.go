package tiller

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Config struct {
	Annotation     string
	DefaultVersion string
}

type Tiller struct {
	annotation     string
	defaultVersion string
}

func New(config *Config) *Tiller {
	return &Tiller{
		annotation:     config.Annotation,
		defaultVersion: config.DefaultVersion,
	}
}

func (t *Tiller) BuildResources(namespace *apiv1.Namespace) (resources []runtime.Object) {
	var version string

	for key, value := range namespace.GetAnnotations() {
		if key == t.annotation {
			switch value {
			case "":
				version = t.defaultVersion
			default:
				version = value
			}
			break
		}
	}

	if version == "" {
		return
	}

	ns := namespace.GetName()

	resources = append(resources, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "helm",
			},
			Name: "tiller",
		},
		Spec: appsv1.DeploymentSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "helm",
						"tier": "cs",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "tiller",
							Image: "gcr.io/kubernetes-helm/tiller:" + version,
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "TILLER_NAMESPACE",
									Value: ns,
								},
							},
							LivenessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.FromInt(44135),
									},
								},
							},
							Ports: []apiv1.ContainerPort{
								apiv1.ContainerPort{
									ContainerPort: 44134,
								},
							},
							ReadinessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.FromInt(44135),
									},
								},
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Requests: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse("250m"),
									apiv1.ResourceMemory: resource.MustParse("64Mi"),
								},
							},
						},
					},
					ServiceAccountName: "tiller",
				},
			},
		},
	})

	resources = append(resources, &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"prometheus.io/path":             "/metrics",
				"prometheus.io/port":             "44135",
				"prometheus.io/scrape":           "true",
				"ticketmaster.com/productcode":   ns,
				"ticketmaster.com/inventorycode": "tiller",
			},
			Labels: map[string]string{
				"app": "helm",
			},
			Name: "tiller",
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port: 44134,
				},
			},
			Selector: map[string]string{
				"app": "helm",
			},
		},
	})

	resources = append(resources, &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "helm",
			},
			Name: "tiller",
		},
	})

	resources = append(resources, &rbac.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "helm",
			},
			Name: "tiller",
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "tm:tiller",
		},
		Subjects: []rbac.Subject{
			rbac.Subject{
				Kind: "ServiceAccount",
				Name: "tiller",
			},
		},
	})

	return
}
