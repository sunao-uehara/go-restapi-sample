package mysql

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Sample interface {
	CreateSample(sample *SampleData) (int64, error)
	GetSample(id int64) (*SampleData, error)
	GetManySample() ([]*SampleData, error)
	UpdateSample(int64, *SampleData) (int64, error)
}

func NewSample(dbConn *sql.DB) Sample {
	return &SQLSample{
		db: dbConn,
	}
}

type SQLSample struct {
	db *sql.DB
}

// SampleData is data structure that is corresponding to the table `sample`
type SampleData struct {
	ID     int64  `json:"id"`
	Foo    string `json:"foo"`
	IntVal int64  `json:"int_val"`
}

func (sc *SQLSample) CreateSample(sample *SampleData) (int64, error) {
	if sample == nil {
		return 0, errors.New("invalid data")
	}

	q := `INSERT INTO sample (foo, int_val) VALUES (?, ?)`
	id, err := insert(sc.db, q, sample.Foo, sample.IntVal)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (sc *SQLSample) GetSample(id int64) (*SampleData, error) {
	data := &SampleData{}

	q := `SELECT id, foo, int_val FROM sample WHERE id = ?`
	err := sc.db.QueryRow(q, id).Scan(&data.ID, &data.Foo, &data.IntVal)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sc *SQLSample) GetManySample() ([]*SampleData, error) {
	q := `SELECT id, foo, int_val FROM sample ORDER BY ID ASC`
	rows, err := sc.db.Query(q)
	if err != nil {
		return nil, err
	}

	res := []*SampleData{}
	for rows.Next() {
		data := &SampleData{}
		err := rows.Scan(&data.ID, &data.Foo, &data.IntVal)
		if err != nil {
			return nil, err
		}

		res = append(res, data)
	}

	return res, nil
}

func (sc *SQLSample) UpdateSample(id int64, sample *SampleData) (int64, error) {
	args := make([]interface{}, 0, 3)

	if sample == nil {
		return 0, errors.New("invalid data")
	}

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

	rowsAffected, err := update(sc.db, q, args)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
