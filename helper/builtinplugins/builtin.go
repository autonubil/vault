package builtinplugins

import (
	"github.com/autonubil/vault/plugins/database/cassandra"
	"github.com/autonubil/vault/plugins/database/mssql"
	"github.com/autonubil/vault/plugins/database/mysql"
	"github.com/autonubil/vault/plugins/database/postgresql"
)

type BuiltinFactory func() (interface{}, error)

var plugins map[string]BuiltinFactory = map[string]BuiltinFactory{
	// These four plugins all use the same mysql implementation but with
	// different username settings passed by the constructor.
	"mysql-database-plugin":        mysql.New(mysql.DisplayNameLen, mysql.UsernameLen),
	"mysql-aurora-database-plugin": mysql.New(mysql.LegacyDisplayNameLen, mysql.LegacyUsernameLen),
	"mysql-rds-database-plugin":    mysql.New(mysql.LegacyDisplayNameLen, mysql.LegacyUsernameLen),
	"mysql-legacy-database-plugin": mysql.New(mysql.LegacyDisplayNameLen, mysql.LegacyUsernameLen),

	"postgresql-database-plugin": postgresql.New,
	"mssql-database-plugin":      mssql.New,
	"cassandra-database-plugin":  cassandra.New,
}

func Get(name string) (BuiltinFactory, bool) {
	f, ok := plugins[name]
	return f, ok
}

func Keys() []string {
	keys := make([]string, len(plugins))

	i := 0
	for k := range plugins {
		keys[i] = k
		i++
	}

	return keys
}
