package sqlengine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/sqlengine/driver"
	"bitbucket.org/goreorto/benchwell/sqlengine/sqlenginetest"
)

func registerDriver(t *testing.T, driverName string) (*sqlenginetest.MockDriver, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	drv := sqlenginetest.NewMockDriver(ctrl)
	driver.RegisterDriver(driverName, drv)

	return drv, ctrl
}

func TestEngineConnect(t *testing.T) {
	type testcase struct {
		name         string
		dsn          string
		conn         driver.Connection
		expectations func(*testing.T, *testcase, *sqlenginetest.MockDriver)
		err          error
	}

	testcases := []testcase{
		{
			name: "connection_error",
			dsn:  "%s://user:pass@localhost:3306/testdb",
			expectations: func(t *testing.T, tc *testcase, d *sqlenginetest.MockDriver) {
				d.EXPECT().
					Connect("test0://user:pass@localhost:3306/testdb").
					Return(nil, tc.err)
			},
			err: errors.New("err"),
		},
		{
			name: "ok",
			dsn:  "%s://user:pass@localhost:3306/testdb",
			expectations: func(t *testing.T, tc *testcase, d *sqlenginetest.MockDriver) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				tc.conn = sqlenginetest.NewMockConnection(ctrl)
				d.EXPECT().
					Connect("test1://user:pass@localhost:3306/testdb").
					Return(tc.conn, tc.err)
			},
			err: nil,
		},
	}

	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			driverName := fmt.Sprintf("test%d", i)
			driver, ctrl := registerDriver(t, driverName)
			defer ctrl.Finish()

			tc.expectations(t, &tc, driver)

			subject := New(&config.Config{})
			actualCtx, actualErr := subject.Connect(context.Background(), fmt.Sprintf(tc.dsn, driverName))

			if !reflect.DeepEqual(actualErr, tc.err) {
				t.Fatalf("\nexpected err: `%+v`\ngot:          `%+v`", tc.err, actualErr)
			}

			if tc.conn != nil && !reflect.DeepEqual(subject.connection(actualCtx), tc.conn) {
				t.Fatalf("\nexpected conn: `%+v`\ngot:           `%+v`", tc.conn, subject.connection(actualCtx))
			}
		})
	}
}

func TestEngineDatabases(t *testing.T) {
	type testcase struct {
		name         string
		dbs          []string
		expectations func(*testing.T, *testcase, *sqlenginetest.MockConnection)
		err          error
	}

	testcases := []testcase{
		{
			name: "databases_error",
			expectations: func(t *testing.T, tc *testcase, c *sqlenginetest.MockConnection) {
				c.EXPECT().
					Databases().
					Return(nil, tc.err)
			},
			err: errors.New("err"),
		},
		{
			name: "ok",
			dbs:  []string{"db1", "db2"},
			expectations: func(t *testing.T, tc *testcase, c *sqlenginetest.MockConnection) {
				db1 := sqlenginetest.NewMockDatabase(gomock.NewController(t))
				db2 := sqlenginetest.NewMockDatabase(gomock.NewController(t))
				db1.EXPECT().Name().Return("db1")
				db2.EXPECT().Name().Return("db2")

				c.EXPECT().
					Databases().
					Return([]driver.Database{db1, db2}, tc.err)
			},
			err: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			subject := New(&config.Config{})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			conn := sqlenginetest.NewMockConnection(ctrl)
			tc.expectations(t, &tc, conn)
			ctx := NewContext(nil, conn)

			actualDBS, actualErr := subject.Databases(ctx)

			if !reflect.DeepEqual(actualErr, tc.err) {
				t.Fatalf("\nexpected err: `%+v`\ngot:          `%+v`", tc.err, actualErr)
			}

			if tc.dbs != nil && !reflect.DeepEqual(actualDBS, tc.dbs) {
				t.Fatalf("\nexpected dbs: `%+v`\ngot:          `%+v`", tc.dbs, actualDBS)
			}
		})
	}
}

func TestEngineFetchTable(t *testing.T) {
	type testcase struct {
		name         string
		dbs          []string
		cols         []string
		rows         [][]string
		expectations func(*testing.T, *testcase, *sqlenginetest.MockConnection)
		err          error
	}

	testcases := []testcase{
		{
			name: "ok",
			dbs:  []string{"db1", "db2"},
			cols: []string{"a", "b"},
			rows: [][]string{{"a1", "b1"}, {"a2", "b2"}},
			expectations: func(t *testing.T, tc *testcase, conn *sqlenginetest.MockConnection) {
				db1 := sqlenginetest.NewMockDatabase(gomock.NewController(t))
				for _, db := range tc.dbs {
					db1.EXPECT().Name().AnyTimes().Return(db)
				}

				conn.EXPECT().
					Databases().MinTimes(2). // Databases direct call and UseDatabase implicit one
					Return([]driver.Database{db1}, tc.err)

				conn.EXPECT().
					UseDatabase("db1").
					Return(nil)

				db1.EXPECT().
					FetchTable(int64(0), int64(10)).
					Return(tc.cols, tc.rows, nil)
			},
			err: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			subject := New(&config.Config{})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			conn := sqlenginetest.NewMockConnection(ctrl)
			tc.expectations(t, &tc, conn)
			ctx := NewContext(nil, conn)

			dbs, err := subject.Databases(ctx)
			if err != nil {
				t.Fatal("unexpected set up error", err)
			}

			ctx, err = subject.UseDatabase(ctx, dbs[0])
			if err != nil {
				t.Fatal("unexpected set up error", err)
			}

			actualCols, actualRows, actualErr := subject.FetchTable(ctx, 0, 10)

			if !reflect.DeepEqual(actualErr, tc.err) {
				t.Fatalf("\nexpected err: `%+v`\ngot:          `%+v`", tc.err, actualErr)
			}

			if tc.cols != nil && !reflect.DeepEqual(actualCols, tc.cols) {
				t.Fatalf("\nexpected dbs: `%+v`\ngot:          `%+v`", tc.cols, actualCols)
			}

			if tc.rows != nil && !reflect.DeepEqual(actualRows, tc.rows) {
				t.Fatalf("\nexpected dbs: `%+v`\ngot:          `%+v`", tc.rows, actualRows)
			}
		})
	}
}
