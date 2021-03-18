package main

type UserLegacy struct {
	name        string
	permissions map[Permission]map[string]bool
	identity    Identity
}

type Identity int

const (
	superAdmin Identity = 0
	admin      Identity = 1
	player     Identity = 2
	banned     Identity = 99
)
