package interfaces

import "encoding/json"

// OPNsense reports the wireless config on non-wireless interfaces as an
// empty string instead of an object, so InterfaceWireless (generated in
// overview.go) needs a tolerant UnmarshalJSON to fall back to the zero value
// when it isn't a JSON object.
func (w *InterfaceWireless) UnmarshalJSON(data []byte) error {
	type alias InterfaceWireless
	var a alias
	if err := json.Unmarshal(data, &a); err == nil {
		*w = InterfaceWireless(a)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*w = InterfaceWireless{}
		return nil
	}

	return json.Unmarshal(data, &a)
}
