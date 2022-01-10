package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Sample struct {
	ID     int64  `json:"id"`
	Foo    string `json:"foo"`
	IntVal int64  `json:"int_val"`
}

func CreateSample(dbConn *sql.DB, sample *Sample) (int64, error) {

	q := `INSERT INTO sample (foo, int_val) VALUES (?, ?)`
	id, err := insert(dbConn, q, sample.Foo, sample.IntVal)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetSample(dbConn *sql.DB, id int) (*Sample, error) {
	data := &Sample{}

	q := `SELECT id, foo, int_val FROM sample WHERE id = ?`
	err := dbConn.QueryRow(q, id).Scan(&data.ID, &data.Foo, &data.IntVal)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetManySample(dbConn *sql.DB) ([]*Sample, error) {
	res := []*Sample{}
	data := &Sample{}

	q := `SELECT id, foo, int_val FROM sample`
	rows, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.Foo, &data.IntVal)
		if err != nil {
			return nil, err
		}

		res = append(res, data)
	}

	return res, nil
}

func UpdateSample(dbConn *sql.DB, id int, sample *Sample) (int64, error) {
	args := make([]interface{}, 0, 3)

	q := `UPDATE sample SET id = id`
	if sample.Foo != "" {
		q += `, foo = ?`
		args = append(args, sample.Foo)
	}
	if sample.IntVal != 0 {
		q += `, int_val = ?`
		args = append(args, sample.IntVal)
	}
	q += ` WHERE id = ?`
	args = append(args, id)

	rowsAffected, err := update(dbConn, q, args)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func DeleteSample(id int) error {
	return nil
}
