package scope_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/goaux/scope"
)

func ExampleTx() {
	for db, err := range scope.Use2(sql.Open("sql_test", "ExampleOpenDB")) {
		if err != nil {
			fmt.Println(err)
			break
		}
		tx, err := db.Begin()
		if err != nil {
			fmt.Println(err)
			break
		}
		for tx := range scope.Tx(tx) {
			// do something.
			_ = tx
			// Not committing implies a rollback
		}
		fmt.Println("end#1")
	}
	fmt.Println("end#2")
	// Output:
	// Driver.OpenConnector("ExampleOpenDB")
	// Connector.Connect(ctx)
	// Conn.BeginTx(ctx, {"Isolation":0,"ReadOnly":false})
	// Tx.Rollback()
	// end#1
	// Conn.Close()
	// end#2
}

func ExampleTx2_begin() {
	for db, err := range scope.Use2(sql.Open("sql_test", "ExampleOpenDB")) {
		if err != nil {
			fmt.Println(err)
			break
		}
		for tx, err := range scope.Tx2(db.Begin()) {
			if err != nil {
				fmt.Println(err)
				break
			}

			// do something.
			_ = tx

			// Under normal circumstances, Commit is called last.
			if err := tx.Commit(); err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("end#1")
		}

		for tx, err := range scope.Tx2(db.Begin()) {
			if err != nil {
				fmt.Println(err)
				break
			}

			// do something.
			_ = tx

			// Not committing implies a rollback
			fmt.Println("end#2")
		}
		fmt.Println("end#3")
	}
	fmt.Println("end#4")
	// Output:
	// Driver.OpenConnector("ExampleOpenDB")
	// Connector.Connect(ctx)
	// Conn.BeginTx(ctx, {"Isolation":0,"ReadOnly":false})
	// Tx.Commit()
	// end#1
	// Conn.BeginTx(ctx, {"Isolation":0,"ReadOnly":false})
	// end#2
	// Tx.Rollback()
	// end#3
	// Conn.Close()
	// end#4
}

func ExampleTx2_beginTx() {
	for db, err := range scope.Use2(sql.Open("sql_test", "ExampleOpenDB")) {
		if err != nil {
			fmt.Println(err)
			break
		}
		ctx := context.TODO()
		for tx, err := range scope.Tx2(db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})) {
			if err != nil {
				fmt.Println(err)
				break
			}

			// do something.
			_ = tx

			if err := tx.Commit(); err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("end#1")
		}

		for tx, err := range scope.Tx2(db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})) {
			if err != nil {
				fmt.Println(err)
				break
			}

			// do something.
			_ = tx

			// Not committing implies a rollback
			fmt.Println("end#2")
		}
		fmt.Println("end#3")
	}
	fmt.Println("end#4")
	// Output:
	// Driver.OpenConnector("ExampleOpenDB")
	// Connector.Connect(ctx)
	// Conn.BeginTx(ctx, {"Isolation":6,"ReadOnly":false})
	// Tx.Commit()
	// end#1
	// Conn.BeginTx(ctx, {"Isolation":6,"ReadOnly":false})
	// end#2
	// Tx.Rollback()
	// end#3
	// Conn.Close()
	// end#4
}

func init() {
	sql.Register("sql_test", Driver{})
}

type Driver struct{}

var _ driver.Driver = Driver{}
var _ driver.DriverContext = Driver{}

func (Driver) Open(name string) (driver.Conn, error) {
	fmt.Printf("Driver.Open(%q)\n", name)
	return Conn{}, nil
}

func (Driver) OpenConnector(name string) (driver.Connector, error) {
	fmt.Printf("Driver.OpenConnector(%q)\n", name)
	return Connector{}, nil
}

type Connector struct{}

var _ driver.Connector = Connector{}

func (Connector) Connect(context.Context) (driver.Conn, error) {
	fmt.Println("Connector.Connect(ctx)")
	return Conn{}, nil
}

func (Connector) Driver() driver.Driver {
	return Driver{}
}

type Conn struct{}

var _ driver.Conn = Conn{}
var _ driver.ConnBeginTx = Conn{}

func (Conn) Prepare(query string) (driver.Stmt, error) {
	fmt.Printf("Conn.Prepare(%q)\n", query)
	return Stmt{}, nil
}

func (Conn) Close() error {
	fmt.Println("Conn.Close()")
	return nil
}

func (Conn) Begin() (driver.Tx, error) {
	fmt.Println("Conn.Begin()")
	return Tx{}, nil
}

func (Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	p, _ := json.Marshal(opts)
	fmt.Printf("Conn.BeginTx(ctx, %s)\n", string(p))
	return Tx{}, nil
}

type Stmt struct{}

var _ driver.Stmt = Stmt{}

func (Stmt) Close() error {
	fmt.Println("Stmt.Close()")
	return nil
}

func (Stmt) NumInput() int {
	fmt.Println("Stmt.NumInput()")
	return 0
}

func (Stmt) Exec(args []driver.Value) (driver.Result, error) {
	p, _ := json.Marshal(args)
	fmt.Printf("Stmt.Exec(%s)", string(p))
	return Result{}, nil
}

func (Stmt) Query(args []driver.Value) (driver.Rows, error) {
	p, _ := json.Marshal(args)
	fmt.Printf("Stmt.Query(%s)", string(p))
	return Rows{}, nil
}

type Tx struct{}

var _ driver.Tx = Tx{}

func (Tx) Commit() error {
	fmt.Println("Tx.Commit()")
	return nil
}

func (Tx) Rollback() error {
	fmt.Println("Tx.Rollback()")
	return nil
}

type Result struct{}

var _ driver.Result = Result{}

func (Result) LastInsertId() (int64, error) {
	fmt.Println("Result.LastInsertId()")
	return 0, nil
}

func (Result) RowsAffected() (int64, error) {
	fmt.Println("Result.RowsAffected()")
	return 0, nil
}

type Rows struct{}

func (Rows) Columns() []string {
	fmt.Println("Rows.Columns()")
	return nil
}

// Close closes the rows iterator.
func (Rows) Close() error {
	fmt.Println("Rows.Close()")
	return nil
}

func (Rows) Next(dest []driver.Value) error {
	p, _ := json.Marshal(dest)
	fmt.Printf("Rows.Next(%s)\n", string(p))
	return nil
}
