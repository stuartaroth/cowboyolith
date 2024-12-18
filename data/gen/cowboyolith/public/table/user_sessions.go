//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var UserSessions = newUserSessionsTable("public", "user_sessions", "")

type userSessionsTable struct {
	postgres.Table

	// Columns
	ID                postgres.ColumnString
	UserID            postgres.ColumnString
	HashedCookieToken postgres.ColumnString
	Salt              postgres.ColumnString
	UserAgent         postgres.ColumnString
	InsertedAt        postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UserSessionsTable struct {
	userSessionsTable

	EXCLUDED userSessionsTable
}

// AS creates new UserSessionsTable with assigned alias
func (a UserSessionsTable) AS(alias string) *UserSessionsTable {
	return newUserSessionsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UserSessionsTable with assigned schema name
func (a UserSessionsTable) FromSchema(schemaName string) *UserSessionsTable {
	return newUserSessionsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UserSessionsTable with assigned table prefix
func (a UserSessionsTable) WithPrefix(prefix string) *UserSessionsTable {
	return newUserSessionsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UserSessionsTable with assigned table suffix
func (a UserSessionsTable) WithSuffix(suffix string) *UserSessionsTable {
	return newUserSessionsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUserSessionsTable(schemaName, tableName, alias string) *UserSessionsTable {
	return &UserSessionsTable{
		userSessionsTable: newUserSessionsTableImpl(schemaName, tableName, alias),
		EXCLUDED:          newUserSessionsTableImpl("", "excluded", ""),
	}
}

func newUserSessionsTableImpl(schemaName, tableName, alias string) userSessionsTable {
	var (
		IDColumn                = postgres.StringColumn("id")
		UserIDColumn            = postgres.StringColumn("user_id")
		HashedCookieTokenColumn = postgres.StringColumn("hashed_cookie_token")
		SaltColumn              = postgres.StringColumn("salt")
		UserAgentColumn         = postgres.StringColumn("user_agent")
		InsertedAtColumn        = postgres.TimestampColumn("inserted_at")
		allColumns              = postgres.ColumnList{IDColumn, UserIDColumn, HashedCookieTokenColumn, SaltColumn, UserAgentColumn, InsertedAtColumn}
		mutableColumns          = postgres.ColumnList{UserIDColumn, HashedCookieTokenColumn, SaltColumn, UserAgentColumn, InsertedAtColumn}
	)

	return userSessionsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                IDColumn,
		UserID:            UserIDColumn,
		HashedCookieToken: HashedCookieTokenColumn,
		Salt:              SaltColumn,
		UserAgent:         UserAgentColumn,
		InsertedAt:        InsertedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
