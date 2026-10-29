package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mshagirov/goldap/internal/cache"
	"github.com/mshagirov/goldap/internal/tabs"
	"github.com/mshagirov/goldap/ldapapi"
)

func createNewTabsModel(tableIdx int, rowIndices []int, api *ldapapi.LdapApi) (tabs.Model, error) {
	// Stage the cache as well as the tables: a failed reload must leave both intact.
	loadingAPI := *api
	loadingAPI.Cache = cache.NewCache()
	contents, err := loadTables(ldapapi.TableNames, loadingAPI.GetTableInfo)
	if err != nil {
		return tabs.Model{}, err
	}
	dn := make([][]string, len(contents))
	for i, t := range contents {
		dn[i] = t.DN
	}

	m := tabs.NewTabsModel(ldapapi.TableNames, contents, dn, &loadingAPI)

	m.State.TabId = tableIdx
	m.State.Table = tabs.NewTable(contents[tableIdx])
	for i, rowId := range rowIndices {
		m.State.TabSates[i].Cursor = min(len(dn[i]), rowId)
	}
	m.SetTable()
	api.Cache = loadingAPI.Cache
	m.LdapApi = api
	return m, nil
}

func loadTables(names []string, load func(string) (ldapapi.TableInfo, error)) ([]ldapapi.TableInfo, error) {
	contents := make([]ldapapi.TableInfo, len(names))
	errs := make([]error, len(names))
	var wg sync.WaitGroup
	for i, name := range names {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t, err := load(name)
			if err != nil {
				errs[i] = fmt.Errorf("Load %s Error; %w", name, err)
				return
			}
			contents[i] = t
		}()
	}
	// Even on failure, all workers must finish before discarding the staged cache.
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return contents, nil
}
