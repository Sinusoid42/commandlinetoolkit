package commandlinetoolkit

type Key byte

var KEY_DELETE = Key(127)

var KEY_LF = Key(13)

var KEY_ESC = Key(27)

var KEY_RETURN = Key(10)

var KEY_SPACE = Key(32)

var KEY_TAB = Key(9)

var ARROW_UP = []Key{Key(27), Key(91), Key(65)}
var ARROW_DOWN = []Key{Key(27), Key(91), Key(66)}
var ARROW_LEFT = []Key{Key(27), Key(91), Key(68)}
var ARROW_RIGHT = []Key{Key(27), Key(91), Key(67)}

func ConvertKeysToString(keys []Key) string {
	s := ""
	for k := range keys {
		s += string(k)
	}
	return s
}
