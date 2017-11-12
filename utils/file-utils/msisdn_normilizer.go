package file_utils

import (
	"strings"
)

func MSISDNNormalizer(str string) []string {
	numberList := make([]string, 0)

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
        if strings.HasPrefix(num, "+380"){
            num = strings.Replace(num, "+380", "0", 1)
        }

        if strings.HasPrefix(num, "380"){
            num = strings.Replace(num, "380", "0", 1)
        }

		numberList = append(numberList, num)
	}

	return numberList
}
