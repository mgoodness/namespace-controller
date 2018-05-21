package plugins

// AnnotationValidator validates the existence & format of Kubernetes Namespace Annotations
type AnnotationValidator interface {
	ValidateAnnotations(annotations map[string]string) error
}

// NamespaceValidator validates the format of a Kubernetes Namespace
type NamespaceValidator interface {
	ValidateNamespace(string) error
}

// NamespaceAnnotator prepares a map of Annotations to be added to a Kubernetes Namespace
type NamespaceAnnotator interface {
	PrepareAnnotations(string) (map[string]string, error)
}
