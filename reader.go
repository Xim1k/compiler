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
  55: "itemNode",
  56: "itemCalledLibrary",
  57: "itemEqual",
  58: "itemMupltiply",
  59: "itemDivide",
  60: "itemNot",
  61: "itemGreater",
  62: "itemLower",
  63: "itemGreaterOrEqual",
  64: "itemLowerOrEqual",
  65: "itemOr",
  66: "itemAnd",
  67: "itemReturn",
  68: "itemBoolType",
  69: "itemIndex",
  70: "itemRest",
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
  tree := &AstTree{
    level: 0,
    typ: -1,
  }

  token := getNextToken(lex, true)

  token = parseProgram(tree, tree, token, lex, 1)
  printTree(tree)
}

func debug(token *item) {
    fmt.Println(valuesTranslations[int(token.typ)])
    fmt.Println(token.val)
}


func parseProgram(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ != itemPackage {
        parseErrorPrint(token, itemPackage)
    } else {
        token = parsePackage(tree, node, token, lex, currentLevel)
    }

    if token.typ != itemImport && token.typ != itemFunctionDefine {
        //syntax error: non-declaration statement outside function body
        parseErrorPrint(token, itemFunctionDefine)
    }

    if token.typ == itemImport {
        token = parseImport(tree, node, token, lex, currentLevel)
    }

    if token.typ == itemFunctionDefine {
        token = parseFunctionsList(tree, node, token, lex, currentLevel)
    } else {
        // expect at least main()
        parseErrorPrint(token, itemFunctionDefine)
    }

    return token
}

func parseImport(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    packageNode := node.addChild(&AstTree{
        key: time.Now().String(),
        typ: itemNode,
        level: currentLevel,
        text: "Import definition",
    })

    token = getNextToken(lex, false)

    if token.val == "(" && token.typ == itemLeftParen {
        libsNode := packageNode.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel + 2,
                text: "Imported libs",
        })

        token = getNextToken(lex, true)

        token = parseImportsValue(tree, libsNode, token, lex, currentLevel + 4)

        if token.typ != itemRightParen {
            parseErrorPrint(token, itemRightParen)
        }

        token = getNextToken(lex, true)
    } else {
        parseErrorPrint(token, itemLeftParen)
    }

    return token
}

func parseImportsValue(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemString {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemString,
                level: currentLevel,
                text: "Lib",
                data: token.val,
        })

        token = getNextToken(lex, false)

        if token.typ != itemNewLine && token.typ != itemRightParen {
            parseErrorPrint(token, itemRightParen)
        }

        if token.typ == itemRightParen {
            return token
        }

        token = getNextToken(lex, false)

        if token.typ == itemRightParen {
            return token
        }

        if token.typ == itemString {
            token = parseImportsValue(tree, node, token, lex, currentLevel)
        } else {
            parseErrorPrint(token, itemString)
        }
    } else {
        parseErrorPrint(token, itemString)
    }

    return token
}

func parsePackage(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    packageNode := node.addChild(&AstTree{
        key: time.Now().String(),
        typ: itemNode,
        level: currentLevel,
        text: "Package definition",
    })

    token = getNextToken(lex, false)

    if token.typ == itemPackageValue {
        packageNode.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemPackageValue,
            level: currentLevel + 2,
            text: "itemPackageValue",
            data: token.val,
        })

        token = getNextToken(lex, true)
    } else {
        parseErrorPrint(token, itemPackageValue)
    }

    return token
}

func parseFunction(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    node.addChild(&AstTree{
        key: time.Now().String(),
        typ: itemFunctionName,
        level: currentLevel,
        text: "Function name",
        data: token.val,
    })

    isMain := false

    if token.val == "main" {
        isMain = true
    }

    token = getNextToken(lex, false)

    if token.val == "(" && token.typ == itemLeftParen {
        parametersNode := node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel,
                text: "Parameters of function",
        })

        token = getNextToken(lex, false)

        if token.typ == itemIdentifier && isMain == false {
            token = parseParameters(tree, parametersNode, token, lex, currentLevel + 2)
        }

        if token.typ != itemRightParen {
            parseErrorPrint(token, itemRightParen)
        }

        token = getNextToken(lex, false)

        if token.typ == itemLeftDelim {
            bodyNode := node.addChild(&AstTree{
                    key: time.Now().String(),
                    typ: itemNode,
                    level: currentLevel,
                    text: "Body of function",
            })

            token = getNextToken(lex, true)

            token = parseInstructionList(tree, bodyNode, token, lex, currentLevel + 2)

            if token.typ != itemRightDelim {
                parseErrorPrint(token, itemRightDelim)
            }
        } else {
            parseErrorPrint(token, itemLeftDelim)
        }
    } else {
        parseErrorPrint(token, itemLeftParen)
    }

    return getNextToken(lex, false)
}

//main parse function
func parseFunctionsList(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemError {
        return token
    }

    if token.typ == itemFunctionDefine {
        functionNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNode,
            level: currentLevel,
            text: "Function",
        })

        token = getNextToken(lex, false)

        if token.typ == itemFunctionName {
            token = parseFunction(tree, functionNode, token, lex, currentLevel + 2)
            token = getNextToken(lex, true)
        } else {
            parseErrorPrint(token, itemFunctionName)
        }
    } else {
        parseErrorPrint(token, itemFunctionDefine)
    }

    return parseFunctionsList(tree, node, token, lex, currentLevel)
}

func parseParameters(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    token = parseParameter(tree, node, token, lex, currentLevel)

    if token.val == "," && token.typ == itemChar {
        token = getNextToken(lex, false)
        token = parseParameters(tree, node, token, lex, currentLevel)
    }

    return token
}

func parseParameter(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ != itemIdentifier {
        parseErrorPrint(token, itemIdentifier)
    }

    identifierNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemIdentifier,
            level: currentLevel,
            text: "Identifier",
            data: token.val,
    })

    token = getNextToken(lex, false)

    if token.typ == itemIntType || token.typ == itemStringType || token.typ == itemBoolType || token.typ == itemByteType {
        identifierNode.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemVariableType,
            level: currentLevel + 2,
            text: "Variable type",
            data: token.val,
        })
    } else {
        parseErrorPrint(token, itemVariableType)
    }

    return getNextToken(lex, false)
}

//main parse function
func parseInstructionList(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemRightDelim {
        return token
    }

    instructionNode := node.addChild(&AstTree{
        key: time.Now().String(),
        typ: itemNode,
        level: currentLevel,
        text: "itemInstruction",
    })

    token = parseInstruction(tree, instructionNode, token, lex, currentLevel + 2)

    if token.typ == itemRightDelim {
        return token
    }

    if token.typ != itemNewLine && token.typ != itemSemiColon {
        parseErrorPrint(token, itemNewLine)
    }

    token = getNextToken(lex, true)

    return parseInstructionList(tree, node, token, lex, currentLevel)
}

func parseInstruction(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.val == "return" && token.typ == itemReturn {
        childNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemReturn,
            level: currentLevel,
            text: "itemReturn",
            data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseExpression(tree, childNode, token, lex, currentLevel + 2)

        return token
    }

    if itemVar == token.typ {
        return parseDeclaration(tree, node, token, lex, currentLevel)
    }

    if itemIdentifier == token.typ {
        return parseInstructionExpression(tree, node, token, lex, currentLevel)
    }

    if token.val == "if" && token.typ == itemIf {
        structureNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNode,
            level: currentLevel,
            text: "If structure",
        })

        token = getNextToken(lex, false)

        conditionNode := structureNode.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNode,
            level: currentLevel + 2,
            text: "Condition",
        })

        token = parseExpression(tree, conditionNode, token, lex, currentLevel + 4)

        if token.typ == itemLeftDelim {
            bodyNode := structureNode.addChild(&AstTree{
                    key: time.Now().String(),
                    typ: itemNode,
                    level: currentLevel + 2,
                    text: "Body of structure",
            })

            token = getNextToken(lex, true)

            token = parseInstructionList(tree, bodyNode, token, lex, currentLevel + 4)

            if token.typ != itemRightDelim {
                parseErrorPrint(token, itemRightDelim)
            }

            token = getNextToken(lex, false)

            if token.typ == itemElse {
                elseNode := structureNode.addChild(&AstTree{
                        key: time.Now().String(),
                        typ: itemNode,
                        level: currentLevel + 2,
                        text: "Else structure",
                })

                token = getNextToken(lex, true)
                token = getNextToken(lex, true)

                token = parseInstructionList(tree, elseNode, token, lex, currentLevel + 4)

                if token.typ != itemRightDelim {
                    parseErrorPrint(token, itemRightDelim)
                }

                token = getNextToken(lex, false)
            }
        } else {
            parseErrorPrint(token, itemLeftDelim)
        }

        return token
    }

    if token.val == "for" && token.typ == itemFor {
        structureNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNode,
            level: currentLevel,
            text: "For (while) structure",
        })

        token = getNextToken(lex, false)

        conditionNode := structureNode.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNode,
            level: currentLevel + 2,
            text: "Condition",
        })

        token = parseExpression(tree, conditionNode, token, lex, currentLevel + 4)

        if token.typ == itemLeftDelim {
            bodyNode := structureNode.addChild(&AstTree{
                    key: time.Now().String(),
                    typ: itemNode,
                    level: currentLevel + 2,
                    text: "Body of structure",
            })

            token = getNextToken(lex, true)

            token = parseInstructionList(tree, bodyNode, token, lex, currentLevel + 4)

            if token.typ != itemRightDelim {
                parseErrorPrint(token, itemRightDelim)
            }

            token = getNextToken(lex, false)
        } else {
            parseErrorPrint(token, itemLeftDelim)
        }

        return token
    }

    return token
}

func parseExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    parentNode := node.addChild(&AstTree{
        key: time.Now().String(),
        typ: itemNode,
        level: currentLevel,
        text: "Expression",
    })

    token = parseLogicalExpression(tree, parentNode, token, lex, currentLevel + 2)
    token = parseExtendedExpression(tree, parentNode, token, lex, currentLevel + 2)

    return token
}

func parseExtendedExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemOr || token.typ == itemAnd {
        token = parseLogicalOperator(tree, node, token, lex, currentLevel)
        token = parseLogicalExpression(tree, node, token, lex, currentLevel)
        token = parseExtendedExpression(tree, node, token, lex, currentLevel)
    }

    return token
}

func parseLogicalOperator(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemOr || token.val == "||" {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemOr,
            level: currentLevel,
            text: "itemOr",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    if token.typ == itemAnd || token.val == "&&" {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemAnd,
            level: currentLevel,
            text: "itemAnd",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    return token
}

func parseLogicalExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.val == "!" && token.typ == itemNot {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNot,
                level: currentLevel,
                text: "itemNot",
                data: token.val,
        })

        token = getNextToken(lex, false)
    }

    token = parseSimpleExpression(tree, node, token, lex, currentLevel)
    token = parseExtendedLogicalExpression(tree, node, token, lex, currentLevel)

    return token
}

func parseDeclaration(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    declarationNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNode,
            level: currentLevel,
            text: "Declaration",
    })

    token = getNextToken(lex, false)

    if token.typ != itemIdentifier {
        parseErrorPrint(token, itemIdentifier)
    }

    identifierNode := declarationNode.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemIdentifier,
            level: currentLevel + 2,
            text: "itemIdentifier",
            data: token.val,
    })

    token = getNextToken(lex, false)

    if token.typ == itemIntType || token.typ == itemStringType || token.typ == itemBoolType || token.typ == itemByteType {
        declarationNode.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemVariableType,
            level: currentLevel + 2,
            text: "Variable type",
            data: token.val,
        })

        token = getNextToken(lex, false)

        if token.typ == itemNewLine {
            return token
        }

        if token.typ != itemAssign {
            parseErrorPrint(token, itemAssign)
        }

        token = getNextToken(lex, false)

        return parseExpression(tree, declarationNode, token, lex, currentLevel + 2)
    }

    if token.typ == itemAssign {
        token = getNextToken(lex, false)

        if token.val == "[" && token.typ == itemChar {
            indexNode := identifierNode.addChild(&AstTree{
                    key: time.Now().String(),
                    typ: itemNode,
                    level: currentLevel + 2,
                    text: "itemArraySize",
                    data: token.val,
            })

            token = getNextToken(lex, false)

            token = parseSimpleExpression(tree, indexNode, token, lex, currentLevel + 4)

            if token.val != "]" || token.typ != itemChar {
                parseErrorPrint(token, itemChar)
            }

            token = getNextToken(lex, false)

            if token.typ == itemIntType || token.typ == itemStringType || token.typ == itemBoolType || token.typ == itemByteType {
                declarationNode.addChild(&AstTree{
                    key: time.Now().String(),
                    typ: itemVariableType,
                    level: currentLevel + 2,
                    text: "Variable type",
                    data: token.val,
                })
            } else {
                parseErrorPrint(token, itemVariableType)
            }

            token = getNextToken(lex, false)

            if token.typ != itemLeftDelim {
                parseErrorPrint(token, itemLeftDelim)
            }

            token = getNextToken(lex, false)

            variablesNode := declarationNode.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel + 2,
                text: "Array's variables",
            })

            token = parseExpressions(tree, variablesNode, token, lex, currentLevel + 4)

            if token.typ != itemRightDelim {
                parseErrorPrint(token, itemRightDelim)
            }

            return getNextToken(lex, false)
        } else {
            parseErrorPrint(token, itemChar)
        }
    }

    if token.val == "[" && token.typ == itemChar {
        indexNode := identifierNode.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel,
                text: "itemArraySize",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseSimpleExpression(tree, indexNode, token, lex, currentLevel + 2)

        if token.val != "]" || token.typ != itemChar {
            parseErrorPrint(token, itemChar)
        }

        token = getNextToken(lex, false)

        if token.typ == itemIntType || token.typ == itemStringType || token.typ == itemBoolType || token.typ == itemByteType {
            declarationNode.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemVariableType,
                level: currentLevel + 2,
                text: "Variable type",
                data: token.val,
            })
        } else {
            parseErrorPrint(token, itemVariableType)
        }

        return getNextToken(lex, false)
    }

    parseErrorPrint(token, itemVariableType)

    return token
}

func parseInstructionExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemIdentifier {
        expressionNode := node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel,
                text: "Expression",
        })
        childNode := expressionNode.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemIdentifier,
                level: currentLevel + 2,
                text: "Identifier",
                data: token.val,
        })

        token = parseExtendedFactor(tree, childNode, token, lex, currentLevel + 2)

        if token.typ == itemAssign {
            node.addChild(&AstTree{
                    key: time.Now().String(),
                    typ: itemAssign,
                    level: currentLevel,
                    text: "itemAssign",
                    data: token.val,
            })

            token = getNextToken(lex, false)

            token = parseExpression(tree, node, token, lex, currentLevel)
        }

        return token
    }

    parseErrorPrint(token, itemIdentifier)

    return token
}

func parseExtendedLogicalExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.typ == itemEqual || token.typ == itemGreater || token.typ == itemLower || token.typ == itemGreaterOrEqual || token.typ == itemLowerOrEqual || token.typ == itemNotEqual {
        token = parseComparison(tree, node, token, lex, currentLevel)
        token = parseSimpleExpression(tree, node, token, lex, currentLevel)
        token = parseExtendedLogicalExpression(tree, node, token, lex, currentLevel)
    }

    return token
}

func parseComparison(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.val == "==" && token.typ == itemEqual {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemEqual,
            level: currentLevel,
            text: "itemEqual",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    if token.val == "!=" && token.typ == itemNotEqual {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemNotEqual,
            level: currentLevel,
            text: "itemNotEqual",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    if token.val == ">" && token.typ == itemGreater {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemGreater,
            level: currentLevel,
            text: "itemGreater",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    if token.val == "<" && token.typ == itemLower {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemLower,
            level: currentLevel,
            text: "itemLower",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    if token.val == ">=" && token.typ == itemGreaterOrEqual {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemGreaterOrEqual,
            level: currentLevel,
            text: "itemGreaterOrEqual",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    if token.val == "<=" && token.typ == itemLowerOrEqual {
        node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemLowerOrEqual,
            level: currentLevel,
            text: "itemLowerOrEqual",
            data: token.val,
        })

        token = getNextToken(lex, false)
    }

    return token
}

func parseSimpleExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if (token.val == "-") && (token.typ == itemMinus) {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemMinus,
                level: currentLevel,
                text: "itemMinus",
                data: token.val,
        })

        token = getNextToken(lex, false)
    }

    token = parseTerm(tree, node, token, lex, currentLevel)
    token = parseExtendedSimpleExpression(tree, node, token, lex, currentLevel)

    return token
}

func parseExtendedSimpleExpression(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if (token.val == "-") && (token.typ == itemMinus) {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemMinus,
                level: currentLevel,
                text: "itemMinus",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseTerm(tree, node, token, lex, currentLevel)
    }

    if (token.val == "+") && (token.typ == itemPlus) {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemPlus,
                level: currentLevel,
                text: "itemPlus",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseTerm(tree, node, token, lex, currentLevel)
    }

    return token
}

func parseTerm(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    token = parseFactor(tree, node, token, lex, currentLevel)
    token = parseExtendedTerm(tree, node, token, lex, currentLevel)

    return token
}

func parseExtendedTerm(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.val == "*" && token.typ == itemMupltiply {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemMupltiply,
                level: currentLevel,
                text: "itemMupltiply",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseFactor(tree, node, token, lex, currentLevel)
        token = parseExtendedTerm(tree, node, token, lex, currentLevel)
    }

    if token.val == "/" && token.typ == itemDivide {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemDivide,
                level: currentLevel,
                text: "itemDivide",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseFactor(tree, node, token, lex, currentLevel)
        token = parseExtendedTerm(tree, node, token, lex, currentLevel)
    }

    if token.val == "%" {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemRest,
                level: currentLevel,
                text: "itemRest",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseFactor(tree, node, token, lex, currentLevel)
        token = parseExtendedTerm(tree, node, token, lex, currentLevel)
    }

    return token
}

func parseFactor(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.val == "-" && token.typ == itemMinus {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemChar,
                level: currentLevel,
                text: "itemNot",
        })

        token = getNextToken(lex, false)

        return token
    }

    if token.typ == itemIdentifier {
        childNode := node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemIdentifier,
                level: currentLevel,
                text: "Identifier",
                data: token.val,
        })

        token = parseExtendedFactor(tree, childNode, token, lex, currentLevel)

        return token
    }

    if token.typ == itemBool {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemBool,
                level: currentLevel,
                text: "Boolean",
                data: token.val,
        })

        token = getNextToken(lex, false)

        return token
    }

    if token.typ == itemNumber {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNumber,
                level: currentLevel,
                text: "Number",
                data: token.val,
        })

        token = getNextToken(lex, false)

        return token
    }

    if token.typ == itemString {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemString,
                level: currentLevel,
                text: "itemString",
                data: token.val,
        })

        token = getNextToken(lex, false)

        return token
    }

    if token.val == "(" && token.typ == itemLeftParen {
        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemLeftParen,
                level: currentLevel,
                text: "itemLeftParen",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseExpression(tree, node, token, lex, currentLevel)

        if token.val != ")" || token.typ != itemRightParen {
            parseErrorPrint(token, itemRightParen)
        }

        node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemRightParen,
                level: currentLevel,
                text: "itemRightParen",
                data: token.val,
        })

        token = getNextToken(lex, false)

        return token
    }

    fmt.Println(token.val)
    fmt.Println(valuesTranslations[int(token.typ)])

    token = getNextToken(lex, false)

//// NOTE: Сделать вывод ошибок
//
    parseErrorPrint(token, itemNumber)

    return token
}

func parseExtendedFactor(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    token = getNextToken(lex, false)

    if token.val == "[" && token.typ == itemChar {
        indexNode := node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel,
                text: "itemIndex",
                data: token.val,
        })

        token = getNextToken(lex, false)

        token = parseSimpleExpression(tree, indexNode, token, lex, currentLevel + 2)

        if token.val != "]" || token.typ != itemChar {
            parseErrorPrint(token, itemChar)
        }

        token = getNextToken(lex, false)
    }

    if token.val == "(" && token.typ == itemLeftParen {
        functionParametersNode := node.addChild(&AstTree{
                key: time.Now().String(),
                typ: itemNode,
                level: currentLevel + 4,
                text: "Function parameters",
        })

        token = getNextToken(lex, false)

        token = functionParameter(tree, functionParametersNode, token, lex, currentLevel + 6)

        if token.val != ")" || token.typ != itemRightParen {
            parseErrorPrint(token, itemRightParen)
        }

        token = getNextToken(lex, false)
    }

    if token.typ == itemFunction {
        childNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemFunction,
            level: currentLevel + 2,
            text: "Function of identifier",
            data: token.val,
        })

        token = parseExtendedFactor(tree, childNode, token, lex, currentLevel)
    }

    if token.typ == itemField {
        childNode := node.addChild(&AstTree{
            key: time.Now().String(),
            typ: itemField,
            level: currentLevel + 2,
            text: "Field of identifier",
            data: token.val,
        })

        token = parseExtendedFactor(tree, childNode, token, lex, currentLevel)
    }

    return token
}

func functionParameter(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item {
    if token.val == ")" && token.typ == itemRightParen {
        return token
    }

    return parseExpressions(tree, node, token, lex, currentLevel)
}

func parseExpressions(tree *AstTree, node *AstTree, token *item, lex * lexer, currentLevel int) *item  {
    token = parseExpression(tree, node, token, lex, currentLevel)

    if token.val == "," && token.typ == itemChar {
        token = getNextToken(lex, false)
        token = parseExpressions(tree, node, token, lex, currentLevel)
    }

    return token
}

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

func parseErrorPrint(token *item, item itemType) {
    fmt.Println(token.line, ":", token.pos,"syntax error: unexpected", token.val, ", expecting", valuesTranslations[int(item)])
  // fmt.Println("There should be", valuesTranslations[int(item)], "<", token.pos, ">", "line:", token.line, "got:", valuesTranslations[int(token.typ)])
  os.Exit(1)
}
