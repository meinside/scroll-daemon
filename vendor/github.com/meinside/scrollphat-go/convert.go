package scroll

import (
	"log"
)

// Convert given 0/1 arrays to drawable bytes
func ArraysToBytes(arrays [][]uint8) []byte {
	if len(arrays) < NumRows || len(arrays[0]) < NumCols {
		log.Printf("wrong arrays length")
		return nil
	}

	bytes := []byte{}

	var b byte
	for col := 0; col < NumCols; col++ {
		b = 0
		for row := 0; row < NumRows; row++ {
			if arrays[row][col] > 0 {
				b |= (1 << uint(row))
			}
		}
		bytes = append(bytes, b)
	}

	return bytes
}
