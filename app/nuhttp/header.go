package nuhttp

type headerValue struct {
	name  string
	Value string
}

type Header struct {
	Path   headerPath
	Values []headerValue
}

func (h Header) HasHeader(input string) bool {
	for _, head := range h.Values {
		if head.name == input {
			return true
		}
	}
	return false
}

// TODO: better way to handle... null/tryget type behavior?
func (h Header) GetHeader(input string) (*headerValue, bool) {
	for _, header := range h.Values {
		if header.name == input {
			return &header, true
		}
	}
	return nil, false
}

func (h *Header) SetHeaderValue(name string, value string) {
	hv, exists := h.GetHeader(name)
	if !exists {
		h.Values = append(h.Values, headerValue{name, value})
		return
	}

	hv.Value = value
}
