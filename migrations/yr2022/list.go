package yr2022

import (
	"github.com/ashishkumar68/auction-api/migrations/yr2022/m5"
	"github.com/go-gormigrate/gormigrate/v2"
)

func GetMigrationsList() []*gormigrate.Migration {
	// ordering is important here so please make sure to not mess up the ordering.
	return []*gormigrate.Migration{
		m5.Version0001,
	}
}
