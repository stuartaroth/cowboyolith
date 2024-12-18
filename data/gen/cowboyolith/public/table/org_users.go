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

var OrgUsers = newOrgUsersTable("public", "org_users", "")

type orgUsersTable struct {
	postgres.Table

	// Columns
	OrgID  postgres.ColumnString
	UserID postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type OrgUsersTable struct {
	orgUsersTable

	EXCLUDED orgUsersTable
}

// AS creates new OrgUsersTable with assigned alias
func (a OrgUsersTable) AS(alias string) *OrgUsersTable {
	return newOrgUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new OrgUsersTable with assigned schema name
func (a OrgUsersTable) FromSchema(schemaName string) *OrgUsersTable {
	return newOrgUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new OrgUsersTable with assigned table prefix
func (a OrgUsersTable) WithPrefix(prefix string) *OrgUsersTable {
	return newOrgUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new OrgUsersTable with assigned table suffix
func (a OrgUsersTable) WithSuffix(suffix string) *OrgUsersTable {
	return newOrgUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newOrgUsersTable(schemaName, tableName, alias string) *OrgUsersTable {
	return &OrgUsersTable{
		orgUsersTable: newOrgUsersTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newOrgUsersTableImpl("", "excluded", ""),
	}
}

func newOrgUsersTableImpl(schemaName, tableName, alias string) orgUsersTable {
	var (
		OrgIDColumn    = postgres.StringColumn("org_id")
		UserIDColumn   = postgres.StringColumn("user_id")
		allColumns     = postgres.ColumnList{OrgIDColumn, UserIDColumn}
		mutableColumns = postgres.ColumnList{OrgIDColumn, UserIDColumn}
	)

	return orgUsersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		OrgID:  OrgIDColumn,
		UserID: UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}