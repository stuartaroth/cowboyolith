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

var OrgGroupUsers = newOrgGroupUsersTable("public", "org_group_users", "")

type orgGroupUsersTable struct {
	postgres.Table

	// Columns
	OrgGroupID postgres.ColumnString
	UserID     postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type OrgGroupUsersTable struct {
	orgGroupUsersTable

	EXCLUDED orgGroupUsersTable
}

// AS creates new OrgGroupUsersTable with assigned alias
func (a OrgGroupUsersTable) AS(alias string) *OrgGroupUsersTable {
	return newOrgGroupUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new OrgGroupUsersTable with assigned schema name
func (a OrgGroupUsersTable) FromSchema(schemaName string) *OrgGroupUsersTable {
	return newOrgGroupUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new OrgGroupUsersTable with assigned table prefix
func (a OrgGroupUsersTable) WithPrefix(prefix string) *OrgGroupUsersTable {
	return newOrgGroupUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new OrgGroupUsersTable with assigned table suffix
func (a OrgGroupUsersTable) WithSuffix(suffix string) *OrgGroupUsersTable {
	return newOrgGroupUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newOrgGroupUsersTable(schemaName, tableName, alias string) *OrgGroupUsersTable {
	return &OrgGroupUsersTable{
		orgGroupUsersTable: newOrgGroupUsersTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newOrgGroupUsersTableImpl("", "excluded", ""),
	}
}

func newOrgGroupUsersTableImpl(schemaName, tableName, alias string) orgGroupUsersTable {
	var (
		OrgGroupIDColumn = postgres.StringColumn("org_group_id")
		UserIDColumn     = postgres.StringColumn("user_id")
		allColumns       = postgres.ColumnList{OrgGroupIDColumn, UserIDColumn}
		mutableColumns   = postgres.ColumnList{OrgGroupIDColumn, UserIDColumn}
	)

	return orgGroupUsersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		OrgGroupID: OrgGroupIDColumn,
		UserID:     UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
