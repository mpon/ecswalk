package sliceutil

// ChunkedSlice chunk slice by chunk size
// ref: https://stackoverflow.com/questions/35179656/slice-chunking-in-go
func ChunkedSlice(slice []string, chunkSize int) [][]string {
	chunked := [][]string{}
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunked = append(chunked, slice[i:end])
	}
	return chunked
}

// DistinctSlice to remove duplicate value
func DistinctSlice(slice []string) []string {
	distinct := []string{}
	sliceMap := map[string]bool{}
	for _, v := range slice {
		if !sliceMap[v] {
			sliceMap[v] = true
			distinct = append(distinct, v)
		}
	}
	return distinct
}
