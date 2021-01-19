package packet

import (
	"github.com/eren5960/gophertunnel408/minecraft/nbt"
	"github.com/eren5960/gophertunnel408/minecraft/protocol"
)

// BlockActorData is sent by the server to update data of a block entity client-side, for example the data of
// a chest.
type BlockActorData struct {
	// Position is the position of the block that holds the block entity. If no block entity is at this
	// position, the packet is ignored by the client.
	Position protocol.BlockPos
	// NBTData is the new data of the block that will be encoded to NBT and applied client-side, so that the
	// client can see the block update. The NBTData should contain all properties of the block, not just
	// properties that were changed.
	NBTData map[string]interface{}
}

// ID ...
func (*BlockActorData) ID() uint32 {
	return IDBlockActorData
}

// Marshal ...
func (pk *BlockActorData) Marshal(w *protocol.Writer) {
	w.UBlockPos(&pk.Position)
	w.NBT(&pk.NBTData, nbt.NetworkLittleEndian)
}

// Unmarshal ...
func (pk *BlockActorData) Unmarshal(r *protocol.Reader) {
	pk.NBTData = make(map[string]interface{})
	r.UBlockPos(&pk.Position)
	r.NBT(&pk.NBTData, nbt.NetworkLittleEndian)
}
