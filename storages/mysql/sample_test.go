package mysql

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestCreateSample(t *testing.T) {
	type in struct {
		Sample *Sample
	}
	type out struct {
		Expected int64
		Error    error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case 1",
			&in{
				Sample: &Sample{
					Foo:    "var",
					IntVal: int64(1),
				},
			},
			&out{
				Expected: int64(1),
			},
		},
		{
			"success case 2",
			&in{
				Sample: &Sample{
					Foo:    "var 2",
					IntVal: int64(2),
				},
			},
			&out{
				Expected: int64(2),
			},
		},
		{
			"failure case, invalid data sent",
			&in{
				Sample: nil,
			},
			&out{
				Expected: 0,
				Error:    errors.New("invalid data sent"),
			},
		},
	}

	createTable("sample")
	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		id, err := CreateSample(testDBConn, in.Sample)
		if testCase.Out.Expected != id {
			t.Errorf("test failed, got: %v, want: %v", id, testCase.Out.Expected)
		}

		switch {
		case err != nil && out.Error == nil:
			t.Errorf("expected non error, but some error occurred, %s", err.Error())
		case err == nil && out.Error != nil:
			t.Errorf("expected error %s, but results: no error", out.Error.Error())
		case err != nil && out.Error != nil:
			// fmt.Println(err)
		}

	}
	deleteTable("sample")
}

func TestGetSample(t *testing.T) {
	testData := &Sample{
		ID:     int64(1),
		Foo:    "var",
		IntVal: int64(100),
	}

	type in struct {
		ID int64
	}
	type out struct {
		Expected *Sample
		Error    error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case",
			&in{
				ID: int64(1),
			},
			&out{
				Expected: testData,
			},
		},
		{
			"failure case, data not exits",
			&in{
				ID: int64(2),
			},
			&out{
				Expected: nil,
				Error:    errors.New("sql: no rows in result set"),
			},
		},
	}

	createTable("sample")
	CreateSample(testDBConn, testData)
	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		got, err := GetSample(testDBConn, in.ID)
		if !reflect.DeepEqual(out.Expected, got) {
			t.Errorf("test failed, got: %v, want: %v", got, testCase.Out.Expected)
		}

		switch {
		case err != nil && out.Error == nil:
			t.Errorf("expected non error, but some error occurred, %s", err.Error())
		case err == nil && out.Error != nil:
			t.Errorf("expected error %s, but results: no error", out.Error.Error())
		case err != nil && out.Error != nil:
			// fmt.Println(err)
		}
	}
	deleteTable("sample")
}

func TestGetManySample(t *testing.T) {
	testDataList := []*Sample{
		{
			ID:     int64(1),
			Foo:    "var",
			IntVal: int64(101),
		},
		{
			ID:     int64(2),
			Foo:    "var2",
			IntVal: int64(102),
		},
	}

	type out struct {
		Expected []*Sample
		Error    error
	}
	type testCase struct {
		Scenario string
		Out      *out
	}

	testCases := []testCase{
		{
			"success case",
			&out{
				Expected: testDataList,
			},
		},
	}

	createTable("sample")
	for _, d := range testDataList {
		CreateSample(testDBConn, d)
	}
	for _, testCase := range testCases {
		out := testCase.Out
		got, _ := GetManySample(testDBConn)
		if !reflect.DeepEqual(out.Expected, got) {
			gotStr, _ := json.Marshal(got)
			expectedStr, _ := json.Marshal(out.Expected)
			t.Errorf("test failed, got: %v, want: %v", string(gotStr), string(expectedStr))
		}
	}
	deleteTable("sample")
}

func TestUpdateSample(t *testing.T) {
	testData := &Sample{
		ID:     int64(1),
		Foo:    "var",
		IntVal: int64(100),
	}

	type in struct {
		ID   int64
		Data *Sample
	}
	type out struct {
		RowsAffected int64
		Expected     *Sample
		Error        error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case, update one field",
			&in{
				ID: int64(1),
				Data: &Sample{
					Foo: "var mod 1",
				},
			},
			&out{
				RowsAffected: int64(1),
				Expected: &Sample{
					ID:     int64(1),
					Foo:    "var mod 1",
					IntVal: int64(100),
				},
			},
		},
		{
			"success case, update all field",
			&in{
				ID: int64(1),
				Data: &Sample{
					Foo:    "var mod 2",
					IntVal: int64(101),
				},
			},
			&out{
				RowsAffected: int64(1),
				Expected: &Sample{
					ID:     int64(1),
					Foo:    "var mod 2",
					IntVal: int64(101),
				},
			},
		},
		{
			"success case, update as same record",
			&in{
				ID: int64(1),
				Data: &Sample{
					Foo:    "var mod 2",
					IntVal: int64(101),
				},
			},
			&out{
				RowsAffected: int64(0),
				Expected: &Sample{
					ID:     int64(1),
					Foo:    "var mod 2",
					IntVal: int64(101),
				},
			},
		},
		{
			"success case, data not exits",
			&in{
				ID: int64(2),
				Data: &Sample{
					Foo:    "var mod 3",
					IntVal: int64(102),
				},
			},
			&out{
				RowsAffected: int64(0),
				Expected:     nil,
			},
		},
		{
			"failure case, invalid data sent",
			&in{
				ID:   int64(1),
				Data: nil,
			},
			&out{
				RowsAffected: int64(0),
				// the original data shouldn't be changed
				Expected: &Sample{
					ID:     int64(1),
					Foo:    "var mod 2",
					IntVal: int64(101),
				},
				Error: errors.New("invalid data"),
			},
		},
	}

	createTable("sample")
	CreateSample(testDBConn, testData)
	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		got, err := UpdateSample(testDBConn, in.ID, in.Data)
		if out.RowsAffected != got {
			t.Errorf("test failed, got: %v, want: %v", got, testCase.Out.RowsAffected)
		}

		updated, _ := GetSample(testDBConn, in.ID)
		if !reflect.DeepEqual(out.Expected, updated) {
			t.Errorf("test failed, got: %v, want: %v", updated, testCase.Out.Expected)
		}

		switch {
		case err != nil && out.Error == nil:
			t.Errorf("expected non error, but some error occurred, %s", err.Error())
		case err == nil && out.Error != nil:
			t.Errorf("expected error %s, but results: no error", out.Error.Error())
		case err != nil && out.Error != nil:
			// fmt.Println(err)
		}
	}
	deleteTable("sample")
}
