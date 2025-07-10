package main

import (
	"regexp"
)

var (
	REGEX_DOMAIN                  = `[\w-]+([-\.]{1}[\w-]+)*\.\w{2,63}`
	REGEX_EMAIL_PATTERN            = `^[\w\-\.]{1,255}@(` + REGEX_DOMAIN + `)$`
	WHITELISTED_EMAILS            = []string{}
)

var (
	compiledRegexEmailPattern           = regexp.MustCompile(REGEX_EMAIL_PATTERN)
)