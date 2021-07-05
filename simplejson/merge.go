package simplejson

// MergeFrom merge the given JSON recursively.
func (j *JSON) MergeFrom(from *JSON) {
	j.data = mergeFromInternal(from.data, j.data)
}

func mergeFromInternal(from interface{}, to interface{}) interface{} {
	if to == nil {
		return from
	}
	switch fromV := from.(type) {
	case map[string]interface{}:
		switch toV := to.(type) {
		case map[string]interface{}:
			for k, v := range fromV {
				oldV, ok := toV[k]
				if ok {
					toV[k] = mergeFromInternal(v, oldV)
					continue
				}
				toV[k] = v
			}
			return to
		default:
			return from
		}
	default:
		return from
	}
}
