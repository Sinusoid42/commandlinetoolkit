package commandlinetoolkit

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
**********************************************************************
Return the Tree visually
*/

/*
**********************************************************************
parse a girven Argument Type
*/
func decodeArgType(m map[string]interface{}) (ArgumentType, error) {
	
	theType := ArgumentType(0)
	
	typeStr := ""
	
	var ok bool
	if typeStr, ok = m[TYPEKEY].(string); ok {
		
		typeStrArrOR := strings.Split(typeStr, "|")
		typeStrArrComma := strings.Split(typeStr, ",")
		
		theTypes := []string{}
		
		if len(typeStrArrOR) > 0 {
			theTypes = typeStrArrOR
			
		} else if len(typeStrArrComma) > 0 {
			theTypes = typeStrArrComma
		}
		for i, str := range theTypes {
			theTypes[i] = strings.TrimSpace(str)
		}
		
		for _, value := range theTypes {
			
			if strings.Compare(value, OPTIONSTRING) == 0 {
				
				theType |= OPTION
				
			}
			if strings.Compare(value, PARAMETERSTRING) == 0 {
				
				theType |= PARAMETER
			}
			if strings.Compare(value, WILDCARDSTRING) == 0 {
				
				theType |= WILDCARD
			}
			if strings.Compare(value, FLAGSTRING) == 0 {
				theType |= FLAG
			}
			if strings.Compare(value, COMMANDSTRING) == 0 {
				theType |= COMMAND
			}
		}
	}
	
	if theType&COMMAND > 0 && (theType&(OPTION|WILDCARD|PARAMETER) > 0) {
		return theType, errors.New("Could not parse the Type")
	}
	return theType, nil
}

func validateArgType(argtype ArgumentType) bool {
	if argtype&COMMAND > 0 && (argtype&(OPTION|WILDCARD|FLAG|PARAMETER) > 0) {
		return false
	}
	if argtype&PARAMETER > 0 && (argtype&(COMMAND|WILDCARD|FLAG|PARAMETER) > 0) {
		return false
	}
	if argtype&OPTION > 0 && (argtype&(COMMAND|PARAMETER) > 0) {
		return false
	}
	if argtype&WILDCARD > 0 && (argtype&(COMMAND|PARAMETER) > 0) {
		return false
	}
	if argtype&FLAG > 0 && (argtype&(COMMAND|PARAMETER) > 0) {
		return false
	}
	return true
}

func saveErrorFile(fileName string, content string) {
	
	yr, month, day := time.Now().Date()
	hr := time.Now().Hour()
	min := time.Now().Minute()
	
	theTime := strconv.Itoa(yr) + "_" + strconv.Itoa(int(month)) + "_" + strconv.Itoa(day) + "_" + strconv.Itoa(hr) + "_" + strconv.Itoa(min) + "_"
	
	file, _ := os.OpenFile(theTime+"errsave_"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	
	writer := bufio.NewWriter(file)
	
	writer.WriteString(content)
	
	writer.Flush()
	
	file.Close()
}

func copyProgramArgument(dst *map[string]interface{}, src *map[string]interface{}) {
	_copy(dst, src, TYPEKEY)
	
	_copy(dst, src, LONGFLAGKEY)
	_copy(dst, src, SHORTFLAGKEY)
	_copy(dst, src, HELPKEY)
	_copy(dst, src, SHORTHELPKEY)
	_copy(dst, src, DATATYPEKEY)
	_copy(dst, src, RUNKEY)
	_copy(dst, src, REQUIREDKEY)
	_copy(dst, src, MUTEABLEKEY)
	
	//do this with checkInputProgram casting
	//_copy(dst, src, ARGUMENTSKEY] = (*src)[ARGUMENTSKEY] since here we have no string
	_copy(dst, src, AUTHORKEY)
	_copy(dst, src, DESCRIPTIONKEY)
	_copy(dst, src, MANUALKEY)
	_copy(dst, src, VERSIONKEY)
	_copy(dst, src, EXECUTEABLEKEY)
}

func _copy(dst *map[string]interface{}, src *map[string]interface{}, typeString string) {
	if str, ok := (*src)[typeString].(string); ok {
		(*dst)[typeString] = str
	}
	if str, ok := (*src)[typeString].(int); ok {
		(*dst)[typeString] = str
	}
	if str, ok := (*src)[typeString].(bool); ok {
		(*dst)[typeString] = str
	}
}

func splitBySeperator(str []string) []string {
	
	//split
	if strings.Contains(str[0], "|") {
		str = strings.Split(str[0], "|")
	} else if strings.Contains(str[0], ",") {
		str = strings.Split(str[0], ",")
	} else if strings.Contains(str[0], " ") {
		d := newDebugHandler()
		
		s := "Multitype Error\n\""
		for _, k := range str {
			s += k
		}
		s += "\"\n"
		d.printError(s)
		return []string{}
	}
	
	for i, s := range str {
		
		if strings.Contains(s, " ") {
			s = strings.ReplaceAll(s, " ", "")
		}
		str[i] = s
	}
	return str
}
