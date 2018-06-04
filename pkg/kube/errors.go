package kube

import "fmt"

type KindError struct {
	kind string
}

func (e *KindError) Error() string {
	return fmt.Sprintf("File is not a %s manifest", e.kind)
}
