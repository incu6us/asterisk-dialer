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

	listLen := len(list) - 1

	for _, num := range list[:listLen] {
		numberList = append(numberList, num)
	}

	return numberList
}
