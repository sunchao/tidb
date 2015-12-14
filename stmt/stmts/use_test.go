// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package stmts_test

import (
	. "github.com/pingcap/check"
	"github.com/pingcap/tidb"
	"github.com/pingcap/tidb/stmt/stmts"
)

func (s *testStmtSuite) TestUse(c *C) {
	testSQL := `create database if not exists use_test;`
	mustExec(c, s.testDB, testSQL)

	testSQL = `use test;`
	stmtList, err := tidb.Compile(s.ctx, testSQL)
	c.Assert(err, IsNil)
	c.Assert(stmtList, HasLen, 1)

	testStmt, ok := stmtList[0].(*stmts.UseStmt)
	c.Assert(ok, IsTrue)

	c.Assert(testStmt.IsDDL(), IsFalse)
	c.Assert(len(testStmt.OriginText()), Greater, 0)

	mf := newMockFormatter()
	testStmt.Explain(nil, mf)
	c.Assert(mf.Len(), Greater, 0)

	errTestSQL := `use xxx;`
	tx := mustBegin(c, s.testDB)
	_, err = tx.Exec(errTestSQL)
	c.Assert(err, NotNil)
	tx.Rollback()
}

func (s *testStmtSuite) TestCharsetDatabase(c *C) {
	testSQL := `create database if not exists cd_test_utf8 CHARACTER SET utf8 COLLATE utf8_bin;`
	mustExec(c, s.testDB, testSQL)

	testSQL = `create database if not exists cd_test_latin1 CHARACTER SET latin1 COLLATE latin1_swedish_ci;`
	mustExec(c, s.testDB, testSQL)

	testSQL = `use cd_test_utf8;`
	mustExec(c, s.testDB, testSQL)

	tx := mustBegin(c, s.testDB)
	rows, err := tx.Query(`select @@character_set_database;`)
	c.Assert(err, IsNil)
	matchRows(c, rows, [][]interface{}{{"utf8"}})
	rows, err = tx.Query(`select @@collation_database;`)
	c.Assert(err, IsNil)
	matchRows(c, rows, [][]interface{}{{"utf8_bin"}})
	mustCommit(c, tx)

	testSQL = `use cd_test_latin1;`
	mustExec(c, s.testDB, testSQL)

	tx = mustBegin(c, s.testDB)
	rows, err = tx.Query(`select @@character_set_database;`)
	c.Assert(err, IsNil)
	matchRows(c, rows, [][]interface{}{{"latin1"}})
	rows, err = tx.Query(`select @@collation_database;`)
	c.Assert(err, IsNil)
	matchRows(c, rows, [][]interface{}{{"latin1_swedish_ci"}})
	mustCommit(c, tx)
}
