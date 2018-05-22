package pim

import (
	"fmt"
	"testing"

	pimapi "git.tmaws.io/ProductInventoryManagement/go-pimapi"
)

type MockClient struct{}

func (m *MockClient) Get(params map[string]string) (assetItems pimapi.AssetItems, err error) {
	switch params["prdCodes"] {

	// Complete information
	case "PRD1811":
		assetDoc := &pimapi.AssetDoc{
			Name:         "CICD",
			OnCallInfo:   "techops-devxtools-oncall@tmtoc.pagerduty.com",
			SlackChannel: "#gitlab",
			SupportEmail: "Devxtools@ticketmaster.com",
			TechOwner:    "Andy Chan",
		}

		asset := &pimapi.Asset{
			Doc:             *assetDoc,
			ProductCodeName: "PRD1811",
		}

		assetItems.Items = append(assetItems.Items, *asset)

	// Incomplete information
	case "PRD367":
		assetDoc := &pimapi.AssetDoc{
			Name:         "CICD",
			OnCallInfo:   "techops-devxtools-oncall@tmtoc.pagerduty.com",
			SlackChannel: "#gitlab",
			SupportEmail: "Devxtools@ticketmaster.com",
		}

		asset := &pimapi.Asset{
			Doc:             *assetDoc,
			ProductCodeName: "PRD367",
		}

		assetItems.Items = append(assetItems.Items, *asset)

	// Error
	case "TEST9":
		err = fmt.Errorf("")

	// No Product
	default:
	}

	return assetItems, err
}

func TestValidateNamespace(t *testing.T) {
	c := &MockClient{}

	// Test default regex
	p := New("", c)

	name := "prd1811"
	if err := p.ValidateNamespace(name); err != nil {
		t.Error(err)
	}

	name = "test9"
	if err := p.ValidateNamespace(name); err == nil {
		t.Errorf("Namespace %s should be invalid", name)
	}

	// Test custom regex
	p = New(`^test[1-9]{1}$`, c)

	name = "test9"
	if err := p.ValidateNamespace(name); err != nil {
		t.Error(err)
	}

	name = "prd1811"
	if err := p.ValidateNamespace(name); err == nil {
		t.Errorf("Namespace %s should be invalid", name)
	}
}

func TestPrepareAnnotations(t *testing.T) {
	validKeys := [6]string{
		"ticketmaster.com/name",
		"ticketmaster.com/oncall-info",
		"ticketmaster.com/productcode",
		"ticketmaster.com/slack-channel",
		"ticketmaster.com/support-email",
		"ticketmaster.com/tech-owner",
	}

	c := &MockClient{}
	p := New("", c)

	name := "prd1811"
	annotations, err := p.PrepareAnnotations(name)
	if err != nil {
		t.Error(err)
	}

	found := false
	for annotationKey := range annotations {
		for _, validKey := range validKeys {
			if annotationKey == validKey {
				found = true
			}
		}
	}
	if !found {
		t.Error("Annotation(s) missing")
	}

	// No Product found
	name = "prd354"
	_, err = p.PrepareAnnotations(name)
	if err == nil {
		t.Errorf("ProductCode %s should have no Product", name)
	}

	// Incomplete information
	name = "prd367"
	_, err = p.PrepareAnnotations(name)
	if err == nil {
		t.Errorf("Product %s should be missing a field", name)
	}

	// Error
	name = "test9"
	_, err = p.PrepareAnnotations(name)
	if err == nil {
		t.Errorf("Product %s should throw an error", name)
	}
}
