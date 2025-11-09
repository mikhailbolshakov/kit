package localization

func ToTranslatable[T any](v ...T) []Translatable {
	r := make([]Translatable, 0, len(v))
	for _, vv := range v {
		var vAny any = vv
		if _, ok := vAny.(Translatable); ok {
			r = append(r, vAny.(Translatable))
		}
	}
	return r
}
