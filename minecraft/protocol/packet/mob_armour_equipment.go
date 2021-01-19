package packet

import (
	"github.com/eren5960/gophertunnel408/minecraft/protocol"
)

// MobArmourEquipment is sent by the server to the client to update the armour an entity is wearing. It is
// sent for both players and other entities, such as zombies.
type MobArmourEquipment struct {
	// EntityRuntimeID is the runtime ID of the entity. The runtime ID is unique for each world session, and
	// entities are generally identified in packets using this runtime ID.
	EntityRuntimeID uint64
	// Helmet is the equipped helmet of the entity. Items that are not wearable on the head will not be
	// rendered by the client. Unlike in Java Edition, blocks cannot be worn.
	Helmet protocol.ItemStack
	// Chestplate is the chestplate of the entity. Items that are not wearable as chestplate will not be
	// rendered.
	Chestplate protocol.ItemStack
	// Leggings is the item worn as leggings by the entity. Items not wearable as leggings will not be
	// rendered client-side.
	Leggings protocol.ItemStack
	// Boots is the item worn as boots by the entity. Items not wearable as boots will not be rendered.
	Boots protocol.ItemStack
}

// ID ...
func (*MobArmourEquipment) ID() uint32 {
	return IDMobArmourEquipment
}

// Marshal ...
func (pk *MobArmourEquipment) Marshal(w *protocol.Writer) {
	w.Varuint64(&pk.EntityRuntimeID)
	w.Item(&pk.Helmet)
	w.Item(&pk.Chestplate)
	w.Item(&pk.Leggings)
	w.Item(&pk.Boots)
}

// Unmarshal ...
func (pk *MobArmourEquipment) Unmarshal(r *protocol.Reader) {
	r.Varuint64(&pk.EntityRuntimeID)
	r.Item(&pk.Helmet)
	r.Item(&pk.Chestplate)
	r.Item(&pk.Leggings)
	r.Item(&pk.Boots)
}
