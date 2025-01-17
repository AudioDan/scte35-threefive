package threefive

import (
	"fmt"
	"github.com/futzu/bitter"
)

// Cue a SCTE35 cue.
type Cue struct {
	InfoSection
	Command
	Descriptors []Descriptor `json:",omitempty"`
	Packet      PacketData   `json:",omitempty"`
}

// Decode extracts bits for the Cue values.
func (cue *Cue) Decode(bites []byte) bool {
	var bitn bitter.Bitn
	bitn.Load(bites)
	if !cue.InfoSection.Decode(&bitn) {
		return false
	}
	if CommandIsValid(cue.InfoSection.SpliceCommandType) {
		cue.Command = CommandDecoder(cue.InfoSection.SpliceCommandType)
		cue.Command.Decode(&bitn)
		cue.InfoSection.DescriptorLoopLength = bitn.AsUInt64(16)
		cue.dscptrLoop(&bitn)
		return true
	}
	return false
}

// DscptrLoop loops over any splice descriptors
func (cue *Cue) dscptrLoop(bitn *bitter.Bitn) {
	var i uint64
	i = 0
	l := cue.InfoSection.DescriptorLoopLength
	for i < l {
		tag := bitn.AsUInt8(8)
		i++
		length := bitn.AsUInt64(8)
		i++
		i += length
		if DescriptorIsValid(tag) {
			sdr := DescriptorDecoder(tag)
			sdr.Decode(bitn, tag, uint8(length))
			cue.Descriptors = append(cue.Descriptors, sdr)
		}
	}
}

//Show display SCTE-35 data as JSON.
func (cue *Cue) Show() {
	fmt.Println(MkJson(&cue))
}
