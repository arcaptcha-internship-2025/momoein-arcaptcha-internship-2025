package types

import "database/sql"

type Model struct{
	ID string
	CreateAt sql.NullTime
	UpdateAt sql.NullTime
	DeleteAt sql.NullTime
}