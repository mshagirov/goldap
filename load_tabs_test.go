package main

import (
	"errors"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mshagirov/goldap/internal/cache"
	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/ldapapi"
)

func TestCreateNewTabsModelFailurePreservesCache(t *testing.T) {
	api := &ldapapi.LdapApi{
		Config: config.Config{LdapUrl: "unsupported://localhost"},
		Cache:  cache.NewCache(),
	}
	entries := map[string]string{
		"uid=alice":      "uid=alice,dc=goldap,dc=sh",
		"gidNumber=1000": "cn=staff,dc=goldap,dc=sh",
		"nextUidNumber":  "1001",
		"nextGidNumber":  "1001",
	}
	for key, value := range entries {
		api.Cache.Add(key, value)
	}

	model, err := createNewTabsModel(1, []int{2, 3, 0}, api)
	if err == nil || !strings.Contains(err.Error(), "DialURL Error") {
		t.Fatalf("error = %v, want LDAP dial failure", err)
	}
	if model.State != nil {
		t.Fatal("failed reload returned a partial model")
	}
	for key, want := range entries {
		if got, ok := api.Cache.Get(key); !ok || got != want {
			t.Errorf("cache[%q] = %q, %v; want %q", key, got, ok, want)
		}
	}
}

func TestLoadTablesConcurrent(t *testing.T) {
	for _, fail := range []bool{false, true} {
		name := "preserves order"
		if fail {
			name = "waits for all workers and reports failures"
		}
		t.Run(name, func(t *testing.T) {
			names := []string{"Users", "Groups", "OrgUnits"}
			started := make(chan string, len(names))
			finished := make(chan string, len(names))
			release := make(chan struct{})
			unblock := sync.OnceFunc(func() { close(release) })
			defer unblock()
			failure := errors.New("LDAP unavailable")
			type result struct {
				contents []ldapapi.TableInfo
				err      error
			}
			done := make(chan result, 1)
			go func() {
				contents, err := loadTables(names, func(name string) (ldapapi.TableInfo, error) {
					started <- name
					// An early failure must not let loadTables return while other workers run.
					if fail && name == "Users" {
						finished <- name
						return ldapapi.TableInfo{}, failure
					}
					<-release
					finished <- name
					if fail && name == "Groups" {
						return ldapapi.TableInfo{}, failure
					}
					return ldapapi.TableInfo{DN: []string{name}}, nil
				})
				done <- result{contents, err}
			}()

			seen := make(map[string]bool)
			for range names {
				select {
				case name := <-started:
					seen[name] = true
				case <-time.After(5 * time.Second):
					t.Fatal("table loads did not all start concurrently")
				}
			}
			if len(seen) != len(names) {
				t.Fatalf("started tables = %v, want %v", seen, names)
			}
			select {
			case <-done:
				t.Fatal("returned before all loaders finished")
			default:
			}

			unblock()
			select {
			case got := <-done:
				if len(finished) != len(names) {
					t.Fatal("returned before all loaders completed")
				}
				if fail {
					if !errors.Is(got.err, failure) || !strings.Contains(got.err.Error(), "Users") || !strings.Contains(got.err.Error(), "Groups") {
						t.Fatalf("error = %v, want both table failures", got.err)
					}
					if got.contents != nil {
						t.Fatal("failed load returned partial contents")
					}
					return
				}
				if got.err != nil {
					t.Fatal(got.err)
				}
				var order []string
				for _, table := range got.contents {
					order = append(order, table.DN[0])
				}
				if !slices.Equal(order, names) {
					t.Errorf("table order = %v, want %v", order, names)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("table loads did not finish")
			}
		})
	}
}
