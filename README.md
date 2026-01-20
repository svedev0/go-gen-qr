# go-gen-qr

A QR code generator written in Go using only the standard library.

Notes:
- Mode: Byte mode only (UTF-8 bytes).
- Error correction: Level L only.
- Mask pattern: Chosen by penalty scoring across masks 0..7.
