package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-ldap/ldap"
)

type MockClient struct{}

func (m *MockClient) Start() {}

func (m *MockClient) StartTLS(config *tls.Config) error {
	return nil
}

func (m *MockClient) Close()                   {}
func (m *MockClient) SetTimeout(time.Duration) {}

func (m *MockClient) Bind(username, password string) error {
	return nil
}

func (m *MockClient) SimpleBind(simpleBindRequest *ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	return nil, nil
}

func (m *MockClient) Add(addRequest *ldap.AddRequest) error {
	switch addRequest.DN {
	case "OU=prd367,OU=Kubernetes,OU=Techops,DC=techops,DC=info":
		return errors.New("Unable to add OU")
	case "CN=kubernetes-prd367-test,OU=prd367,OU=Kubernetes,OU=Techops,DC=techops,DC=info":
		return errors.New("Unable to add group")
	}

	return nil
}

func (m *MockClient) Del(delRequest *ldap.DelRequest) error {
	return nil
}

func (m *MockClient) Modify(modifyRequest *ldap.ModifyRequest) error {
	return nil
}

func (m *MockClient) Compare(dn, attribute, value string) (bool, error) {
	return false, nil
}

func (m *MockClient) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	return nil, nil
}

func (m *MockClient) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	switch searchRequest.Filter {
	case "(&(objectClass=organizationalUnit)(ou=prd1811))":
		return nil, errors.New(`LDAP Result Code 32 "No Such Object"`)
	case "(&(objectClass=user)(|(sAMAccountName=mike.goodness)(CN=mike.goodness)))":
		return &ldap.SearchResult{Entries: []*ldap.Entry{&ldap.Entry{DN: "mike.goodness"}}}, nil
	}

	return &ldap.SearchResult{}, nil
}

func (m *MockClient) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, nil
}

func getTestLdap() *LDAP {
	config := &Config{
		AnnotationPrefix: "ldap.ticketmaster.com",
		BaseDN:           "OU=Techops,DC=techops,DC=info",
		CommonOrgUnits:   "OU=Kubernetes",
	}

	return New(config, &MockClient{})
}

func getGoodAnnotations() map[string]string {
	return map[string]string{
		"ldap.ticketmaster.com/admins": "mike.goodness",
		// "ldap.ticketmaster.com/users":  "test4,test5,test6",
	}
}

func getBadAnnotations() map[string]string {
	return map[string]string{
		"admins":     "test,test2,test3",
		"test/users": "test4,test5,test6",
	}
}

func TestAddOU(t *testing.T) {
	l := getTestLdap()

	dn := fmt.Sprintf("OU=prd1811,%s,%s", l.commonOrgUnits, l.baseDN)
	if err := l.addOU(dn); err != nil {
		t.Error(err)
	}

	dn = fmt.Sprintf("OU=prd367,%s,%s", l.commonOrgUnits, l.baseDN)
	if err := l.addOU(dn); err == nil {
		t.Error("Should fail to create OU")
	}
}

func TestAddGroup(t *testing.T) {
	l := getTestLdap()

	dn := fmt.Sprintf("CN=kubernetes-prd1811-test,OU=prd1811,%s,%s", l.commonOrgUnits, l.baseDN)
	if err := l.addOU(dn); err != nil {
		t.Error(err)
	}

	dn = fmt.Sprintf("CN=kubernetes-prd367-test,OU=prd367,%s,%s", l.commonOrgUnits, l.baseDN)
	if err := l.addOU(dn); err == nil {
		t.Error("Should fail to create group")
	}
}

func TestModifyGroupSetMembers(t *testing.T) {
	l := getTestLdap()

	dn := fmt.Sprintf("CN=kubernetes-prd1811-test,OU=prd1811,%s,%s", l.commonOrgUnits, l.baseDN)
	if err := l.modifyGroupSetMembers(dn, []string{}); err == nil {
		t.Error("Should have no group members")
	}

	if err := l.modifyGroupSetMembers(dn, []string{"mike.goodness"}); err != nil {
		t.Error(err)
	}

	if err := l.modifyGroupSetMembers(dn, []string{"test"}); err == nil {
		t.Error("Should have no users found")
	}
}

func TestAddEntries(t *testing.T) {
	l := getTestLdap()

	if err := l.AddEntries("prd1811", getGoodAnnotations()); err != nil {
		t.Error(err)
	}

	if err := l.AddEntries("prd1811", getBadAnnotations()); err == nil {
		t.Error("Should have no valid annotations")
	}
}
