package main

import "testing"

type testPairKey struct {
    string string
    expectedKeys []itemType
}

var tests = []testPairKey{
    { "if i := 0 {\n}", []itemType{itemIf, itemSpace, itemIdentifier, itemSpace, itemDeclare, itemSpace, itemNumber, itemSpace, itemLeftDelim, itemNewLine, itemRightDelim} },
    { "range", []itemType{itemRange} },
    { "{####}", []itemType{itemLeftDelim, itemUnknownToken, itemRightDelim} },
    { "{/*asdasd*/asdasd.Atoi()}", []itemType{itemLeftDelim, itemComment, itemIdentifier, itemFunction, itemLeftParen, itemRightParen, itemRightDelim} },
}

func TestKey(t *testing.T) {
    for pairNumber, pair := range tests {
        lexer := lex(pair.string, "", "")
        i := 0

        for x := range lexer.items {
          if x.typ != pair.expectedKeys[i] {
              t.Error(
                  "Expected", valuesTranslations[int(pair.expectedKeys[i])],
                  "got", valuesTranslations[int(x.typ)],
                  "in pair", pairNumber + 1,
              )
          }

          i++
        }
    }
}
