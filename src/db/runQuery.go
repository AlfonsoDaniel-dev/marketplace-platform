package db

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

func scanQuery(rows *sql.Rows) ([]any, error) {

	var Results []any

	for rows.Next() {
		var result any
		err := rows.Scan(&result)
		if err != nil {
			return nil, err
		}

		Results = append(Results, result)
	}

	return Results, nil
}

func RunQuery(tx *sql.Tx, query string, params ...any) ([]any, error) {
	if query == "" || len(params) == 0 {
		return nil, errors.New("empty query or params")
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	queryParams := make([]interface{}, len(params))
	for i, queryParam := range params {
		queryParams[i] = queryParam
	}

	// pasa los elementos del slice como si fueran parametros individuales
	rows, err := stmt.Query(queryParams...) // sintaxis de elipse ...
	if err != nil {
		return nil, err
	}

	results, err := scanQuery(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func ExecQuery(tx *sql.Tx, query string, params ...any) (int64, error) {
	if query == "" {
		return 0, errors.New("empty query or params")
	} else if len(params) == 0 {
		stmt, err := tx.Prepare(query)
		if err != nil {
			return 0, err
		}

		defer stmt.Close()

		res, err := stmt.Exec()
		rowsAff, err := res.RowsAffected()
		if err != nil {
			return rowsAff, err
		}

		return rowsAff, nil
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	queryParams := make([]interface{}, len(params))
	for i, queryParam := range params {
		queryParams[i] = queryParam
	}

	res, err := stmt.Exec(queryParams...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func ParseAnyToString(str any) (string, error) {
	if _, ok := str.(string); !ok {
		return "", errors.New("string expected")
	}

	return str.(string), nil
}

func ParseAnyToUUID(uid any) (uuid.UUID, error) {
	if _, ok := uid.(uuid.UUID); !ok {
		return uuid.Nil, errors.New("uuid expected")
	}

	return uid.(uuid.UUID), nil
}

func ParseAnyToInt64(number any) (int64, error) {
	if _, ok := number.(int64); !ok {
		return 0, errors.New("int64 expected")
	}

	return number.(int64), nil
}

func ParseAnyToInt(number any) (int, error) {
	if _, ok := number.(int); !ok {
		return 0, errors.New("int expected")
	}

	return number.(int), nil
}

func ParseAnytToFloat(flt any) (float64, error) {
	if _, ok := flt.(string); !ok {
		return 0.0, errors.New("float expected")
	}

	return flt.(float64), nil
}

func ParseAnyToBool(value any) (bool, error) {
	if _, ok := value.(bool); !ok {
		return false, errors.New("bool expected")
	}

	return value.(bool), nil
}

func ParseAnyToTime(date any) (time.Time, error) {
	if _, ok := date.(time.Time); !ok {
		return time.Time{}, errors.New("time expected")
	}

	return date.(time.Time), nil
}

func ParseAnyToBytes(data any) ([]byte, error) {
	if _, ok := data.([]byte); !ok {
		return nil, errors.New("bytes expected")
	}

	return data.([]byte), nil
}
