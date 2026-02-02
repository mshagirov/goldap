package ldapapi

func (api *LdapApi) GetNextIdNumber(tableName string) (string, string, bool) {
	var name, nextId string
	var ok bool

	switch tableName {
	case "Users":
		nextId, ok = api.Cache.Get("nextUidNumber")
		name = "uidNumber"
	case "Groups":
		nextId, ok = api.Cache.Get("nextGidNumber")
		name = "gidNumber"
	default:
		return "", "", false
	}
	return name, nextId, ok
}
