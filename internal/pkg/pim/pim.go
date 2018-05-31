package pim

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	pimapi "git.tmaws.io/ProductInventoryManagement/go-pimapi"
	"github.com/cenkalti/backoff"
)

const defaultPimRegex string = `^prd[1-9]{1,5}$`

type AssetClient interface {
	Get(map[string]string) (pimapi.AssetItems, error)
}

type PIM struct {
	assetClient AssetClient
	regex       string
}

func New(regex string, assetClient AssetClient) *PIM {
	if regex == "" {
		regex = defaultPimRegex
	}

	return &PIM{
		assetClient: assetClient,
		regex:       regex,
	}
}

func (p *PIM) ValidateNamespace(name string) error {
	if match, _ := regexp.MatchString(p.regex, name); match != true {
		return fmt.Errorf("Must match regex %s", p.regex)
	}

	return nil
}

func (p *PIM) PrepareAnnotations(name string) (map[string]string, error) {
	products, err := p.fetchProducts(&name)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("No Product found using ProductCode %s", name)
	}

	if len(products[0].Doc.TechOwner) == 0 || len(products[0].Doc.SlackChannel) == 0 ||
		len(products[0].Doc.SupportEmail) == 0 || len(products[0].Doc.OnCallInfo) == 0 {
		return nil, errors.New("Some mandatory information is missing. Update PIM with your TechOwner, SlackChannel, SupportEmail, OnCallInfo")
	}

	return map[string]string{
		"ticketmaster.com/name":          products[0].Doc.Name,
		"ticketmaster.com/oncall-info":   products[0].Doc.OnCallInfo,
		"ticketmaster.com/productcode":   products[0].ProductCodeName,
		"ticketmaster.com/slack-channel": products[0].Doc.SlackChannel,
		"ticketmaster.com/support-email": products[0].Doc.SupportEmail,
		"ticketmaster.com/tech-owner":    products[0].Doc.TechOwner,
	}, nil
}

func (p *PIM) fetchProducts(prdCode *string) ([]pimapi.Asset, error) {
	var products pimapi.AssetItems

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second

	getProducts := func() (err error) {
		if products, err = p.assetClient.Get(
			map[string]string{
				"assetTypeName": "Product",
				"deleted":       "false",
				"prdCodes":      strings.ToUpper(*prdCode),
			},
		); err != nil {
			return err
		}

		return nil
	}

	if err := backoff.Retry(getProducts, b); err != nil {
		return nil, err
	}

	return products.Items, nil
}
