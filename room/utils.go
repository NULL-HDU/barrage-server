package room

import (
	"encoding/binary"
)

// mergeInfoListBytes merge InfoList byte into buffer.
// This function will rewrite the number of Infoes and connect rest bytes.
func mergeInfoListBytes(buf *[]byte, nb []byte) {
	if len(*buf) <= 0 {
		*buf = nb
		return
	}
	if len(nb) <= 4 {
		return
	}

	// calculate new number
	n1 := binary.BigEndian.Uint32(*buf)
	n2 := binary.BigEndian.Uint32(nb)
	num := n1 + n2
	if n1 > num {
		// n1 > num,  there are too many informations, drop them
		logger.Warnf("Too many informations in a list! Get %d.", n2)
		return
	}

	// set new number
	binary.BigEndian.PutUint32(*buf, num)
	// connect them
	*buf = append(*buf, nb[4:]...)
}
