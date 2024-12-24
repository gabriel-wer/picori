package api

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type LDAP struct {
	Server string
	Port   int
	baseDN string
	Conn   *ldap.Conn
}

func NewLDAPConnection() (*LDAP, error) {
	//https://www.forumsys.com/2022/05/10/online-ldap-test-server/
	ldap_data := &LDAP{
		Server: "ldap://ldap.forumsys.com",
		Port:   389,
		baseDN: "dc=example,dc=com",
	}
	conn, err := ldap.DialURL(fmt.Sprintf("%s:%d", ldap_data.Server, ldap_data.Port))
	if err != nil {
		return nil, fmt.Errorf("Failed to Connect to LDAP: %v", err)
	}

	ldap_data.Conn = conn

	return ldap_data, nil
}

func (s *LDAP) BindToLDAP() error {
	userDN := fmt.Sprintf("cn=read-only-admin,%s", s.baseDN)
	err := s.Conn.Bind(userDN, "password")
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	return nil
}

func (s *LDAP) SearchUser(username string) (string, error) {
	filter := fmt.Sprintf("(uid=%s)", username)

	searchRequest := ldap.NewSearchRequest(
		s.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"dn"},
		nil,
	)

	result, err := s.Conn.Search(searchRequest)
	if err != nil {
		return "", fmt.Errorf("failed to search for user: %v", err)
	}

	if len(result.Entries) == 0 {
		return "", fmt.Errorf("user not found")
	}

	return result.Entries[0].DN, nil
}

func (s *LDAP) AuthenticateUser(username, password string) error {
	err := s.BindToLDAP()
	if err != nil {
		return err
	}

	userDN, err := s.SearchUser(username)
	if err != nil {
		return err
	}

	err = s.Conn.Bind(userDN, password)
	if err != nil {
		return err
	}

	return nil
}
