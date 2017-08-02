package rrframework

import (
	"fmt"
)

type framework struct {
	ver   string
	gover string
}

var (
	Framework = &framework{
		ver:   "0.1",
		gover: "1.5+",
	}
)

func (s *framework) Description() string {
	return fmt.Sprintf("rrframework verion %s, build on go %s", s.ver, s.gover)
}
