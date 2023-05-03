package commandlinetoolkit

type settings struct {
	version     string
	author      string
	executeable string
	description string
	title       string

	_styleTitle bool

	customTitle bool
	_printTitle bool
	man         string
}

func (s *settings) build(m map[string]interface{}) {

	s.author = check(AUTHORKEY, "author", m)

	s.version = check(VERSIONKEY, "0.0.1", m)

	s.man = check(MANUALKEY, "Man Page", m)

	//s.title = check(TITLEKEY, "", m)

	if m[TITLEKEY] != nil && len(m[TITLEKEY].(string)) >= 0 {
		s.customTitle = true
		s.title = m[TITLEKEY].(string)
	} else {
		if len(s.title) <= 0 {
			s.title = "Command Line: " + s.version
		}
		s.customTitle = false
	}

	s.description = check(DESCRIPTIONKEY, "Description", m)

	s.executeable = check(EXECUTEABLEKEY, "Executeable", m)

	if v, ok := m[STYLETITLEKEY].(bool); ok {
		s._styleTitle = v
	} else {
	}
	if v, ok := m[PRINTTITLE].(bool); ok {
		s._printTitle = v
	} else {
	}
}

func (s *settings) clone() *settings {
	ns := &settings{
		version:     s.version,
		author:      s.author,
		executeable: s.executeable,
		description: s.description,
		man:         s.man,
		_styleTitle: s._styleTitle,
		customTitle: s.customTitle,
		title:       s.title,
		_printTitle: s._printTitle,
	}
	return ns
}
