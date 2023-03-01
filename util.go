/*
//  Implementation of a few utilitary functions:
//    - getSeconds
*/

package main

// GetSeconds splits a time in milliseconds into seconds and milliseconds
func GetSeconds(d int64) (s, ms int64) {
	s = d / 1000
	ms = d - s*1000
	return s, ms
}
