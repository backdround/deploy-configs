package config

import (
	"fmt"
	"strings"
)

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func assertNoTab(data string) {
	if strings.Contains(data, "\t") {
		message := fmt.Sprintf("data contains tab: \n%v", data)
		panic(message)
	}
}
