package slip39

// #cgo CFLAGS: -I/usr/local/include/
// #cgo LDFLAGS: -L/usr/local/lib/ -lbc-slip39 -lbc-crypto-base -lbc-shamir
// #include <bc-slip39/bc-slip39.h>
//
// int slip39_combine_1d(
//  const uint16_t *mnemonics_1d, // array of pointers to 10-bit words
//  uint32_t mnemonics_words,     // number of words in each shard
//  uint32_t mnemonics_shards,    // total number of shards
//  const char *passphrase,       // passphrase to unlock master secret
//  const char **passwords,       // passwords protecting shards
//  uint8_t *buffer,              // working space, and place to return secret
//  uint32_t buffer_length        // total amount of working space
// ) {
//  // convert 1D mnemonics_1d array to 2D array
//  const uint16_t *mnemonics[mnemonics_shards];
//  uint16_t mnemonics_temp[mnemonics_shards][mnemonics_words];
//  int c = 0;
//  for (int i = 0; i < mnemonics_shards; i++) {
//          for (int j = 0; j < mnemonics_words; j++) {
//                  mnemonics_temp[i][j] = mnemonics_1d[c];
//                  c++;
//          }
//          mnemonics[i] = mnemonics_temp[i];
//  }
//  return slip39_combine(mnemonics, mnemonics_words, mnemonics_shards, passphrase, passwords, buffer, buffer_length);
// }
import "C"

import (
	"errors"
)

var ErrInsufficientSpace = errors.New("insufficient space")

// Generate calls `slip39_generate` C function.
func Generate(groupThreshold uint8, groups []GroupDescriptor, secret []byte, password string, iterationExponent uint8, randomGenerator *[0]byte) (result, wordsInEachShare int, sharesBuffer []uint16, err error) {
	groupCount := uint8(len(groups))
	secretLen := len(secret)
	wordsInEachShare = 0
	sharesBufferSize := 1024
	sharesBuffer = make([]uint16, sharesBufferSize)

	groupThresholdC := C.uint8_t(groupThreshold)
	groupsC := make([]C.group_descriptor, groupCount)
	for i := uint8(0); i < groupCount; i++ {
		groupsC[i].threshold = C.uint8_t(groups[i].Threshold)
		groupsC[i].count = C.uint8_t(groups[i].Count)
	}
	groupCountC := C.uint8_t(groupCount)
	secretDataC := (*C.uint8_t)(&secret[0])
	secretLenC := C.uint32_t(secretLen)
	passwordC := C.CString(password)
	iterationExponentC := C.uint8_t(iterationExponent)
	wordsInEachShareC := C.uint32_t(wordsInEachShare)
	sharesBufferC := (*C.ushort)(&sharesBuffer[0])
	sharesBufferSizeC := C.uint32_t(sharesBufferSize)

	res, err := C.slip39_generate(groupThresholdC, &groupsC[0], groupCountC, secretDataC, secretLenC, passwordC, iterationExponentC, &wordsInEachShareC, sharesBufferC, sharesBufferSizeC, nil, randomGenerator) //nolint:nlreturn
	result = int(res)
	wordsInEachShare = int(wordsInEachShareC)

	return
}

// StringsForWords calls `slip39_strings_for_words` C function.
func StringsForWords(words []uint16, wordsInEachShare int) (result string, err error) {
	wordsC := (*C.ushort)(&words[0])
	res, err := C.slip39_strings_for_words(wordsC, C.size_t(wordsInEachShare))
	result = C.GoString(res)

	return
}

// WordsForStrings calls `slip39_words_for_strings` C function.
func WordsForStrings(strings string, wordsInEachShare int) (words []uint16, err error) {
	words = make([]uint16, wordsInEachShare)

	stringsC := C.CString(strings)
	wordsC := (*C.uint16_t)(&words[0])
	wordsInEachShareC := C.uint32_t(wordsInEachShare)
	_, err = C.slip39_words_for_strings(stringsC, wordsC, wordsInEachShareC)

	return
}

// Combine calls `slip39_combine` C function.
func Combine(mnemonics [][]uint16, passphrase string) (secret []byte, err error) {
	mnemonicsShards := len(mnemonics)
	if mnemonicsShards == 0 {
		err = ErrInvalidLength

		return
	}
	mnemonicsWords := len(mnemonics[0])
	var mnemonics1D []uint16
	for i := 0; i < mnemonicsShards; i++ {
		mnemonics1D = append(mnemonics1D, mnemonics[i]...)
	}
	outputSecretDataLen := 1024
	secret = make([]uint8, outputSecretDataLen)

	mnemonics1DC := (*C.uint16_t)(&mnemonics1D[0])
	mnemonicsWordsC := C.uint32_t(mnemonicsWords)
	mnemonicsShardsC := C.uint32_t(mnemonicsShards)
	passphraseC := C.CString(passphrase)
	outputSecretDataC := (*C.uint8_t)(&secret[0])
	outputSecretDataLenC := C.uint32_t(outputSecretDataLen)
	secretLen, err := C.slip39_combine_1d(mnemonics1DC, mnemonicsWordsC, mnemonicsShardsC, passphraseC, nil, outputSecretDataC, outputSecretDataLenC)
	if err != nil {
		return
	}
	if secretLen == C.ERROR_INSUFFICIENT_SPACE {
		err = ErrInsufficientSpace

		return
	}
	secret = secret[:secretLen]

	return
}
