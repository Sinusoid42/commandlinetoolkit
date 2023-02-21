package commandlinetoolkit

type settings struct {
	version     string
	author      string
	executeable string
	description string

	man string
}

func (s *settings) build(m map[string]interface{}) {
	if str, ok := m["version"].(string); ok {

		s.version = str

	} else {

		s.version = "0.0.1"

	}

	s.author = check("author", "author", m)

	s.version = check("version", "0.0.1", m)

	s.man = check("man", "Man Page", m)

}
