package excel

import (
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeExcelOpenFile  = "XLS-001"
	ErrCodeExcelOpenSheet = "XLS-002"
	ErrCodeExcelReadRows  = "XLS-003"
)

var (
	ErrExcelOpenFile = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeExcelOpenFile, "").Wrap(cause).Err()
	}
	ErrExcelOpenSheet = func() error {
		return kit.NewAppErrBuilder(ErrCodeExcelOpenSheet, "unable to find excel sheet to open").Err()
	}
	ErrExcelReadRows = func(cause error) error {
		return kit.NewAppErrBuilder(ErrCodeExcelReadRows, "").Wrap(cause).Err()
	}
)
