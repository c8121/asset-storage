package storage

type XorEncoder interface {
	Encode(b []byte)
}

type Xor struct {
	Key []byte
	kl  int //key length
	ki  int //key index
}

// Encode does b ^ e.Key if key length greater 0
func (e *Xor) Encode(b []byte) {

	if e.kl == 0 {
		return
	}

	for i := range b {

		b[i] ^= e.Key[e.ki]

		if e.ki++; e.ki == e.kl {
			e.ki = 0
		}
	}
}
