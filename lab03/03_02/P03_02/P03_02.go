package P03_02

import "fmt"

type Stats struct {
	postCount int
	getCount  int
}

func (s *Stats) IncGet() {
	s.getCount++
}

func (s *Stats) IncPost() {
	s.postCount++
}

func (s *Stats) GenStr() string {
	return fmt.Sprintf("GET count = %d\nPOST count %d", s.getCount, s.postCount)
}
