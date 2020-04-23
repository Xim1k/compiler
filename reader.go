package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

var valuesTranslations = map[int]string{
  0: "itemError",
  1: "itemBool",
  2: "itemChar",
  3: "itemCharConstant",
  4: "itemAssign",
  5: "itemDeclare",
  6: "itemEOF",
  7: "itemFunction",
  8: "itemField" ,
  9: "itemIdentifier",
  10: "itemLeftDelim",
  11: "itemLeftParen",
  12: "itemNumber",
  13: "itemPipe",
  14: "itemRawString",
  15: "itemRightDelim",
  16: "itemRightParen",
  17: "itemSpace",
  18: "itemString",
  19: "itemText",
  20: "itemVariable",
  21: "itemKeyword",
  22: "itemBlock",
  23: "itemDot",
  24: "itemDefine",
  25: "itemElse",
  26: "itemEnd",
  27: "itemIf",
  28: "itemNil",
  29: "itemRange",
  30: "itemTemplate",
  31: "itemWith",
  32: "itemFor",
  33: "itemDoublePlus",
  34: "itemDoubleMinus",
  35: "itemMinus",
  36: "itemPlus",
  37: "itemNewLine",
  38: "itemPackage",
  39: "itemPackageValue",
  40: "itemImport",
  41: "itemImportValue",
  42: "itemFunctionDefine",
  43: "itemNotEqual",
  44: "itemFunctionName",
  45: "itemVariableType",
  46: "itemMap",
  47: "itemVar",
  48: "itemColon",
  49: "itemByteType",
  50: "itemStringType",
  51: "itemIntType",
  52: "itemUnknownToken",
  53: "itemSemiColon",
  54: "itemComment",
}

func main() {
    filename := os.Args[1]
    data, err := ioutil.ReadFile(filename)

    if err != nil {
        fmt.Println("File reading error", err)
        return
    }

    runLexer(lex(string(data), "", ""))
}

func runLexer(lex *lexer) {
  for x := range lex.items {
    if x.val == "\n" {
      x.val = "\\n"
    }

    fmt.Println("type -", valuesTranslations[int(x.typ)], "; value - \"", x.val , "\"; position -", x.pos + 1, "; line -", x.line)
  }
}
