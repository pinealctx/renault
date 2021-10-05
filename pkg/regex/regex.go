package regex

import "regexp"

func MatchMoreStr(regexStr, source string, indexes ...int) ([]string, bool) {
	regex := regexp.MustCompile(regexStr)
	results := regex.FindAllStringSubmatch(source, -1)
	if len(results) == 0 {
		return nil, false
	}
	values := make([]string, 0)
	for _, index := range indexes {
		if len(results[0]) <= index {
			return nil, false
		}
		values = append(values, results[0][index])
	}
	return values, true
}

func MatchStr(regexStr, source string, index int) (string, bool) {
	regex := regexp.MustCompile(regexStr)
	results := regex.FindAllStringSubmatch(source, -1)
	if len(results) == 0 {
		return "", false
	}
	if len(results[0]) <= index {
		return "", false
	}
	return results[0][index], true
}

func Match(regexStr, source string) bool {
	var regex = regexp.MustCompile(regexStr)
	var results = regex.FindAllStringSubmatch(source, -1)
	return len(results) != 0
}
