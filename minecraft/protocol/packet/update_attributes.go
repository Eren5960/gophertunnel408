package packet

import (
	"github.com/eren5960/gophertunnel408/minecraft/protocol"
)

// UpdateAttributes is sent by the server to update an amount of attributes of any entity in the world. These
// attributes include ones such as the health or the movement speed of the entity.
type UpdateAttributes struct {
	// EntityRuntimeID is the runtime ID of the entity. The runtime ID is unique for each world session, and
	// entities are generally identified in packets using this runtime ID.
	EntityRuntimeID uint64
	// Attributes is a slice of new attributes that the entity gets. It includes attributes such as its
	// health, movement speed, etc. Note that only changed attributes have to be sent in this packet. It is
	// not required to send attributes that did not have their values changed.
	Attributes []protocol.Attribute
}

// ID ...
func (*UpdateAttributes) ID() uint32 {
	return IDUpdateAttributes
}

// Marshal ...
func (pk *UpdateAttributes) Marshal(w *protocol.Writer) {
	w.Varuint64(&pk.EntityRuntimeID)
	protocol.WriteAttributes(w, &pk.Attributes)
}

// Unmarshal ...
func (pk *UpdateAttributes) Unmarshal(r *protocol.Reader) {
	r.Varuint64(&pk.EntityRuntimeID)
	protocol.Attributes(r, &pk.Attributes)
}
