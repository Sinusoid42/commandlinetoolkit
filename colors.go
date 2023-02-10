package commandlinetoolkit

type Color string

const COLOR_HEADER = string("\033[")

const REGULAR_COLORFONT = string("0")
const LIGHT_COLORFONT = string("1")
const DARK_COLORFONT = string("2")
const ITALIC_COLORFONT = string("3")
const UNDERLINE_COLORFONT = string("4")
const BLINKING_COLORFONT = string("5")

// @deprecated
// const BLINKING_FAST_FONT = string("6")

const REVERSE_COLORFONT = string("7")
const HIDE_COLORFONT = string("8")
const CROSS_OUT_COLORFONT = string("9")

const REGULAR_COLORTYPE = string("3")
const BACKGROUND_COLORTYPE = string("4")
const INTENSITY_COLORTYPE = string("9")
const INTENSITY_BACKGROUND_COLORTYPE = string("10")

const BLACK_COLOR = string("0")
const RED_COLOR = string("1")
const GREEN_COLOR = string("2")
const YELLOW_COLOR = string("3")
const BLUE_COLOR = string("4")
const PINK_COLOR = string("5")
const CYAN_COLOR = string("6")
const GRAY_COLOR = string("7")

func GenColor(fontType string, colType string, color string) Color {

	return Color(COLOR_HEADER + fontType + ";" + colType + color + "m")

}

// TODO
var COLOR_BLACK = Color(COLOR_HEADER + "0;30m")
var COLOR_BLACK_L = Color(COLOR_HEADER + "1;30m")
var COLOR_BLACK_D = Color(COLOR_HEADER + "2;30m")
var COLOR_BLACK_I = Color(COLOR_HEADER + "3;30m")
var COLOR_BLACK_U = Color(COLOR_HEADER + "4;30m")
var COLOR_BLACK_BL = Color(COLOR_HEADER + "5;30m")
var COLOR_BLACK_X = Color(COLOR_HEADER + "9;30m")
var COLOR_BLACK_BG = Color(COLOR_HEADER + "0;40m")
var COLOR_BLACK_IBG = Color(COLOR_HEADER + "0;100m")

var COLOR_RED = Color(COLOR_HEADER + "0;31m")
var COLOR_RED_L = Color(COLOR_HEADER + "1;31m")
var COLOR_RED_D = Color(COLOR_HEADER + "2;31m")
var COLOR_RED_I = Color(COLOR_HEADER + "3;31m")
var COLOR_RED_U = Color(COLOR_HEADER + "4;31m")
var COLOR_RED_BL = Color(COLOR_HEADER + "5;31m")
var COLOR_RED_X = Color(COLOR_HEADER + "9;31m")
var COLOR_RED_BG = Color(COLOR_HEADER + "0;41m")
var COLOR_RED_IBG = Color(COLOR_HEADER + "0;101m")

var COLOR_GREEN = Color(COLOR_HEADER + "0;32m")
var COLOR_GREEN_L = Color(COLOR_HEADER + "1;32m")
var COLOR_GREEN_D = Color(COLOR_HEADER + "2;32m")
var COLOR_GREEN_I = Color(COLOR_HEADER + "3;32m")
var COLOR_GREEN_U = Color(COLOR_HEADER + "4;32m")
var COLOR_GREEN_BL = Color(COLOR_HEADER + "5;32m")
var COLOR_GREEN_X = Color(COLOR_HEADER + "9;32m")
var COLOR_GREEN_BG = Color(COLOR_HEADER + "0;42m")
var COLOR_GREEN_IBG = Color(COLOR_HEADER + "0;102m")

var COLOR_YELLOW = Color(COLOR_HEADER + "0;33m")
var COLOR_YELLOW_L = Color(COLOR_HEADER + "1;33m")
var COLOR_YELLOW_D = Color(COLOR_HEADER + "2;33m")
var COLOR_YELLOW_I = Color(COLOR_HEADER + "3;33m")
var COLOR_YELLOW_U = Color(COLOR_HEADER + "4;33m")
var COLOR_YELLOW_BL = Color(COLOR_HEADER + "5;33m")
var COLOR_YELLOW_X = Color(COLOR_HEADER + "9;33m")
var COLOR_YELLOW_BG = Color(COLOR_HEADER + "0;43m")
var COLOR_YELLOW_IBG = Color(COLOR_HEADER + "0;103m")

var COLOR_BLUE = Color(COLOR_HEADER + "0;34m")
var COLOR_BLUE_L = Color(COLOR_HEADER + "1;34m")
var COLOR_BLUE_D = Color(COLOR_HEADER + "2;34m")
var COLOR_BLUE_I = Color(COLOR_HEADER + "3;34m")
var COLOR_BLUE_U = Color(COLOR_HEADER + "4;34m")
var COLOR_BLUE_BL = Color(COLOR_HEADER + "5;34m")
var COLOR_BLUE_X = Color(COLOR_HEADER + "9;34m")
var COLOR_BLUE_BG = Color(COLOR_HEADER + "0;44m")
var COLOR_BLUE_IBG = Color(COLOR_HEADER + "0;104m")

var COLOR_PINK = Color(COLOR_HEADER + "0;35m")
var COLOR_PINK_L = Color(COLOR_HEADER + "1;35m")
var COLOR_PINK_D = Color(COLOR_HEADER + "2;35m")
var COLOR_PINK_I = Color(COLOR_HEADER + "3;35m")
var COLOR_PINK_U = Color(COLOR_HEADER + "4;35m")
var COLOR_PINK_BL = Color(COLOR_HEADER + "5;35m")
var COLOR_PINK_X = Color(COLOR_HEADER + "9;35m")
var COLOR_PINK_BG = Color(COLOR_HEADER + "0;45m")
var COLOR_PINK_IBG = Color(COLOR_HEADER + "0;105m")

var COLOR_CYAN = Color(COLOR_HEADER + "0;36m")
var COLOR_CYAN_L = Color(COLOR_HEADER + "1;36m")
var COLOR_CYAN_D = Color(COLOR_HEADER + "2;36m")
var COLOR_CYAN_I = Color(COLOR_HEADER + "3;36m")
var COLOR_CYAN_U = Color(COLOR_HEADER + "4;36m")
var COLOR_CYAN_BL = Color(COLOR_HEADER + "5;36m")
var COLOR_CYAN_X = Color(COLOR_HEADER + "9;36m")
var COLOR_CYAN_BG = Color(COLOR_HEADER + "0;46m")
var COLOR_CYAN_IBG = Color(COLOR_HEADER + "0;106m")

var COLOR_GRAY = Color(COLOR_HEADER + "0;37m")
var COLOR_GRAY_L = Color(COLOR_HEADER + "1;37m")
var COLOR_GRAY_D = Color(COLOR_HEADER + "2;37m")
var COLOR_GRAY_I = Color(COLOR_HEADER + "3;37m")
var COLOR_GRAY_U = Color(COLOR_HEADER + "4;37m")
var COLOR_GRAY_BL = Color(COLOR_HEADER + "5;37m")
var COLOR_GRAY_X = Color(COLOR_HEADER + "9;37m")
var COLOR_GRAY_BG = Color(COLOR_HEADER + "0;47m")
var COLOR_GRAY_IBG = Color(COLOR_HEADER + "0;107m")

var COLOR_RESET = Color(COLOR_HEADER + "0;0m")
