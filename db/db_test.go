package db

import (
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grantsavage/ip-lookup-api/graph/model"
	uuid "github.com/satori/go.uuid"
)

func TestSetupDatbase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("should setup address_results table", func(t *testing.T) {
		mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS address_results(.+)").WillReturnError(nil)
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS address_results(.+)").WillReturnResult(sqlmock.NewResult(0, 0))

		err = SetupDatabase(db)
		if err != nil {
			t.Fatalf("error: '%s'", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}
	})

	t.Run("should return error on SQL execution error", func(t *testing.T) {
		executionError := errors.New("sql error")

		mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS address_results(.+)").WillReturnError(executionError)

		err = SetupDatabase(db)
		if err != executionError {
			t.Errorf("got error '%s', wanted '%s'", err, executionError)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}
	})
}

func TestGetIPLookupResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("should return lookup result", func(t *testing.T) {
		ip := net.ParseIP("1.2.3.4")
		result := &model.IPLookupResult{
			UUID:         uuid.NewV4().String(),
			IPAddress:    ip.String(),
			ResponseCode: "127.0.0.4",
			CreatedAt:    time.Now().Format(time.RFC3339),
			UpdatedAt:    time.Now().Format(time.RFC3339),
		}

		rows := sqlmock.
			NewRows([]string{"uuid", "ip_address", "response_code", "created_at", "updated_at"}).
			AddRow(
				result.UUID,
				result.IPAddress,
				result.ResponseCode,
				result.CreatedAt,
				result.UpdatedAt,
			)
		mock.
			ExpectQuery(`SELECT(.+)FROM address_results(.+)`).
			WithArgs(ip.String()).
			WillReturnRows(rows)

		lookupResult, err := GetIPLookupResult(db, ip)
		if err != nil {
			t.Fatalf("error: '%s'", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}

		if !reflect.DeepEqual(result, lookupResult) {
			t.Errorf("got '%q', want '%q", lookupResult, result)
		}
	})

	t.Run("should return error if no lookup result is found", func(t *testing.T) {
		ip := net.ParseIP("5.6.7.8")

		rows := sqlmock.
			NewRows([]string{"uuid", "ip_address", "response_code", "created_at", "updated_at"})
		mock.
			ExpectQuery(`SELECT(.+)FROM address_results(.+)`).
			WithArgs(ip.String()).
			WillReturnRows(rows)

		_, err := GetIPLookupResult(db, ip)
		if err == nil {
			t.Error("error not returned")
		}

		if err != ErrorNotFound {
			t.Errorf("got '%s', want %s", err, ErrorNotFound)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		ip := net.ParseIP("5.6.7.8")

		queryError := errors.New("unable to query")
		mock.
			ExpectQuery(`SELECT(.+)FROM address_results(.+)`).
			WithArgs(ip.String()).
			WillReturnError(queryError)

		_, err := GetIPLookupResult(db, ip)
		if err != queryError {
			t.Errorf("got error '%s', wanted '%s'", err, queryError)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}
	})
}

func TestUpsertIPLookupResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("should upsert result", func(t *testing.T) {
		ip := net.ParseIP("1.2.3.4")
		result := model.IPLookupResult{
			UUID:         uuid.NewV4().String(),
			IPAddress:    ip.String(),
			ResponseCode: "127.0.0.4",
			CreatedAt:    time.Now().Format(time.RFC3339),
			UpdatedAt:    time.Now().Format(time.RFC3339),
		}

		mock.
			ExpectPrepare(`INSERT INTO address_results(.+)ON CONFLICT(.+)`).
			WillReturnError(nil)
		mock.
			ExpectExec(`INSERT INTO address_results(.+)ON CONFLICT(.+)`).
			WithArgs(
				result.UUID,
				result.IPAddress,
				result.ResponseCode,
				result.CreatedAt,
				result.UpdatedAt,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err = UpsertIPLookupResult(db, result)
		if err != nil {
			t.Fatalf("error: '%s'", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}
	})

	t.Run("should return error when SQL exception occurs", func(t *testing.T) {
		ip := net.ParseIP("1.2.3.4")
		result := model.IPLookupResult{
			UUID:         uuid.NewV4().String(),
			IPAddress:    ip.String(),
			ResponseCode: "127.0.0.4",
			CreatedAt:    time.Now().Format(time.RFC3339),
			UpdatedAt:    time.Now().Format(time.RFC3339),
		}

		executionError := errors.New("sql error")
		mock.
			ExpectPrepare(`INSERT INTO address_results(.+)ON CONFLICT(.+)`).
			WillReturnError(executionError)

		err = UpsertIPLookupResult(db, result)
		if err != executionError {
			t.Errorf("got error '%s', wanted '%s'", err, executionError)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("expectations were not met: '%s'", err)
		}
	})
}
