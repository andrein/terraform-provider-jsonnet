package jsonnet

func ExpandJPath(in []interface{}) []string {
	jpath := make([]string, len(in))

	for i, v := range in {
		jpath[i] = v.(string)
	}

	return jpath
}
