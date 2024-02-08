package slip39

// #include <stdint.h>
//
// // Clearly not random. Only use for tests.
// void fake_random(uint8_t *buf, size_t count, void* ctx) {
//  uint8_t b = 0;
//  for(int i = 0; i < count; i++) {
//   buf[i] = b;
//   b = b + 17;
//  }
// }
import "C"

import (
	"errors"
)

var ErrInvalidLength = errors.New("invalid length")

// FakeRandom returns `fake_random` C function.
func FakeRandom() (fakeRandom *[0]byte) {
	return (*[0]byte)(C.fake_random)
}
