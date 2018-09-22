/*
 * Copyright (c) 2018 Jeffrey Walter <jeffreydwalter@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func PrettyPrint(data interface{}) string {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprint("error:", err)
	}
	return fmt.Sprint(string(j))
}

func FloatToHex(x float64) string {
	var result []byte
	quotient := int(x)
	fraction := x - float64(quotient)

	for quotient > 0 {
		quotient = int(x / 16)
		remainder := int(x - (float64(quotient) * 16))

		if remainder > 9 {
			result = append([]byte{byte(remainder + 55)}, result...)
		} else {
			for _, c := range strconv.Itoa(int(remainder)) {
				result = append([]byte{byte(c)}, result...)
			}
		}

		x = float64(quotient)
	}

	if fraction == 0 {
		return string(result)
	}

	result = append(result, '.')

	for fraction > 0 {
		fraction = fraction * 16
		integer := int(fraction)
		fraction = fraction - float64(integer)

		if integer > 9 {
			result = append(result, byte(integer+55))
		} else {
			for _, c := range strconv.Itoa(int(integer)) {
				result = append(result, byte(c))
			}
		}
	}

	return string(result)
}

func HeaderToArray(header http.Header) (res []string) {
	for name, values := range header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	return
}

func HeaderToMap(header http.Header) map[string]string {
	h := make(map[string]string)
	for name, values := range header {
		for _, value := range values {
			h[name] = value
		}
	}
	return h
}
