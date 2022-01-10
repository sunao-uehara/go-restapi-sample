package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func Initialize(url string) (*sql.DB, error) {
	sqlURL, err := mysql.ParseDSN(url)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", sqlURL.FormatDSN())
	if err != nil {
		return nil, err
	}
	// defer db.Close()

	return db, nil
}

func insert(dbConn *sql.DB, sql string, args ...interface{}) (int64, error) {
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func update(dbConn *sql.DB, sql string, args []interface{}) (int64, error) {
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// func getOne(dbConn *sql.DB, sql string, data *Sample, args ...interface{}) (*Sample, error) {
// 	err := dbConn.QueryRow(sql, args).Scan(data)
// 	// err := row.Scan(&data.ID, &data.Foo, &data.IntVal)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
