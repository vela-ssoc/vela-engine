package header

import (
	"github.com/vela-ssoc/vela-kit/grep"
	"strings"
)

type Slice []string

func (s *Slice) UnmarshalYAML(fn func(interface{}) error) error {
	var v string
	if err := fn(&v); err != nil {
		return err
	}

	*s = append((*s)[:0], strings.Split(v, ",")...)
	return nil
}

func (s *Slice) String() string {
	return strings.Join(*s, ",")
}

func (s *Slice) have(tag string) bool {
	a := *s

	var have func(string) bool

	if strings.Index(tag, "*") == -1 {
		have = func(b string) bool { return tag == b }
	} else {
		have = grep.New(tag)
	}

	for i := 0; i < len(a); i++ {
		if have(a[i]) {
			return true
		}
	}
	return false

}

func (s *Slice) Have(v ...string) bool {
	a := *s
	n := len(a)
	if n == 0 {
		return false
	}

	if len(v) == 0 {
		return false
	}

	for _, tag := range v {
		if s.have(tag) {
			return true
		}
	}

	return false
}
