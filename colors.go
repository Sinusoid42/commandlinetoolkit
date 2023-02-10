package commandlinetoolkit

type Color string

var COL_HEADER = string("\033[")

var COLOR_BLACK = Color(COL_HEADER + "0;30m")
var COLOR_BLACK_L = Color(COL_HEADER + "1;30m")
var COLOR_BLACK_D = Color(COL_HEADER + "2;30m")
var COLOR_BLACK_I = Color(COL_HEADER + "3;30m")
var COLOR_BLACK_U = Color(COL_HEADER + "4;30m")
var COLOR_BLACK_BL = Color(COL_HEADER + "5;30m")

var COLOR_RED = Color(COL_HEADER + "0;31m")
var COLOR_RED_L = Color(COL_HEADER + "1;31m")
var COLOR_RED_D = Color(COL_HEADER + "2;31m")
var COLOR_RED_I = Color(COL_HEADER + "3;31m")
var COLOR_RED_U = Color(COL_HEADER + "4;31m")
var COLOR_RED_BL = Color(COL_HEADER + "5;31m")

var COLOR_GREEN = Color(COL_HEADER + "0;32m")
var COLOR_GREEN_L = Color(COL_HEADER + "1;32m")
var COLOR_GREEN_D = Color(COL_HEADER + "2;32m")
var COLOR_GREEN_I = Color(COL_HEADER + "3;32m")
var COLOR_GREEN_U = Color(COL_HEADER + "4;32m")
var COLOR_GREEN_BL = Color(COL_HEADER + "5;32m")

var COLOR_ORANGE = Color(COL_HEADER + "0;33m")
var COLOR_ORANGE_L = Color(COL_HEADER + "1;33m")
var COLOR_ORANGE_D = Color(COL_HEADER + "2;33m")
var COLOR_ORANGE_I = Color(COL_HEADER + "3;33m")
var COLOR_ORANGE_U = Color(COL_HEADER + "4;33m")
var COLOR_ORANGE_BL = Color(COL_HEADER + "5;33m")

var COLOR_BLUE = Color(COL_HEADER + "0;34m")
var COLOR_BLUE_L = Color(COL_HEADER + "1;34m")
var COLOR_BLUE_D = Color(COL_HEADER + "2;34m")
var COLOR_BLUE_I = Color(COL_HEADER + "3;34m")
var COLOR_BLUE_U = Color(COL_HEADER + "4;34m")
var COLOR_BLUE_BL = Color(COL_HEADER + "5;34m")

var COLOR_PINK = Color(COL_HEADER + "0;35m")
var COLOR_PINK_L = Color(COL_HEADER + "1;35m")
var COLOR_PINK_D = Color(COL_HEADER + "2;35m")
var COLOR_PINK_I = Color(COL_HEADER + "3;35m")
var COLOR_PINK_U = Color(COL_HEADER + "4;35m")
var COLOR_PINK_BL = Color(COL_HEADER + "5;35m")

var COLOR_CYAN = Color(COL_HEADER + "0;36m")
var COLOR_CYAN_L = Color(COL_HEADER + "1;36m")
var COLOR_CYAN_D = Color(COL_HEADER + "2;36m")
var COLOR_CYAN_I = Color(COL_HEADER + "3;36m")
var COLOR_CYAN_U = Color(COL_HEADER + "4;36m")
var COLOR_CYAN_BL = Color(COL_HEADER + "5;36m")

var COLOR_GRAY = Color(COL_HEADER + "0;37m")
var COLOR_GRAY_L = Color(COL_HEADER + "1;37m")
var COLOR_GRAY_D = Color(COL_HEADER + "2;37m")
var COLOR_GRAY_I = Color(COL_HEADER + "3;37m")
var COLOR_GRAY_U = Color(COL_HEADER + "4;37m")
var COLOR_GRAY_BL = Color(COL_HEADER + "5;37m")

var COLOR_RESET = Color(COL_HEADER + "0;0m")
