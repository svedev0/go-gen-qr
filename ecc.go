// Reed–Solomon ECC and interleaving using a 256 element Galois field with
// primitive polynomial 0x11D.

package main

// ECC tables for level L only. Index by version 0..40 (0 unused). Values
// correspond to standard QR (Model 2) parameters.
var eccCodewordsPerBlockL = [...]int{
	-1,
	7, 10, 15, 20, 26, 18, 20, 24, 30,
	18, 20, 24, 26, 30, 22, 24, 28, 30, 28,
	28, 28, 28, 30, 30, 26, 28, 30, 30, 30,
	30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
	30,
}

var numEccBlocksL = [...]int{
	-1,
	1, 1, 1, 1, 1, 2, 2, 2, 2,
	4, 4, 4, 4, 4, 6, 6, 6, 6, 7,
	8, 8, 9, 9, 10, 12, 12, 12, 13, 14,
	15, 16, 17, 18, 19, 19, 20, 21, 22, 24,
	25,
}

func gfMul(x, y byte) byte {
	var z byte = 0
	for range 8 {
		if (y & 1) != 0 {
			z ^= x
		}

		var carry bool = (x & 0x80) != 0
		x <<= 1
		if carry {
			x ^= 0x1D // 0x11D reduced without the top bit
		}
		y >>= 1
	}
	return z
}

// rsGenerator returns generator polynomial coefficients of given degree,
// excluding the leading term (which is 1). Coefficients are in a form directly
// used by rsRemainder().
func rsGenerator(degree int) []byte {
	var generator []byte = make([]byte, degree)
	generator[degree-1] = 1
	var root byte = 1
	for range degree {
		for j := range degree {
			generator[j] = gfMul(generator[j], root)
			if j+1 < degree {
				generator[j] ^= generator[j+1]
			}
		}
		root = gfMul(root, 0x02) // 0x02 = alpha
	}
	return generator
}

func rsRemainder(data []byte, generator []byte) []byte {
	var degree int = len(generator)
	var result []byte = make([]byte, degree)
	for _, b := range data {
		var factor byte = b ^ result[0]
		copy(result, result[1:])
		result[degree-1] = 0
		for j := range degree {
			result[j] ^= gfMul(generator[j], factor)
		}
	}
	return result
}

func addEccAndInterleave(data []byte, version int) []byte {
	var numBlocks int = numEccBlocksL[version]
	var eccLen int = eccCodewordsPerBlockL[version]
	var generator []byte = rsGenerator(eccLen)

	// Split into blocks with "short" and "long" data lengths. In QR codes,
	// some blocks have 1 extra data codeword.
	var shortLen int = len(data) / numBlocks
	var numLong int = len(data) % numBlocks
	var longLen int = shortLen + 1

	var blocks [][]byte = make([][]byte, numBlocks)
	var eccs [][]byte = make([][]byte, numBlocks)
	var k int = 0
	for i := range numBlocks {
		var blockSize int = shortLen
		if i >= numBlocks-numLong {
			blockSize = longLen
		}

		var block []byte = make([]byte, blockSize)
		copy(block, data[k:k+blockSize])
		k += blockSize
		blocks[i] = block
		eccs[i] = rsRemainder(block, generator)
	}

	// Interleave data
	var result []byte = make([]byte, 0, (getNumRawDataModules(version) / 8))
	for i := range longLen {
		for b := range numBlocks {
			if i < len(blocks[b]) {
				result = append(result, blocks[b][i])
			}
		}
	}

	// Interleave ECC
	for i := range eccLen {
		for b := range numBlocks {
			result = append(result, eccs[b][i])
		}
	}

	return result
}
