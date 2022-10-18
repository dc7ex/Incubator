package hack

import (
	"sort"
	"strings"
)

func Sign(data map[string]interface{}, key string) (sign string) {

	if len(data) == 0 {
		return ""
	}

	var d = make(map[string]string)
	for k, v := range data {
		d[k] = ToString(v)
	}

	return SignString(d, key)
}

func SignString(data map[string]string, key string) (sign string) {

	if len(data) == 0 {
		return ""
	}

	var strs []string
	for k, _ := range data {
		strs = append(strs, k)
	}

	sort.Strings(strs)

	var str string = ""
	for _, k := range strs {
		if k == "sign" || k == "signature" {
			continue
		}
		str += k + "=" + data[k] + "&"
	}

	str += "key=" + key
	sign = strings.ToLower(Sha1s(str))
	return
}
