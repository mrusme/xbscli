// LZ77 implementation for Go
//
// The implementation here was basically "translated" from
// https://github.com/rotemdan/lzutf8.js and hence contains the original
// author's namings and comments. It's a very crude and hacky implementation
// because I really didn't care to do it properly and just wanted a working
// LZ77 implementation, which Go seems to lack. In fact this had cost me two
// hours of my life I'll never get back, meh.

package lz77

import "math"

var MaximumMatchDistance = 32767

var outputBuffer []byte
var outputPosition = 0

var inputBufferRemainder []byte = nil
var outputBufferRemainder []byte

func DecompressBlockToString(input []byte) string {
	return string(decompressBlock(input))
}

func decompressBlock(input []byte) []byte {
	if inputBufferRemainder != nil {
		input = append(inputBufferRemainder, input...)
		inputBufferRemainder = nil
	}

	var outputStartPosition = cropOutputBufferToWindowAndInitialize(int(math.Max(float64(len(input)*4), 1024)))

	inputLength := len(input)
	for readPosition := 0; readPosition < inputLength; readPosition++ {
		var inputValue = input[readPosition]

		if (inputValue >> 6) != 3 {
			// If at the continuation byte of a UTF-8 codepoint sequence, output the literal value and continue
			outputByte(inputValue)
			continue
		}

		// At this point it is known that the current byte is the lead byte of either a UTF-8 codepoint or a sized pointer sequence.
		var sequenceLengthIdentifier = inputValue >> 5 // 6 for 2 bytes, 7 for at least 3 bytes

		// If bytes in read position imply the start of a truncated input sequence (either a literal codepoint or a pointer)
		// keep the remainder to be decoded with the next buffer
		if readPosition == inputLength-1 ||
			(readPosition == inputLength-2 && sequenceLengthIdentifier == 7) {
			inputBufferRemainder = input[readPosition:]
			break
		}

		// If at the leading byte of a UTF-8 codepoint byte sequence
		if input[readPosition+1]>>7 == 1 {
			// Output the literal value
			outputByte(inputValue)
		} else {
			// Beginning of a pointer sequence
			var matchLength = inputValue & 31
			var matchDistance = 0

			if sequenceLengthIdentifier == 6 { // 2 byte pointer type, distance was smaller than 128
				matchDistance = int(input[readPosition+1])
				readPosition += 1
			} else { // 3 byte pointer type, distance was greater or equal to 128
				matchDistance = (int(input[readPosition+1]) << 8) | int(input[readPosition+2]) // Big endian
				readPosition += 2
			}

			var matchPosition = outputPosition - matchDistance

			// Copy the match bytes to output
			for offset := 0; offset < int(matchLength); offset++ {
				outputByte(outputBuffer[matchPosition+offset])
			}
		}
	}

	rollBackIfOutputBufferEndsWithATruncatedMultibyteSequence()
	return getCroppedBuffer(outputBuffer, outputStartPosition, outputPosition-outputStartPosition, 0)
}

func outputByte(value byte) {
	outputPosition++
	outputBuffer = append(outputBuffer, value)
}

func cropOutputBufferToWindowAndInitialize(initialCapacity int) int {
	if outputBuffer == nil {
		outputBuffer = []byte{}
		return 0
	}

	var cropLength = int(math.Min(float64(outputPosition), float64(MaximumMatchDistance)))
	outputBuffer = getCroppedBuffer(outputBuffer, outputPosition-(cropLength), cropLength, initialCapacity)

	outputPosition = cropLength

	if outputBufferRemainder != nil {
		for i := 0; i < len(outputBufferRemainder); i++ {
			outputByte(outputBufferRemainder[i])
		}

		outputBufferRemainder = nil
	}

	return int(cropLength)
}

func rollBackIfOutputBufferEndsWithATruncatedMultibyteSequence() {
	for offset := 1; offset <= 4 && outputPosition-offset >= 0; offset++ {
		var value = outputBuffer[outputPosition-offset]

		if (offset < 4 && (value>>3) == 30) || // Leading byte of a 4 byte UTF8 sequence
			(offset < 3 && (value>>4) == 14) || // Leading byte of a 3 byte UTF8 sequence
			(offset < 2 && (value>>5) == 6) { // Leading byte of a 2 byte UTF8 sequence

			outputBufferRemainder = outputBuffer[outputPosition-offset : outputPosition]
			outputPosition -= offset

			return
		}
	}
}

func getCroppedBuffer(buffer []byte, cropStartOffset int, cropLength int, additionalCapacity int) []byte {
	var croppedBuffer = []byte{}
	croppedBuffer = append(croppedBuffer, buffer[cropStartOffset:cropStartOffset+cropLength]...)

	return croppedBuffer
}
