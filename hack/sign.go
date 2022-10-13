package hack

import (
	"sort"
	"strings"
)

func Sign(data map[string]interface{}, key string) (sign string) {

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
		val := ToString(data[k])
		str += k + "=" + val + "&"
	}

	str += "key=" + key
	sign = strings.ToLower(Sha1s(str))
	return
}
