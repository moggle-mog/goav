// sips from https://github.com/google/gopacket/blob/master/layers/base.go
package sips

// field token
const (
	_FieldBase = iota
	_FieldName
	_FieldUser
	_FieldHost
	_FieldPort
	_FieldID
	_FieldMethod
)

// BaseLayer is a convenience struct which implements the LayerData and
// LayerPayload functions of the Layer interface.
type BaseLayer struct {
	// Contents is the set of bytes that make up this layer.  IE: for an
	// Ethernet packet, this would be the set of bytes making up the
	// Ethernet frame.
	Contents []byte
	// Payload is the set of bytes contained by (but not part of) this
	// Layer.  Again, to take Ethernet as an example, this would be the
	// set of bytes encapsulated by the Ethernet protocol.
	Payload []byte
}

// LayerContents returns the bytes of the packet layer.
func (b *BaseLayer) LayerContents() []byte { return b.Contents }

// LayerPayload returns the bytes contained within the packet layer.
func (b *BaseLayer) LayerPayload() []byte { return b.Payload }

// Substring gets a string from a slice of bytes
// Checks the bounds to avoid any range errors
func Substring(v string, from, to int) string {

	// Remove negative values
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = 0
	}

	// Limit if over len
	if from > len(v) || from > to {
		return ""
	}
	if to > len(v) {
		return v[from:]
	}

	return v[from:to]
}
