package slip39_test

import (
	"strconv"
	"testing"

	slip39 "github.com/rddl-network/bc-slip39-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndCombine(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		shareCount        int
		groupThreshold    uint8
		groups            []slip39.GroupDescriptor
		secret            []byte
		password          string
		iterationExponent uint8
		resultCount       int
		resultStrings     []string
		resultWords       [][]uint16
	}{
		{
			shareCount:     5,
			groupThreshold: 1,
			groups: []slip39.GroupDescriptor{
				{
					Threshold: 3,
					Count:     5,
				},
			},
			secret:            append([]byte("totally secret!"), []byte{0}...),
			password:          "",
			iterationExponent: 0,
			resultCount:       5,
			resultStrings: []string{
				"check academic academic acne academic beard lying desert blessing grin syndrome magazine havoc resident depend tactics romantic engage erode guilt",   //nolint:dupword
				"check academic academic agree coastal plot human satisfy paces military science helpful therapy judicial usher liquid angel deny breathe mansion",     //nolint:dupword
				"check academic academic amazing buyer woman merit idle fiber evaluate knife bolt reunion true prize location valid payroll withdraw friar",            //nolint:dupword
				"check academic academic arcade dough endorse island losing tactics retailer grin verdict depict cover hesitate sweater else square prize robin",       //nolint:dupword
				"check academic academic axle crush teacher coding center aircraft image reaction clay unusual ancestor slow burden penalty moisture mountain peasant", //nolint:dupword
			},
			resultWords: [][]uint16{
				{136, 0, 0, 18, 156, 682, 445, 780, 636, 582, 789, 433, 912, 487, 963, 531, 42, 203, 98, 559},
				{136, 0, 0, 50, 237, 282, 479, 539, 891, 754, 409, 973, 206, 164, 436, 882, 272, 857, 698, 765},
				{136, 0, 0, 66, 177, 898, 157, 131, 20, 455, 730, 144, 958, 40, 825, 110, 657, 591, 596, 655},
			},
		},
	}

	for index, testCase := range testCases {
		index, testCase := index, testCase
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			t.Parallel()
			resultCount, wordsInEachShare, sharesBuffer, err := slip39.Generate(testCase.groupThreshold, testCase.groups, testCase.secret, testCase.password, testCase.iterationExponent, slip39.FakeRandom())
			require.NoError(t, err)
			assert.Equal(t, testCase.resultCount, resultCount)

			resultStrings := make([]string, len(testCase.resultStrings))
			for index := 0; index < testCase.shareCount; index++ {
				start := index * wordsInEachShare
				end := start + wordsInEachShare
				words := sharesBuffer[start:end]
				resultString, err := slip39.StringsForWords(words, wordsInEachShare)
				require.NoError(t, err)
				assert.Equal(t, testCase.resultStrings[index], resultString)
				resultStrings[index] = resultString
			}

			selectedShareIndexes := []int{1, 3, 4}
			selectedSharesLen := len(selectedShareIndexes)
			selectedSharesWords := make([][]uint16, len(testCase.resultWords))
			for index := 0; index < selectedSharesLen; index++ {
				selectedShareIndex := selectedShareIndexes[index]
				selectedShareString := resultStrings[selectedShareIndex]
				resultWords, err := slip39.WordsForStrings(selectedShareString, wordsInEachShare)
				require.NoError(t, err)
				assert.Equal(t, testCase.resultWords[index], resultWords)
				selectedSharesWords[index] = resultWords
			}

			secret, err := slip39.Combine(selectedSharesWords, testCase.password)
			require.NoError(t, err)
			assert.Equal(t, testCase.secret, secret)
		})
	}
}
