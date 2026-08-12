package migrations

import (
	"embed"
	"fmt"
)

//go:embed *.sql
var files embed.FS

type Migration struct {
	Name string
	Up   string
	Down string
}

func Ordered() ([]Migration, error) {
	names := []string{
		"001_waitlist_entries",
		"002_waitlist_entries_name_optional",
		"003_waitlist_referral_consent",
		"004_referral_events_access_tokens",
		"005_waitlist_marketing_consent",
		"006_status_access_sessions",
	}
	migrations := make([]Migration, 0, len(names))
	for _, name := range names {
		up, err := files.ReadFile(name + ".up.sql")
		if err != nil {
			return nil, fmt.Errorf("read %s up migration: %w", name, err)
		}
		down, err := files.ReadFile(name + ".down.sql")
		if err != nil {
			return nil, fmt.Errorf("read %s down migration: %w", name, err)
		}
		migrations = append(migrations, Migration{Name: name, Up: string(up), Down: string(down)})
	}
	return migrations, nil
}
