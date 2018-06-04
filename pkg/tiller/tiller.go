package tiller

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Config struct {
	Annotation     string
	DefaultVersion string
}

type Tiller struct {
	Annotation     *string
	DefaultVersion *string
}

var tillerLogger zerolog.Logger

func init() {
	tillerLogger = log.With().Str("component", "tiller").Logger()
}

func (t *Tiller) GetVersion(annotations map[string]string) *string {
	for key, value := range annotations {
		if key == *t.Annotation {
			spew.Dump(key, value)

			switch value {
			case "":
				return t.DefaultVersion
			default:
				return &value
			}
		}
	}

	return nil
}

func (t *Tiller) BuildDeployment(namespace *string, version *string) (deployment *appsv1.Deployment) {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "helm",
			},
			Name:      "tiller",
			Namespace: *namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "helm"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "helm",
						"tier": "cs",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "tiller",
							Image: "gcr.io/kubernetes-helm/tiller:" + *version,
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "TILLER_NAMESPACE",
									Value: *namespace,
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.FromInt(44135),
									},
								},
							},
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									ContainerPort: 44134,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.FromInt(44135),
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
							},
						},
					},
					ServiceAccountName: "tiller",
				},
			},
		},
	}
}

func (t *Tiller) BuildRoleBinding(namespace *string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "helm",
			},
			Name:      "tiller",
			Namespace: *namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "tm:tiller",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind: "ServiceAccount",
				Name: "tiller",
			},
		},
	}
}

func (t *Tiller) BuildService(namespace *string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"prometheus.io/path":             "/metrics",
				"prometheus.io/port":             "44135",
				"prometheus.io/scrape":           "true",
				"ticketmaster.com/productcode":   *namespace,
				"ticketmaster.com/inventorycode": "tiller",
			},
			Labels: map[string]string{
				"app": "helm",
			},
			Name:      "tiller",
			Namespace: *namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: 44134,
				},
			},
			Selector: map[string]string{
				"app": "helm",
			},
		},
	}
}

func (t *Tiller) BuildServiceAccount(namespace *string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "helm",
			},
			Name:      "tiller",
			Namespace: *namespace,
		},
	}
}
