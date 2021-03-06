package util

import (
	"fmt"
	"math"
)

var idealFreqs = []float64{
	.0817, .0149, .0278, .0425, .1270, .0223, .0202, .0609, .0697, .0015, .0077, .0402, .0241,
	.0675, .0751, .0193, .0009, .0599, .0633, .0906, .0276, .0098, .0236, .0015, .0197, .0007}

func xorByte(a []byte, k byte) []byte {
	res := make([]byte, len(a))
	for i := range a {
		res[i] = a[i] ^ k
	}
	return res
}

func dotVec(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func lenVec(a []float64) float64 {
	return math.Sqrt(dotVec(a, a))
}

func cosine(a, b []float64) float64 {
	return dotVec(a, b) / (lenVec(a) * lenVec(b))
}

// ScoreText returns integer representing how likely seq a to be a regular english text
func ScoreText(a []byte) float64 {
	cts := make([]int, 26)
	for _, ch := range a {
		if 'A' <= ch && ch <= 'Z' {
			ch -= 32
		}
		if 'a' <= ch && ch <= 'z' {
			cts[int(ch)-'a']++
		}
	}
	amount := float64(len(a))
	freqs := make([]float64, 26)
	for i, c := range cts {
		freqs[i] = float64(c) / amount
	}
	// fmt.Println(freqs)
	return cosine(freqs, idealFreqs)
}

// Break1Xor returns most likely key for the sequence of bytes XORed with 1 byte
func BreakOneXor(a []byte) (byte, []byte) {
	var maxScore float64
	var maxKey byte
	var maxDecoded []byte
	for k := 0; k <= 255; k++ {
		decoded := xorByte(a, byte(k))
		score := ScoreText(decoded)
		if score > maxScore {
			maxScore = score
			maxKey = byte(k)
			maxDecoded = decoded
			if score > 0.5 {
				//fmt.Println(k, score, string(decoded))
			}
		}
	}
	return maxKey, maxDecoded
}

func hammingDistance(a, b []byte) int {
	sum := 0
	for i := range a {
		r := a[i] ^ b[i]
		for r > 0 {
			if r&1 == 1 {
				sum++
			}
			r = r >> 1
		}
	}
	return sum
}

// BreakRepeatingXor returns key used to XOR every block of size in the sequence bytes
func BreakRepeatingXor(bytes []byte, size int) []byte {
	res := make([]byte, size)
	blockSize := len(bytes) / size
	for i := 0; i < size; i++ {
		// determine i-th byte of key
		block := make([]byte, blockSize)
		for j := 0; j < blockSize; j++ {
			block[j] = bytes[i+j*size]
		}
		// fmt.Println(block)
		key, _ := BreakOneXor(block)
		res[i] = key
		// fmt.Println("block #", i, string(decoded))
	}

	return res
}

func RepeatingXorSize(bytes []byte) int {
	// determine key size
	bestLen := 0
	bestDist := 10000.0
	maxSize := 40
	blocks := len(bytes) / maxSize
	for keylen := 1; keylen < maxSize; keylen++ {
		dist := 0.0
		for i := 0; i < blocks; i++ {
			a := i * keylen
			b := (i + 1) * keylen
			c := (i + 2) * keylen
			dist += float64(hammingDistance(bytes[a:b], bytes[b:c])) / float64(keylen) // normalizing
			// fmt.Println(bytes[a:b], bytes[b:c])
		}
		dist /= float64(blocks) // averaging
		if dist < bestDist {
			bestDist = dist
			bestLen = keylen
		}
		fmt.Printf("size: %2d bits: %4.2f\n", keylen, dist)
	}
	fmt.Printf("best: %2d %4.2f\n", bestLen, bestDist)
	return bestLen
	// for ll := 2; ll < 40; ll++ {
	//	key := calcKeyXor(bytes, ll)
	//	fmt.Println(key, string(key))
	//	//decoded := encodeKeyXor(bytes, key)
	//	//fmt.Println(string(decoded[:10]))
	// }
}

func XorEncrypt(key, text []byte) []byte {
	ptr := 0
	res := make([]byte, len(text))
	for i := range text {
		res[i] = text[i] ^ key[ptr]
		ptr++
		if ptr >= len(key) {
			ptr = 0
		}
	}
	return res
}
