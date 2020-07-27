package reedsolomon

// Addition, subtraction, multiplication, and division in GF(2^8).
// Operations are performed modulo x^8 + x^4 + x^3 + x^2 + 1.

// http://en.wikipedia.org/wiki/Finite_field_arithmetic

import "log"

const (
	gfZero = gfElement(0)
	gfOne  = gfElement(1)
)

var (
	gfExpTable = [256]gfElement{
		/*   0 -   9 */ 1, 2, 4, 8, 16, 32, 64, 128, 29, 58,
		/*  10 -  19 */ 116, 232, 205, 135, 19, 38, 76, 152, 45, 90,
		/*  20 -  29 */ 180, 117, 234, 201, 143, 3, 6, 12, 24, 48,
		/*  30 -  39 */ 96, 192, 157, 39, 78, 156, 37, 74, 148, 53,
		/*  40 -  49 */ 106, 212, 181, 119, 238, 193, 159, 35, 70, 140,
		/*  50 -  59 */ 5, 10, 20, 40, 80, 160, 93, 186, 105, 210,
		/*  60 -  69 */ 185, 111, 222, 161, 95, 190, 97, 194, 153, 47,
		/*  70 -  79 */ 94, 188, 101, 202, 137, 15, 30, 60, 120, 240,
		/*  80 -  89 */ 253, 231, 211, 187, 107, 214, 177, 127, 254, 225,
		/*  90 -  99 */ 223, 163, 91, 182, 113, 226, 217, 175, 67, 134,
		/* 100 - 109 */ 17, 34, 68, 136, 13, 26, 52, 104, 208, 189,
		/* 110 - 119 */ 103, 206, 129, 31, 62, 124, 248, 237, 199, 147,
		/* 120 - 129 */ 59, 118, 236, 197, 151, 51, 102, 204, 133, 23,
		/* 130 - 139 */ 46, 92, 184, 109, 218, 169, 79, 158, 33, 66,
		/* 140 - 149 */ 132, 21, 42, 84, 168, 77, 154, 41, 82, 164,
		/* 150 - 159 */ 85, 170, 73, 146, 57, 114, 228, 213, 183, 115,
		/* 160 - 169 */ 230, 209, 191, 99, 198, 145, 63, 126, 252, 229,
		/* 170 - 179 */ 215, 179, 123, 246, 241, 255, 227, 219, 171, 75,
		/* 180 - 189 */ 150, 49, 98, 196, 149, 55, 110, 220, 165, 87,
		/* 190 - 199 */ 174, 65, 130, 25, 50, 100, 200, 141, 7, 14,
		/* 200 - 209 */ 28, 56, 112, 224, 221, 167, 83, 166, 81, 162,
		/* 210 - 219 */ 89, 178, 121, 242, 249, 239, 195, 155, 43, 86,
		/* 220 - 229 */ 172, 69, 138, 9, 18, 36, 72, 144, 61, 122,
		/* 230 - 239 */ 244, 245, 247, 243, 251, 235, 203, 139, 11, 22,
		/* 240 - 249 */ 44, 88, 176, 125, 250, 233, 207, 131, 27, 54,
		/* 250 - 255 */ 108, 216, 173, 71, 142, 1}

	gfLogTable = [256]int{
		/*   0 -   9 */ -1, 0, 1, 25, 2, 50, 26, 198, 3, 223,
		/*  10 -  19 */ 51, 238, 27, 104, 199, 75, 4, 100, 224, 14,
		/*  20 -  29 */ 52, 141, 239, 129, 28, 193, 105, 248, 200, 8,
		/*  30 -  39 */ 76, 113, 5, 138, 101, 47, 225, 36, 15, 33,
		/*  40 -  49 */ 53, 147, 142, 218, 240, 18, 130, 69, 29, 181,
		/*  50 -  59 */ 194, 125, 106, 39, 249, 185, 201, 154, 9, 120,
		/*  60 -  69 */ 77, 228, 114, 166, 6, 191, 139, 98, 102, 221,
		/*  70 -  79 */ 48, 253, 226, 152, 37, 179, 16, 145, 34, 136,
		/*  80 -  89 */ 54, 208, 148, 206, 143, 150, 219, 189, 241, 210,
		/*  90 -  99 */ 19, 92, 131, 56, 70, 64, 30, 66, 182, 163,
		/* 100 - 109 */ 195, 72, 126, 110, 107, 58, 40, 84, 250, 133,
		/* 110 - 119 */ 186, 61, 202, 94, 155, 159, 10, 21, 121, 43,
		/* 120 - 129 */ 78, 212, 229, 172, 115, 243, 167, 87, 7, 112,
		/* 130 - 139 */ 192, 247, 140, 128, 99, 13, 103, 74, 222, 237,
		/* 140 - 149 */ 49, 197, 254, 24, 227, 165, 153, 119, 38, 184,
		/* 150 - 159 */ 180, 124, 17, 68, 146, 217, 35, 32, 137, 46,
		/* 160 - 169 */ 55, 63, 209, 91, 149, 188, 207, 205, 144, 135,
		/* 170 - 179 */ 151, 178, 220, 252, 190, 97, 242, 86, 211, 171,
		/* 180 - 189 */ 20, 42, 93, 158, 132, 60, 57, 83, 71, 109,
		/* 190 - 199 */ 65, 162, 31, 45, 67, 216, 183, 123, 164, 118,
		/* 200 - 209 */ 196, 23, 73, 236, 127, 12, 111, 246, 108, 161,
		/* 210 - 219 */ 59, 82, 41, 157, 85, 170, 251, 96, 134, 177,
		/* 220 - 229 */ 187, 204, 62, 90, 203, 89, 95, 176, 156, 169,
		/* 230 - 239 */ 160, 81, 11, 245, 22, 235, 122, 117, 44, 215,
		/* 240 - 249 */ 79, 174, 213, 233, 230, 231, 173, 232, 116, 214,
		/* 250 - 255 */ 244, 234, 168, 80, 88, 175}
)

// gfElement is an element in GF(2^8).
type gfElement uint8

// newGFElement creates and returns a new gfElement.
func newGFElement(data byte) gfElement {
	return gfElement(data)
}

// gfAdd returns a + b.
func gfAdd(a, b gfElement) gfElement {
	return a ^ b
}

// gfSub returns a - b.
//
// Note addition is equivalent to subtraction in GF(2).
func gfSub(a, b gfElement) gfElement {
	return a ^ b
}

// gfMultiply returns a * b.
func gfMultiply(a, b gfElement) gfElement {
	if a == gfZero || b == gfZero {
		return gfZero
	}

	return gfExpTable[(gfLogTable[a]+gfLogTable[b])%255]
}

// gfDivide returns a / b.
//
// Divide by zero results in a panic.
func gfDivide(a, b gfElement) gfElement {
	if a == gfZero {
		return gfZero
	} else if b == gfZero {
		log.Panicln("Divide by zero")
	}

	return gfMultiply(a, gfInverse(b))
}

// gfInverse returns the multiplicative inverse of a, a^-1.
//
// a * a^-1 = 1
func gfInverse(a gfElement) gfElement {
	if a == gfZero {
		log.Panicln("No multiplicative inverse of 0")
	}

	return gfExpTable[255-gfLogTable[a]]
}

// a^i   | bits      | polynomial                                   | decimal
// --------------------------------------------------------------------------
// 0     | 000000000 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 0
// a^0   | 000000001 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 1
// a^1   | 000000010 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 2
// a^2   | 000000100 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 4
// a^3   | 000001000 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 8
// a^4   | 000010000 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 16
// a^5   | 000100000 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 32
// a^6   | 001000000 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 64
// a^7   | 010000000 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 128
// a^8   | 000011101 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 29
// a^9   | 000111010 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 58
// a^10  | 001110100 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 116
// a^11  | 011101000 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 232
// a^12  | 011001101 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 205
// a^13  | 010000111 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 135
// a^14  | 000010011 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 19
// a^15  | 000100110 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 38
// a^16  | 001001100 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 76
// a^17  | 010011000 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 152
// a^18  | 000101101 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 45
// a^19  | 001011010 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 90
// a^20  | 010110100 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 180
// a^21  | 001110101 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 117
// a^22  | 011101010 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 234
// a^23  | 011001001 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 201
// a^24  | 010001111 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 143
// a^25  | 000000011 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 3
// a^26  | 000000110 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 6
// a^27  | 000001100 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 12
// a^28  | 000011000 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 24
// a^29  | 000110000 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 48
// a^30  | 001100000 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 96
// a^31  | 011000000 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 192
// a^32  | 010011101 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 157
// a^33  | 000100111 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 39
// a^34  | 001001110 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 78
// a^35  | 010011100 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 156
// a^36  | 000100101 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 37
// a^37  | 001001010 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 74
// a^38  | 010010100 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 148
// a^39  | 000110101 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 53
// a^40  | 001101010 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 106
// a^41  | 011010100 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 212
// a^42  | 010110101 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 181
// a^43  | 001110111 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 119
// a^44  | 011101110 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 238
// a^45  | 011000001 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 193
// a^46  | 010011111 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 159
// a^47  | 000100011 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 35
// a^48  | 001000110 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 70
// a^49  | 010001100 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 140
// a^50  | 000000101 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 5
// a^51  | 000001010 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 10
// a^52  | 000010100 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 20
// a^53  | 000101000 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 40
// a^54  | 001010000 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 80
// a^55  | 010100000 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 160
// a^56  | 001011101 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 93
// a^57  | 010111010 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 186
// a^58  | 001101001 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 105
// a^59  | 011010010 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 210
// a^60  | 010111001 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 185
// a^61  | 001101111 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 111
// a^62  | 011011110 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 222
// a^63  | 010100001 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 161
// a^64  | 001011111 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 95
// a^65  | 010111110 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 190
// a^66  | 001100001 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 97
// a^67  | 011000010 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 194
// a^68  | 010011001 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 153
// a^69  | 000101111 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 47
// a^70  | 001011110 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 94
// a^71  | 010111100 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 188
// a^72  | 001100101 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 101
// a^73  | 011001010 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 202
// a^74  | 010001001 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 137
// a^75  | 000001111 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 15
// a^76  | 000011110 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 30
// a^77  | 000111100 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 60
// a^78  | 001111000 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 120
// a^79  | 011110000 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 240
// a^80  | 011111101 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 253
// a^81  | 011100111 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 231
// a^82  | 011010011 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 211
// a^83  | 010111011 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 187
// a^84  | 001101011 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 107
// a^85  | 011010110 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 214
// a^86  | 010110001 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 177
// a^87  | 001111111 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 127
// a^88  | 011111110 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 254
// a^89  | 011100001 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 225
// a^90  | 011011111 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 223
// a^91  | 010100011 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 163
// a^92  | 001011011 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 91
// a^93  | 010110110 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 182
// a^94  | 001110001 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 113
// a^95  | 011100010 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 226
// a^96  | 011011001 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 217
// a^97  | 010101111 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 175
// a^98  | 001000011 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 67
// a^99  | 010000110 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 134
// a^100 | 000010001 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 17
// a^101 | 000100010 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 34
// a^102 | 001000100 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 68
// a^103 | 010001000 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 136
// a^104 | 000001101 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 13
// a^105 | 000011010 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 26
// a^106 | 000110100 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 52
// a^107 | 001101000 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 104
// a^108 | 011010000 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 208
// a^109 | 010111101 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 189
// a^110 | 001100111 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 103
// a^111 | 011001110 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 206
// a^112 | 010000001 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 129
// a^113 | 000011111 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 31
// a^114 | 000111110 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 62
// a^115 | 001111100 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 124
// a^116 | 011111000 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 248
// a^117 | 011101101 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 237
// a^118 | 011000111 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 199
// a^119 | 010010011 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 147
// a^120 | 000111011 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 59
// a^121 | 001110110 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 118
// a^122 | 011101100 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 236
// a^123 | 011000101 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 197
// a^124 | 010010111 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 151
// a^125 | 000110011 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 51
// a^126 | 001100110 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 102
// a^127 | 011001100 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 204
// a^128 | 010000101 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 133
// a^129 | 000010111 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 23
// a^130 | 000101110 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 46
// a^131 | 001011100 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 92
// a^132 | 010111000 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 184
// a^133 | 001101101 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 109
// a^134 | 011011010 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 218
// a^135 | 010101001 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 169
// a^136 | 001001111 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 79
// a^137 | 010011110 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 158
// a^138 | 000100001 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 33
// a^139 | 001000010 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 66
// a^140 | 010000100 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 132
// a^141 | 000010101 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 21
// a^142 | 000101010 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 42
// a^143 | 001010100 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 84
// a^144 | 010101000 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 168
// a^145 | 001001101 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 77
// a^146 | 010011010 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 154
// a^147 | 000101001 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 41
// a^148 | 001010010 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 82
// a^149 | 010100100 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 164
// a^150 | 001010101 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 85
// a^151 | 010101010 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 170
// a^152 | 001001001 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 73
// a^153 | 010010010 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 146
// a^154 | 000111001 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 57
// a^155 | 001110010 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 114
// a^156 | 011100100 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 228
// a^157 | 011010101 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 213
// a^158 | 010110111 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 183
// a^159 | 001110011 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 115
// a^160 | 011100110 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 230
// a^161 | 011010001 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 209
// a^162 | 010111111 | 0x^8 1x^7 0x^6 1x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 191
// a^163 | 001100011 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 99
// a^164 | 011000110 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 198
// a^165 | 010010001 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 145
// a^166 | 000111111 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 63
// a^167 | 001111110 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 1x^2 1x^1 0x^0 | 126
// a^168 | 011111100 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 252
// a^169 | 011100101 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 229
// a^170 | 011010111 | 0x^8 1x^7 1x^6 0x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 215
// a^171 | 010110011 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 179
// a^172 | 001111011 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 123
// a^173 | 011110110 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 246
// a^174 | 011110001 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 241
// a^175 | 011111111 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 1x^2 1x^1 1x^0 | 255
// a^176 | 011100011 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 227
// a^177 | 011011011 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 219
// a^178 | 010101011 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 171
// a^179 | 001001011 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 75
// a^180 | 010010110 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 150
// a^181 | 000110001 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 49
// a^182 | 001100010 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 98
// a^183 | 011000100 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 196
// a^184 | 010010101 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 149
// a^185 | 000110111 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 55
// a^186 | 001101110 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 110
// a^187 | 011011100 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 220
// a^188 | 010100101 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 165
// a^189 | 001010111 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 87
// a^190 | 010101110 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 174
// a^191 | 001000001 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 65
// a^192 | 010000010 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 130
// a^193 | 000011001 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 25
// a^194 | 000110010 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 50
// a^195 | 001100100 | 0x^8 0x^7 1x^6 1x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 100
// a^196 | 011001000 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 200
// a^197 | 010001101 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 141
// a^198 | 000000111 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 7
// a^199 | 000001110 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 14
// a^200 | 000011100 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 1x^2 0x^1 0x^0 | 28
// a^201 | 000111000 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 56
// a^202 | 001110000 | 0x^8 0x^7 1x^6 1x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 112
// a^203 | 011100000 | 0x^8 1x^7 1x^6 1x^5 0x^4 0x^3 0x^2 0x^1 0x^0 | 224
// a^204 | 011011101 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 221
// a^205 | 010100111 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 167
// a^206 | 001010011 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 83
// a^207 | 010100110 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 1x^2 1x^1 0x^0 | 166
// a^208 | 001010001 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 0x^2 0x^1 1x^0 | 81
// a^209 | 010100010 | 0x^8 1x^7 0x^6 1x^5 0x^4 0x^3 0x^2 1x^1 0x^0 | 162
// a^210 | 001011001 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 89
// a^211 | 010110010 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 178
// a^212 | 001111001 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 121
// a^213 | 011110010 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 242
// a^214 | 011111001 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 0x^2 0x^1 1x^0 | 249
// a^215 | 011101111 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 239
// a^216 | 011000011 | 0x^8 1x^7 1x^6 0x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 195
// a^217 | 010011011 | 0x^8 1x^7 0x^6 0x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 155
// a^218 | 000101011 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 43
// a^219 | 001010110 | 0x^8 0x^7 1x^6 0x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 86
// a^220 | 010101100 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 172
// a^221 | 001000101 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 1x^2 0x^1 1x^0 | 69
// a^222 | 010001010 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 0x^2 1x^1 0x^0 | 138
// a^223 | 000001001 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 9
// a^224 | 000010010 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 0x^2 1x^1 0x^0 | 18
// a^225 | 000100100 | 0x^8 0x^7 0x^6 1x^5 0x^4 0x^3 1x^2 0x^1 0x^0 | 36
// a^226 | 001001000 | 0x^8 0x^7 1x^6 0x^5 0x^4 1x^3 0x^2 0x^1 0x^0 | 72
// a^227 | 010010000 | 0x^8 1x^7 0x^6 0x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 144
// a^228 | 000111101 | 0x^8 0x^7 0x^6 1x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 61
// a^229 | 001111010 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 122
// a^230 | 011110100 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 1x^2 0x^1 0x^0 | 244
// a^231 | 011110101 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 1x^2 0x^1 1x^0 | 245
// a^232 | 011110111 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 1x^2 1x^1 1x^0 | 247
// a^233 | 011110011 | 0x^8 1x^7 1x^6 1x^5 1x^4 0x^3 0x^2 1x^1 1x^0 | 243
// a^234 | 011111011 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 251
// a^235 | 011101011 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 235
// a^236 | 011001011 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 203
// a^237 | 010001011 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 139
// a^238 | 000001011 | 0x^8 0x^7 0x^6 0x^5 0x^4 1x^3 0x^2 1x^1 1x^0 | 11
// a^239 | 000010110 | 0x^8 0x^7 0x^6 0x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 22
// a^240 | 000101100 | 0x^8 0x^7 0x^6 1x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 44
// a^241 | 001011000 | 0x^8 0x^7 1x^6 0x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 88
// a^242 | 010110000 | 0x^8 1x^7 0x^6 1x^5 1x^4 0x^3 0x^2 0x^1 0x^0 | 176
// a^243 | 001111101 | 0x^8 0x^7 1x^6 1x^5 1x^4 1x^3 1x^2 0x^1 1x^0 | 125
// a^244 | 011111010 | 0x^8 1x^7 1x^6 1x^5 1x^4 1x^3 0x^2 1x^1 0x^0 | 250
// a^245 | 011101001 | 0x^8 1x^7 1x^6 1x^5 0x^4 1x^3 0x^2 0x^1 1x^0 | 233
// a^246 | 011001111 | 0x^8 1x^7 1x^6 0x^5 0x^4 1x^3 1x^2 1x^1 1x^0 | 207
// a^247 | 010000011 | 0x^8 1x^7 0x^6 0x^5 0x^4 0x^3 0x^2 1x^1 1x^0 | 131
// a^248 | 000011011 | 0x^8 0x^7 0x^6 0x^5 1x^4 1x^3 0x^2 1x^1 1x^0 | 27
// a^249 | 000110110 | 0x^8 0x^7 0x^6 1x^5 1x^4 0x^3 1x^2 1x^1 0x^0 | 54
// a^250 | 001101100 | 0x^8 0x^7 1x^6 1x^5 0x^4 1x^3 1x^2 0x^1 0x^0 | 108
// a^251 | 011011000 | 0x^8 1x^7 1x^6 0x^5 1x^4 1x^3 0x^2 0x^1 0x^0 | 216
// a^252 | 010101101 | 0x^8 1x^7 0x^6 1x^5 0x^4 1x^3 1x^2 0x^1 1x^0 | 173
// a^253 | 001000111 | 0x^8 0x^7 1x^6 0x^5 0x^4 0x^3 1x^2 1x^1 1x^0 | 71
// a^254 | 010001110 | 0x^8 1x^7 0x^6 0x^5 0x^4 1x^3 1x^2 1x^1 0x^0 | 142
// a^255 | 000000001 | 0x^8 0x^7 0x^6 0x^5 0x^4 0x^3 0x^2 0x^1 1x^0 | 1
