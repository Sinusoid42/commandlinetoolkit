package commandlinetoolkit

type Color string

var COL_HEADER = string("\033[")

var BLACK = Color(COL_HEADER + "0;30m")
var BLACK_L = Color(COL_HEADER + "1;30m")
var BLACK_D = Color(COL_HEADER + "2;30m")
var BLACK_I = Color(COL_HEADER + "3;30m")
var BLACK_U = Color(COL_HEADER + "4;30m")
var BLACK_BL = Color(COL_HEADER + "5;30m")

var RED = Color(COL_HEADER + "0;31m")
var RED_L = Color(COL_HEADER + "1;31m")
var RED_D = Color(COL_HEADER + "2;31m")
var RED_I = Color(COL_HEADER + "3;31m")
var RED_U = Color(COL_HEADER + "4;31m")
var RED_BL = Color(COL_HEADER + "5;31m")

var GREEN = Color(COL_HEADER + "0;32m")
var GREEN_L = Color(COL_HEADER + "1;32m")
var GREEN_D = Color(COL_HEADER + "2;32m")
var GREEN_I = Color(COL_HEADER + "3;32m")
var GREEN_U = Color(COL_HEADER + "4;32m")
var GREEN_BL = Color(COL_HEADER + "5;32m")

var ORANGE = Color(COL_HEADER + "0;33m")
var ORANGE_L = Color(COL_HEADER + "1;33m")
var ORANGE_D = Color(COL_HEADER + "2;33m")
var ORANGE_I = Color(COL_HEADER + "3;33m")
var ORANGE_U = Color(COL_HEADER + "4;33m")
var ORANGE_BL = Color(COL_HEADER + "5;33m")

var BLUE = Color(COL_HEADER + "0;34m")
var BLUE_L = Color(COL_HEADER + "1;34m")
var BLUE_D = Color(COL_HEADER + "2;34m")
var BLUE_I = Color(COL_HEADER + "3;34m")
var BLUE_U = Color(COL_HEADER + "4;34m")
var BLUE_BL = Color(COL_HEADER + "5;34m")

var PINK = Color(COL_HEADER + "0;35m")
var PINK_L = Color(COL_HEADER + "1;35m")
var PINK_D = Color(COL_HEADER + "2;35m")
var PINK_I = Color(COL_HEADER + "3;35m")
var PINK_U = Color(COL_HEADER + "4;35m")
var PINK_BL = Color(COL_HEADER + "5;35m")

var CYAN = Color(COL_HEADER + "0;36m")
var CYAN_L = Color(COL_HEADER + "1;36m")
var CYAN_D = Color(COL_HEADER + "2;36m")
var CYAN_I = Color(COL_HEADER + "3;36m")
var CYAN_U = Color(COL_HEADER + "4;36m")
var CYAN_BL = Color(COL_HEADER + "5;36m")

var COLOR_RESET = Color(COL_HEADER + "0;0m")
