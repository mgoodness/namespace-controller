package ldap

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	ldapapi "github.com/go-ldap/ldap"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type searchRequest struct {
	attributes []string
	baseDN     *string
	filter     string
	sizeLimit  int
}

// Config struct contains LDAP config
type Config struct {
	AnnotationPrefix string
	BaseDN           string
	CommonOrgUnits   string
	Enabled          bool
	Hostname         string
	Password         string
	Username         string
}

// LDAP struct contains an LDAP client
type LDAP struct {
	annotationPrefix *string
	baseDN           *string
	client           ldapapi.Client
	commonOrgUnits   *string
	namespace        *string
}

var ldapLogger zerolog.Logger

func init() {
	ldapLogger = log.With().Str("component", "ldap").Logger()
}

// New creates a new LDAP struct
func New(config *Config, ldapClient ldapapi.Client) (*LDAP, error) {

	if err := ldapClient.Bind(config.Username, config.Password); err != nil {
		return nil, err
	}

	return &LDAP{
		annotationPrefix: &config.AnnotationPrefix,
		baseDN:           &config.BaseDN,
		client:           ldapClient,
		commonOrgUnits:   &config.CommonOrgUnits,
	}, nil
}

// AddEntries creates LDAP entries
func (l *LDAP) AddEntries(namespace *string, annotations map[string]string) error {
	var found bool

	for key, value := range annotations {
		keySplit := strings.Split(key, "/")
		if len(keySplit) != 2 {
			continue
		}

		if keySplit[0] == *l.annotationPrefix {
			found = true

			dn := fmt.Sprintf("CN=kubernetes-%s-%s,OU=%s,%s,%s", *namespace, keySplit[1], *namespace, *l.commonOrgUnits, *l.baseDN)
			members := strings.Split(strings.Replace(value, " ", "", -1), ",")

			_, ouDN := splitDN(&dn)
			if err := l.createOU(&ouDN); err != nil {
				return err
			}

			if err := l.createGroup(&dn); err != nil {
				return err
			}

			if err := l.modifyGroupSetMembers(&dn, members); err != nil {
				return err
			}
		}
	}

	if !found {
		return fmt.Errorf("No annotations found with prefix %s", *l.annotationPrefix)
	}

	return nil
}

func (l *LDAP) createOU(dn *string) error {
	name, base := splitDN(dn)

	entries, err := l.searchOU(&base, &name)
	if err != nil {
		if strings.HasPrefix(err.Error(), `LDAP Result Code 32 "No Such Object"`) {
			if err = l.createOU(&base); err != nil {
				return err
			}
		}
		return err
	}

	if len(entries) != 0 {
		ldapLogger.Debug().Msgf("Found %s", *dn)
		return nil
	}

	request := ldapapi.NewAddRequest(*dn)
	request.Attribute("objectClass", []string{"organizationalUnit", "top"})
	request.Attribute("ou", []string{name})
	request.Attribute("name", []string{name})

	if err = l.client.Add(request); err != nil {
		return err
	}

	ldapLogger.Debug().Msgf("Created %s", *dn)
	return nil
}

func (l *LDAP) createGroup(dn *string) error {
	name, base := splitDN(dn)

	var entries []*ldapapi.Entry
	entries, err := l.searchGroup(&base, &name)
	if err != nil {
		return err
	}

	if len(entries) != 0 {
		ldapLogger.Debug().Msgf("Found %s", *dn)
		return nil
	}

	request := ldapapi.NewAddRequest(*dn)
	request.Attribute("objectClass", []string{"group", "top"})
	request.Attribute("cn", []string{name})
	request.Attribute("name", []string{name})
	request.Attribute("sAMAccountName", []string{name})

	err = l.client.Add(request)
	if err != nil {
		return err
	}

	ldapLogger.Debug().Msgf("Created %s", *dn)
	return nil
}

func (l *LDAP) modifyGroupSetMembers(dn *string, members []string) error {
	if len(members) == 0 {
		return errors.New("No group members provided")
	}

	var accounts []string

	for i := 0; i < len(members); i++ {
		results, err := l.searchMember(&members[i])
		if err != nil {
			return err
		}

		if len(results) == 0 {
			return fmt.Errorf("User %s not found", members[i])
		}

		for _, entry := range results {
			accounts = append(accounts, entry.DN)
		}
	}

	// TODO: get current membership, only update if there are changes

	request := ldapapi.NewModifyRequest(*dn)
	request.Replace("member", accounts)

	if err := l.client.Modify(request); err != nil {
		return err
	}

	ldapLogger.Debug().Msgf("Updated members of %s", *dn)
	return nil
}

func (l *LDAP) search(r *searchRequest) ([]*ldapapi.Entry, error) {
	request := ldapapi.NewSearchRequest(
		*r.baseDN, ldapapi.ScopeWholeSubtree, ldapapi.NeverDerefAliases, r.sizeLimit, 0, false, r.filter, r.attributes, nil,
	)

	search, err := l.client.Search(request)
	if err != nil {
		return nil, err
	}

	return search.Entries, nil
}

func (l *LDAP) searchOU(dn, ou *string) (results []*ldapapi.Entry, err error) {
	request := &searchRequest{
		baseDN:     dn,
		sizeLimit:  2,
		filter:     fmt.Sprintf("(&(objectClass=organizationalUnit)(ou=%s))", *ou),
		attributes: []string{"dn"},
	}

	return l.search(request)
}

func (l *LDAP) searchGroup(dn, cn *string) (results []*ldapapi.Entry, err error) {
	request := &searchRequest{
		baseDN:     dn,
		sizeLimit:  2,
		filter:     fmt.Sprintf("(&(objectClass=group)(cn=%s))", *cn),
		attributes: []string{"dn", "member"},
	}

	return l.search(request)
}

func (l *LDAP) searchMember(cn *string) (results []*ldapapi.Entry, err error) {
	request := &searchRequest{
		baseDN:     l.baseDN,
		sizeLimit:  1000,
		filter:     fmt.Sprintf("(&(objectClass=user)(|(sAMAccountName=%s)(CN=%s)))", *cn, *cn),
		attributes: []string{"dn"},
	}

	return l.search(request)
}

func splitDN(dn *string) (name, base string) {
	re := regexp.MustCompile(`^[^=]+=([^,]+),(.+)`)
	submatch := re.FindStringSubmatch(*dn)
	if len(submatch) == 3 {
		name = submatch[1]
		base = submatch[2]
	}

	return
}
