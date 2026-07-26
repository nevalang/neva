package messages

// Match reports whether msg matches pattern.
//
// Union patterns without data match any union value with the same tag. All
// other messages use value equality.
func Match(msg Msg, pattern Msg) bool {
	msgUnion, isUnion := msg.(UnionMsg)
	if !isUnion {
		return Equal(msg, pattern)
	}

	patternUnion, isPatternUnion := pattern.(UnionMsg)
	if !isPatternUnion {
		return Equal(msg, pattern)
	}
	if msgUnion.tag != patternUnion.tag {
		return false
	}
	if patternUnion.data == nil {
		return true
	}
	if msgUnion.data == nil {
		return false
	}
	return Equal(msgUnion.data, patternUnion.data)
}
