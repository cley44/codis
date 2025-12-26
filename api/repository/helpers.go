package repository

import (
	"database/sql"
	"errors"
)

func notFoundIsNotAnError(e error) (exist bool, err error) {
	if e == nil {
		return true, nil
	} else if errors.Is(e, sql.ErrNoRows) {
		return false, nil
	}

	return false, e
}
