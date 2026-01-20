// Flags:
//   -c string   content to encode (required)
//   -o string   output PNG path (default "qrcode.png")
//   -s int      module scale (pixels per module) (default 10)

package main

import (
	"flag"
	"fmt"
	"image"
	"os"
)

const (
	modeByte = 0b0100 // Mode indicator for Byte mode
	ecLevelL = 0b01   // Format info ECC level bits (L)
)

func main() {
	var content string
	var outPath string
	var scale int

	flag.StringVar(&content, "c", "", "content to encode (required)")
	flag.StringVar(&outPath, "o", "qrcode.png", "output PNG path")
	flag.IntVar(&scale, "s", 10, "module scale (pixels per module)")
	flag.Parse()

	if content == "" {
		fmt.Fprintln(os.Stderr, "Error: missing required -c content")
		os.Exit(1)
	}
	if scale <= 0 {
		fmt.Fprintln(os.Stderr, "Error: scale must be > 0")
		os.Exit(1)
	}

	var payload []byte = []byte(content)

	// Pick smallest version that fits (Byte mode, ECC L).
	var version int
	var dataCodewords []byte
	version, dataCodewords, err := chooseVersionAndData(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// Full codewords (data + ecc), interleaved.
	var finalCodewords []byte = addEccAndInterleave(dataCodewords, version)
	var bits []bool = codewordsToBits(finalCodewords)

	var size int = 17 + 4*version
	var base *matrix = newMatrix(size)
	placeFunctionPatterns(base, version)
	placeDataBits(base, bits)

	// Choose best mask.
	var bestMask int = 0
	var bestScore int = int(^uint(0) >> 1)
	var bestMatrix *matrix

	for mask := range 8 {
		var cand *matrix = base.clone()
		applyMask(cand, mask)
		writeFormatInfo(cand, mask)
		if version >= 7 {
			writeVersionInfo(cand, version)
		}

		var score int = penaltyScore(cand)
		if score < bestScore {
			bestScore = score
			bestMask = mask
			bestMatrix = cand
		}
	}

	var qrImg *image.RGBA = renderPNG(bestMatrix, scale)
	if err := exportPNG(qrImg, outPath); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote %s (version %d, ECC L, mask %d)\n", outPath, version, bestMask)
}
