package config

import (
	"fmt"
)

func ExampleJson() string {
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Println(err)
	}
	exampleJson := `{
  "LDAP_URL":"ldap://localhost:389",
  "LDAP_BASE_DN":"dc=goldap,dc=sh",
  "LDAP_ADMIN_DN":"cn=admin,dc=goldap,dc=sh"
}

| Configuration also supports simple templates and default values for
| new LDAP entries. E.g.: you may add following after the LDAP server
|	parameters above (append comma "," to the last "LDAP_ADMIN_DN" line):

{
	.... LDAP server parameters ....

  "Users": [
    { "Name": "ou", "Value": [ "People" ] },
    { "Name": "objectClass", "Value": [ "posixAccount", "inetOrgPerson", "top" ] },
    { "Name": "employeeType", "Value": [ "Staff" ] },
    { "Name": "uid", "Value": [ "USERNAME" ] },
    { "Name": "givenName", "Value": [ "NAME" ] },
    { "Name": "sn", "Value": [ "SURNAME" ] },
    { "Name": "mail", "Value": [ "{{uid}}@goldap.sh" ] },
    { "Name": "uidNumber", "Value": [ "1234" ] },
    { "Name": "gidNumber", "Value": [ "Enter group gidNumber, e.g., 100" ] },
    { "Name": "homeDirectory", "Value": [ "/home/{{uid}}" ] },
    { "Name": "userPassword", "Value": [ "password" ] },
    { "Name": "description", "Value": [ "description" ] }
  ],
  "Groups": [
    { "Name": "ou", "Value": [ "" ] },
    { "Name": "objectClass", "Value": [ "top", "posixGroup" ] },
    { "Name": "cn", "Value": [ "Group's name" ] },
    { "Name": "gidNumber", "Value": [ "" ] },
    { "Name": "member", "Value": [ "asimov", "hseldon", "rdolivaw", "..." ] },
    { "Name": "memberUid", "Value": [ "asimov", "hseldon", "rdolivaw", "..." ] }
  ],
  "OrgUnits": [
    { "Name": "ou", "Value": [ "" ] },
    { "Name": "objectClass", "Value": [ "top", "organizationalUnit" ] },
    { "Name": "description", "Value": [ "description" ] }
  ]
}
`
	return fmt.Sprintf(`goldap config: LdapUrl is empty

Create a JSON configuration file in:
  %s
Example configuration file contents:
%s
`, configPath, exampleJson)
}
