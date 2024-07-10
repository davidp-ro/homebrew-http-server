package http

import s "strings"

func ExtractQueryParamsFrom(paramsString string) map[string]string {
	queryParams := s.Split(paramsString, "&")
	params := make(map[string]string)

	for _, param := range queryParams {
		paramKV := s.Split(param, "=")
		if len(paramKV) != 2 {
			// Invalid
			continue
		}
		params[paramKV[0]] = paramKV[1]
	}

	return params
}
