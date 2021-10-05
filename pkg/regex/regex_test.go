package regex

import "testing"

func TestMatch(t *testing.T) {
	const regexStr = ".{1,}://.{1,}"

	t.Log(Match(regexStr, "sss://ssss"))
	t.Log(Match(regexStr, "://"))
	t.Log(Match(regexStr, "sss"))
	t.Log(Match(regexStr, "sss://"))
	t.Log(Match(regexStr, "://ssss"))
}
