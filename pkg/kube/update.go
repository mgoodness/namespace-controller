package kube

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (k *Kubernetes) updateCreateDeployment(namespace *string, deployment *appsv1.Deployment) error {
	name := deployment.GetName()
	client := k.Client.AppsV1().Deployments(*namespace)
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			switch {
			case apiErrors.IsNotFound(err):
				_, err = client.Create(deployment)
				if err != nil {
					return err
				}
				kubeLogger.Debug().Str("operation", "create").Msgf("Created Deployment %s", name)
			default:
				return err
			}
		} else {
			if _, err := client.Update(deployment); err != nil {
				return err
			}
			kubeLogger.Debug().Str("operation", "update").Msgf("Updated Deployment %s", name)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) updateCreateNamespace(namespace *corev1.Namespace) error {
	name := namespace.GetName()
	client := k.Client.CoreV1().Namespaces()
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			switch {
			case apiErrors.IsNotFound(err):
				_, err = client.Create(namespace)
				if err != nil {
					return err
				}
				kubeLogger.Debug().Str("operation", "create").Msgf("Created Namespace %s", name)
			default:
				return err
			}
		} else {
			if _, err := client.Update(namespace); err != nil {
				return err
			}
			kubeLogger.Debug().Str("operation", "update").Msgf("Updated Namespace %s", name)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) updateCreateRoleBinding(namespace *string, roleBinding *rbacv1.RoleBinding) error {
	name := roleBinding.GetName()
	client := k.Client.RbacV1().RoleBindings(*namespace)
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			switch {
			case apiErrors.IsNotFound(err):
				_, err = client.Create(roleBinding)
				if err != nil {
					return err
				}
				kubeLogger.Debug().Str("operation", "create").Msgf("Created RoleBinding %s", name)
			default:
				return err
			}
		} else {
			if _, err := client.Update(roleBinding); err != nil {
				return err
			}
			kubeLogger.Debug().Str("operation", "update").Msgf("Updated RoleBinding %s", name)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) updateCreateService(namespace *string, service *corev1.Service) error {
	name := service.GetName()
	client := k.Client.CoreV1().Services(*namespace)
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		existing, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			switch {
			case apiErrors.IsNotFound(err):
				_, err = client.Create(service)
				if err != nil {
					return err
				}
				kubeLogger.Debug().Str("operation", "create").Msgf("Created Service %s", name)
			default:
				return err
			}
		} else {
			service.SetResourceVersion(existing.GetResourceVersion())
			service.Spec.ClusterIP = existing.Spec.ClusterIP
			if _, err := client.Update(service); err != nil {
				return err
			}
			kubeLogger.Debug().Str("operation", "update").Msgf("Updated Service %s", name)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubernetes) updateCreateServiceAccount(namespace *string, serviceAccount *corev1.ServiceAccount) error {
	name := serviceAccount.GetName()
	client := k.Client.CoreV1().ServiceAccounts(*namespace)
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			switch {
			case apiErrors.IsNotFound(err):
				_, err = client.Create(serviceAccount)
				if err != nil {
					return err
				}
				kubeLogger.Debug().Str("operation", "create").Msgf("Created ServiceAccount %s", name)
			default:
				return err
			}
		} else {
			if _, err := client.Update(serviceAccount); err != nil {
				return err
			}
			kubeLogger.Debug().Str("operation", "update").Msgf("Updated ServiceAccount %s", name)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
