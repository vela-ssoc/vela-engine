package match

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
)

type Condition struct {
	Logic    string     `yaml:"logic"`
	Matchers []*Matcher `yaml:"matchers"`
	cnd      *cond.Cond `yaml:"-"`
}

func (c *Condition) Convert() (*cond.Cond, error) {
	cnd := &cond.Cond{}
	if len(c.Matchers) == 0 {
		return cnd, nil
	}

	for i, match := range c.Matchers {
		if s := match.section(); s == nil {
			return nil, fmt.Errorf("#%d match section convert to cond fail", i)
		} else {
			cnd.Append(s)
		}
	}
	return cnd, nil
}

func (c *Condition) Prepare() error {
	cnd, err := c.Convert()
	if err != nil {
		return err
	}
	c.cnd = cnd
	return nil
}

func (c *Condition) Match(v interface{}, opt ...cond.OptionFunc) bool { //ctx
	if c.cnd == nil {
		return false
	}

	switch c.Logic {
	case "and":
		opt = append(opt, cond.WithLogic(cond.AND))
	case "or":
		opt = append(opt, cond.WithLogic(cond.OR))
	default:
		opt = append(opt, cond.WithLogic(cond.AND))
	}
	return c.cnd.Match(v, opt...)
}
