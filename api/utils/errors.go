package utils

import "codis/utils/slogger"

func LogErrorOrNothing(errors ...error) {
	for _, err := range errors {
		if err != nil {
			slogger.Error(err)
		}
	}
}
