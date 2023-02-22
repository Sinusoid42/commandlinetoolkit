package commandlinetoolkit

type settings struct {
	version     string
	author      string
	executeable string
	description string

	man string
}

func (s *settings) build(m map[string]interface{}) {

	s.author = check(AUTHORKEY, "author", m)

	s.version = check(VERSIONKEY, "0.0.1", m)

	s.man = check(MANUALKEY, "Man Page", m)

	s.description = check(DESCRIPTIONKEY, "Description", m)

	s.executeable = check(EXECUTEABLEKEY, "Executeable", m)

}
