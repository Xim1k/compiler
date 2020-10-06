package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "time"
)

var valuesTranslations = map[int]string{
  -1: "program_head",
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

    lexer := lex(string(data))
    parse(lexer)

    // for item := range lexer.items {
    //     fmt.Println("value: ",item.val, "; ", valuesTranslations[int(item.typ)], "; position:", int(item.line), ":", int(item.pos))
    // }
}

func parse(lex *lexer) {
  // var stack[]int
  tree := &AstTree{
    level: 0,
    typ: -1,
  }
  token := getNextToken(lex, false)

  stmt(tree, token, lex, 1)
  printTree(tree)
}

//main parse function

func stmt(tree *AstTree, token *item, lex *lexer, currentLevel int) {
  if token.typ == itemPackage {
    token = parseItemPackage(tree, token, lex, currentLevel)
  }

  if token.typ == itemImport {
    token = parseItemImport(tree, token, lex, currentLevel)
  }

  if token.typ == itemFunctionDefine {
    token = functionDefine(tree, token, lex, currentLevel)
  }

  // fmt.Println(valuesTranslations[int(token.typ)])
}

//non main parse functions
func functionDefine(tree *AstTree, token *item, lex *lexer, level int) *item {
  token = getNextToken(lex, false)

  if token.typ != itemFunctionName {
      parseErrorPrint(token, itemFunctionName)
  }

  tree.addChild(&AstTree{
      key: time.Now().String(),
      typ: itemFunctionDefine,
      level: level,
      text: "Function declaration",
      data: token.val,
    })

  token = getNextToken(lex, false);

  if token.typ != itemLeftParen {
    parseErrorPrint(token, itemLeftParen)
  }

  token = getNextToken(lex, false);

  if token.typ != itemRightParen {
    parseErrorPrint(token, itemRightParen)
  }

  token = getNextToken(lex, false);

  if token.typ != itemLeftDelim {
    parseErrorPrint(token, itemRightParen)
  }
}

func parseItemImport(tree *AstTree, token *item, lex *lexer, level int) *item {
  node := tree.addChild(&AstTree{
      key: time.Now().String(),
      typ: itemNode,
      level: level,
      text: "Preprocessor directive import",
    })

  node.addChild(&AstTree{
    key: time.Now().String(),
    typ: itemImport,
    data: token.val,
    level: level + 1,
    text: "Preprocessor directive",
  })

  token = getNextToken(lex, false)

  if token.typ != itemLeftParen {
    parseErrorPrint(token, itemLeftParen)
  }

  for token = getNextToken(lex, true); token.typ == itemString; token = getNextToken(lex, true) {
    node.addChild(&AstTree{
      key: time.Now().String(),
      data: token.val,
      typ: itemString,
      level: level + 1,
      text: "Called library",
    })
  }

  if token.typ != itemRightParen {
    parseErrorPrint(token, itemRightParen)
  }

  token = getNextToken(lex, false)

  if token.typ != itemNewLine {
    parseErrorPrint(token, itemNewLine)
  }

  return getNextToken(lex, true)
}

func parseItemPackage(tree *AstTree, token *item, lex *lexer, level int) *item {
  node := tree.addChild(&AstTree{
      key: time.Now().String(),
      typ: itemNode,
      level: level,
      text: "Preprocessor directive package",
    })

  node.addChild(&AstTree{
    key: time.Now().String(),
    typ: itemPackage,
    data: token.val,
    level: level + 1,
    text: "Preprocessor directive",
  })

  token = getNextToken(lex, false)

  if token.typ == itemPackageValue {
    node.addChild(&AstTree{
      key: time.Now().String(),
      typ: itemPackageValue,
      data: token.val,
      level: level + 1,
      text: "Preprocessor package value",
    })
  } else {
    parseErrorPrint(token, itemPackageValue)
  }

  token = getNextToken(lex, false)

  if token.typ != itemNewLine {
    parseErrorPrint(token, itemNewLine)
  }

  return getNextToken(lex, true)
}
//
// //support functions
//
func getNextToken(lex *lexer, skipNewLine bool) *item {
  token := <-lex.items

  if skipNewLine == true {
    for ;token.typ == itemSpace || token.typ == itemComment || token.typ == itemNewLine; {
      token = <-lex.items
    }
  } else {
    for ;token.typ == itemSpace || token.typ == itemComment; {
      token = <-lex.items
    }
  }

  return &token
}

func parseErrorPrint(token *item, item itemType)  {
  fmt.Println("There should be", valuesTranslations[int(item)], "<", token.pos, ">", "line:", token.line)
  os.Exit(1)
}
