package slip39

// #include <errno.h>
// #include <stdint.h>
// #include <stdio.h>
// #include <string.h>
// #include <sys/random.h>
//
// // Clearly not random. Only use for tests.
// void fake_random(uint8_t *buf, size_t count, void* ctx) {
//  uint8_t b = 0;
//  for(int i = 0; i < count; i++) {
//   buf[i] = b;
//   b = b + 17;
//  }
// }
//
// // Call libc to get random data for cryptographic purposes.
// void real_random(uint8_t *buf, size_t count, void* ctx) {
//  if(getrandom(buf, count, 0)==-1) {
//   printf("getrandom failed: %s\n", strerror(errno));
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

// Random returns `real_random` C function.
func Random() (random *[0]byte) {
	return (*[0]byte)(C.real_random)
}
