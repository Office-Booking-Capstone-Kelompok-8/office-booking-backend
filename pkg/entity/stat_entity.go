package entity

import "database/sql"

type TimeframeStat struct {
	Day   sql.NullInt64
	Week  sql.NullInt64
	Month sql.NullInt64
	Year  sql.NullInt64
	All   sql.NullInt64
}
