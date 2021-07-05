package util

import (
	"regexp"
	"sync"
)

var (
	_Mail *regexp.Regexp
	regOnce sync.Once
)

// RegMail follow the singleton pattern
func RegMail() *regexp.Regexp {
	regOnce.Do(initRegMail)
	return _Mail
}

func initRegMail() {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	_Mail = regexp.MustCompile(pattern)
}
