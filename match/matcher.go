package match

import cond "github.com/vela-ssoc/vela-cond"

const (
	And Logic = "and"
	Or  Logic = "or"
)

type Logic string

type Matcher struct {
	Type   string   `yaml:"type"`
	Method string   `yaml:"method"`
	Part   int      `yaml:"part"`
	Value  []string `yaml:"value"`
}

func (m *Matcher) section() *cond.Section {
	section := &cond.Section{}
	section.Keys(m.Type)
	method := m.Method
	if len(method) < 2 {
		return nil
	}

	if method[0] == '!' {
		method = method[1:]
		section.WithNot()
	}

	switch m.Method {
	case "regex":
		section.Method(cond.Re)
		section.Regex(m.Value...)
		section.Partition(m.Part)

	case "equal":
		section.Method(cond.Eq)
		section.Value(m.Value...)
	case "word":
		section.Method(cond.Cn)
		section.Value(m.Value...)
	case "lt":
		section.Method(cond.Lt)
		section.Value(m.Value...)
	case "le":
		section.Method(cond.Le)
		section.Value(m.Value...)
	case "gt":
		section.Method(cond.Gt)
		section.Value(m.Value...)
	case "ge":
		section.Method(cond.Ge)
		section.Value(m.Value...)

	case "call":
		section.Method(cond.Call)
		section.Value(m.Value...)

	default:
		return nil
	}

	return section
}
