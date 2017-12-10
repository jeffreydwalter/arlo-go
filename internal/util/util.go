package util

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(data interface{}) string {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprint("error:", err)
	}
	return fmt.Sprint(string(j))
}

const TransIdPrefix = "web"

func GenTransId(transType string) {
	/*
		func divmod(numerator, denominator int64) (quotient, remainder int64) {
		    quotient = numerator / denominator // integer division, decimals are truncated
		    remainder = numerator % denominator
		    return
		}*/
}
