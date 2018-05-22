package config

import (
	"path/filepath"
	"testing"

	pimapi "git.tmaws.io/ProductInventoryManagement/go-pimapi"
	"github.com/google/go-cmp/cmp"
	"github.com/mgoodness/namespace-controller/pkg/ldap"
)

func TestNew(t *testing.T) {
	configFile := filepath.Join("testdata", "config.toml")
	got, err := New(&configFile)
	if err != nil {
		t.Error(err)
	}

	ldapConfig := ldap.Config{
		AnnotationPrefix: "ldap.ticketmaster.com",
		BaseDN:           "OU=Techops,DC=techops,DC=info",
		CommonOrgUnits:   "OU=Kubernetes",
		Enabled:          false,
		Hostname:         "ldaps.tmaws.io",
		Password:         "fakepass",
		Username:         "fakeuser",
	}

	pimConfig := pimapi.Config{
		Token:      "",
		APIVersion: "v1",
		BaseURL:    "http://pim.tmaws.io",
	}

	want := &Config{
		Debug:      false,
		Kubeconfig: "",
		Ldap:       ldapConfig,
		Manifests:  "",
		Namespaces: "",
		Pim:        pimConfig,
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}