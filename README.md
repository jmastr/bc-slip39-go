# Go wrapper for BlockchainCommons' SLIP39 C implementation

## Usage

### Prerequisite

Get `BlockchainCommons` source code and compile it:
```
$ ./deps.sh
```
Headers will be placed into `/usr/local/include/` and libraries in `/usr/local/lib/`.

### Example - Create

Create a `3 of 5 Shamir Backup`:
```
package main

import (
	"fmt"

	slip39 "github.com/jmastr/bc-slip39-go"
)

func main() {
	groupThreshold := uint8(1)
	groups := []slip39.GroupDescriptor{
		{
			Threshold: 3,
			Count:     5,
		},
	}
	secret := []byte("totally secret!!")
	password := ""
	iterationExponent := uint8(0)
	resultCount, wordsInEachShare, sharesBuffer, err := slip39.Generate(groupThreshold, groups, secret, password, iterationExponent, slip39.Random())
	if err != nil {
		panic(err)
	}

	resultStrings := make([]string, resultCount)
	for index := 0; index < resultCount; index++ {
		start := index * wordsInEachShare
		end := start + wordsInEachShare
		words := sharesBuffer[start:end]
		resultString, err := slip39.StringsForWords(words, wordsInEachShare)
		if err != nil {
			panic(err)
		}
		resultStrings[index] = resultString
	}

	for _, resultString := range resultStrings {
		fmt.Printf("%s\n", resultString)
	}
}
```
Output:
```
$ go run main.go
yoga upgrade academic acne analysis maximum window tolerate slap fact briefing leaf unfair patent mild wine jacket airline making triumph
yoga upgrade academic agree agree surprise marvel disease papa dream silver bulge axis fatigue endless fumes finger fridge install guilt
yoga upgrade academic amazing dining news teammate hunting suitable organize capital spend adapt museum trouble seafood screw ranked budget force
yoga upgrade academic arcade domain silent muscle provide rebound username sled exhaust withdraw grownup dryer glasses true undergo tofu stick
yoga upgrade academic axle dance formal always smith safari mailman union crazy impulse spelling envelope object ending transfer cricket obtain
```

## Development

### Test

```
$ go test -v ./...
```
