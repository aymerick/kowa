package core

import "time"

// ValidTZ returns true if given timezone name is a valid IANA timezone.
func ValidTZ(name string) bool {
	_, err := time.LoadLocation(name)
	return err == nil
}
