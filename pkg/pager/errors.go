package pager

import "errors"

var (
	ErrBadAlignmentSize        = errors.New("pageManagerFile: bad alignment size")
	ErrNoMoreRoomInPage        = errors.New("page: there is not enough room left in the page")
	ErrInvalidRecordID         = errors.New("page: invalid record id")
	ErrRecordHasBeenMarkedFree = errors.New("page: record has been marked free (aka, removed)")
	ErrRecordNotFound          = errors.New("page: record could not be found")
	ErrPageNotFound            = errors.New("pageManagerFile: page could not be found")
	ErrWritingPage             = errors.New("pageManagerFile: error writing page")
	ErrDeletingPage            = errors.New("pageManagerFile: error deleting page")
	ErrMinRecordSize           = errors.New("page: record is smaller than the min record size allowed")
	ErrMaxRecordSize           = errors.New("page: record is larger than the max record size allowed")
	ErrRecordMaxKeySize        = errors.New("record: record key is longer than max size allowed (255)")
	ErrPageIsNotOverflow       = errors.New("pagemanager: error page is not an overflow page")
)
