package helper

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func ToPtrSlice(strs []string) []*string {
	ptrs := make([]*string, len(strs))
	for i, s := range strs {
		ptrs[i] = &s
	}
	return ptrs
}
