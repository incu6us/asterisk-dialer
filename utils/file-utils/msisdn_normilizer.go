package file_utils

import (
	"strings"
)

func MSISDNNormalizer(str string) []string {
	numberList := make([]string, 0)

	if strings.HasPrefix(str, "+380"){
		str = strings.Replace(str, "+380", "0", 1)
	}

	if strings.HasPrefix(str, "380"){
		str = strings.Replace(str, "380", "0", 1)
	}

	list := strings.Split(
		strings.Replace(
			strings.Replace(str,
				"\n",
				"",
				-1),
			"\r",
			"",
			-1),
		";")


	for _, num := range list {
	    if num == "" {
	        continue
        }
		numberList = append(numberList, num)
	}

	return numberList
}
