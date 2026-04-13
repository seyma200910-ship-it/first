package pkg

import "errors"

var (
	ErrUserEmailExists = errors.New("такой email уже существет")
	ErrUserNotFound    = errors.New("такой юзер не найден")
)
