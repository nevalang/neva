//nolint:all // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package funcs

// normalizeSliceBounds applies Python-like slice normalization.
// Negative indices are interpreted from the end, all bounds are clamped to [0, length],
// and out-of-order ranges become empty slices.
func normalizeSliceBounds(from int64, to int64, length int64) (start int64, end int64) {
	start = normalizeSliceIndex(from, length)
	end = normalizeSliceIndex(to, length)
	if start > end {
		start = end
	}
	return start, end
}

// normalizeSliceIndex converts a potentially negative index into [0, length].
func normalizeSliceIndex(idx int64, length int64) int64 {
	if idx < 0 {
		idx += length
	}
	if idx < 0 {
		return 0
	}
	if idx > length {
		return length
	}
	return idx
}
