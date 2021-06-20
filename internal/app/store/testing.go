package store

import (
	"strings"
	"testing"
)

func Teststore(t *testing.T, databaseUrl string) (*Store, func(... string)) {
	t.Helper()
	config := NewConfig()
	config.DataBaseUrl = databaseUrl
	s := NewStore(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	return s, func(tables ... string){
		if len(tables) > 0{
			if _, err:= s.db.Exec("TRUNCATE %s CASCADE", strings.Join(tables, ", ")); err != nil {
				t.Fatal(err)
			}
		}
		s.Close()
	}

}
