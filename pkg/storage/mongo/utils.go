package mongo

import (
	"regexp"
)

// maskDSN masks the sensitive username and password in the mongo dsn
func maskDSN(dsn string) string {
	re := regexp.MustCompile(`//.*@`)
	return re.ReplaceAllString(dsn, "//*****:*****@")
}
