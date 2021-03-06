package sfunc

import "strconv"

// StrToInt merubah str ke int dengan nilai default
func StrToInt(text string, defaultReturn int) int {
	number := defaultReturn
	if text != "" {
		var err error
		number, err = strconv.Atoi(text)
		if err != nil {
			number = defaultReturn
		}
	}
	return number
}
