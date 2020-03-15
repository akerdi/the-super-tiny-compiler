package compiler

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Token struct {
	Type string
	Value string
}

type TypeName struct {
	Type string
	name string
}

type ASTExpression struct {
	Type string
	Value string
	name string
	Body []ASTExpression
	Context *[]ASTExpression
	params []ASTExpression
	callee *ASTExpression
	arguments *[]ASTExpression
	expression *ASTExpression
}

var validNumber = regexp.MustCompile(`[0-9]`)
var validFunction = regexp.MustCompile(`[a-z]`)

func Tokenizer(input string) ([]Token, error) {
	current := 0
	var tokens []Token
	inputBytes := []byte(input)
	for {
		if current >= len(input) {
			break
		}
		char := inputBytes[current]
		if char == '(' {
			_token := Token{
				Type:  "paren",
				Value: "(",
			}
			tokens = append(tokens, _token)
			current++
			continue
		}

		if char == ')' {
			_token := Token{
				Type:  "paren",
				Value: ")",
			}
			tokens = append(tokens, _token)
			current++
			continue
		}
		// whitespace
		if char == ' ' {
			current++
			continue
		}
		if validNumber.MatchString(string(char)) {
			var _value string
			// while 判断是否是number
			for {
				_str := string(char)
				if validNumber.MatchString(_str) {
					_value += string(_str)
					current++
					char = inputBytes[current]
				} else {
					break
				}
			}
			_token := Token{
				Type:  "number",
				Value: _value,
			}
			tokens = append(tokens, _token)
			continue
		}

		if char == '"' {
			var _value string
			current++
			char = inputBytes[current]
			// while 判断是不是string
			for {
				if char != '"' {
					_value += string(char)
					current++
					char = inputBytes[current]
				} else {
					break
				}
			}

			// string 已经判断完毕
			current++
			char = inputBytes[current]
			_token := Token{
				Type:  "string",
				Value: _value,
			}
			tokens = append(tokens, _token)
			continue
		}

		// 接下来判断方法名 add 2 4 中的add
		str := string(char)
		if validFunction.MatchString(str) {
			var _value string
			for {
				_str := string(char)
				if validFunction.MatchString(_str) {
					_value += _str
					current++
					char = inputBytes[current]
				} else {
					break
				}
			}
			_token := Token{
				Type:  "name",
				Value: _value,
			}
			tokens = append(tokens, _token)
			continue
		}
		// 其他类型不支持，直接报错
		return nil, errors.New(fmt.Sprintf("I dont know what you talk about: %x", char))
	}
	return tokens, nil
}

////////////////////////////
////////////////////////////
////////////////////////////

func walk(tokens *[]Token, current *int) (ASTExpression, error) {
	var tToken = *tokens
	_token := tToken[*current]
	if _token.Type == "number" {
		*current++
		return ASTExpression{
			Type:  "NumberLiteral",
			Value: _token.Value,
			name:   "",
			params: nil,
			Context: &[]ASTExpression{},
		}, nil
	}
	if _token.Type == "string" {
		*current++
		return ASTExpression{
			Type:  "StringLiteral",
			Value: _token.Value,
			name:   "",
			params: nil,
			Context: &[]ASTExpression{},
		}, nil
	}
	if _token.Type == "paren" && _token.Value == "(" {
		*current++
		_token = tToken[*current]
		var nodeParams []ASTExpression
		node := ASTExpression{
			Type: "CallExpression",
			name:   _token.Value,
			params: nodeParams,
			Context: &[]ASTExpression{},
		}
		*current++
		_token = tToken[*current]
		for {
			if _token.Type != "paren" || (_token.Type == "paren" && _token.Value != ")") {
				_expression, err := walk(tokens, current)
				if err != nil {
					log.Fatal("[compiler.walk] meet error ", err)
					continue
				}
				node.params = append(node.params, _expression)
				_token = tToken[*current]
			} else {
				break
			}
		}
		*current++
		return node, nil
	}
	return ASTExpression{}, errors.New(fmt.Sprintf("[compiler.walk] dont recognized type %s", _token.Type))
}

func Parser(tokens []Token) ASTExpression {
	var current int = 0
	var _astBody []ASTExpression
	ast := ASTExpression{
		Type: "Program",
		Body: _astBody,
	}
	for {
		if current < len(tokens) {
			_expression, err := walk(&tokens, &current)
			if err != nil {
				log.Fatal("[compiler.walk] meet error ", err)
				continue
			}
			ast.Body = append(ast.Body, _expression)
		} else {
			break
		}
	}
	return ast
}

////////////////////////////
////////////////////////////
////////////////////////////

func traverseNode(node *ASTExpression, parent *ASTExpression, visitor Visitor)  {
	var tNode = *node
	var methods VisitorLiteral
	switch node.Type {
	case "NumberLiteral":
		methods = visitor.NumberLiteral
		break
	case "StringLiteral":
		methods = visitor.StringLiteral
		break
	case "CallExpression":
		methods = visitor.CallExpression
		break
	}
	if methods.enter != nil {
		methods.enter(node, parent)
	}

	switch tNode.Type {
	case	"Program":
		traverseArray(&tNode.Body, node, visitor)
		break
	case "CallExpression":
		traverseArray(&tNode.params, node, visitor)
		fmt.Println("~~~~~~~`", node.Context)
		break
	case "NumberLiteral":
	case "StringLiteral":
		break
	default:
		log.Fatal("[compiler.traverseNode] default error no type", node.Type)
	}
	if methods.exit != nil {
		methods.exit(node, parent)
	}
}

func traverseArray(array *[]ASTExpression, parent *ASTExpression, visitor Visitor)  {
	var tArray = *array
	for _, child := range tArray {
		//fmt.Println("  2222 ", child)
		traverseNode(&child, parent, visitor)
	}
}

func traverser(ast *ASTExpression, visitor Visitor)  {
	traverseNode(ast, &ASTExpression{}, visitor)
}


////////////////////////////
////////////////////////////
////////////////////////////

func Transformer(ast *ASTExpression) ASTExpression {
	var _newAstBody []ASTExpression = []ASTExpression{}
	newAst := ASTExpression{
		Type:    "Program",
		Body:    _newAstBody,
	}
	ast.Context = &newAst.Body
	traverser(ast, Visitor{
		NumberLiteral:  VisitorLiteral{enter:NumberEnter},
		StringLiteral:  VisitorLiteral{enter:StringEnter},
		CallExpression: VisitorLiteral{enter:CallEnter},
	})
	return newAst
}

func NumberEnter(node *ASTExpression, parent *ASTExpression) {
	var tNode = *node
	var tParent = *parent
	*tParent.Context = append(*tParent.Context, ASTExpression{
		Type: "NumberLiteral",
		Value: tNode.Value,
	})
}
func StringEnter(node *ASTExpression, parent *ASTExpression) {
	var tNode = *node
	var tParent = *parent
	*tParent.Context = append(*tParent.Context, ASTExpression{
		Type: "StringLiteral",
		Value: tNode.Value,
	})
}
func CallEnter(node *ASTExpression, parent *ASTExpression) {
	var tNode = *node
	var tParent = *parent
	_expresssion := ASTExpression{
		Type: "CallExpression",
		callee:    &ASTExpression{
			Type: "Identifier",
			name: tNode.name,
		},
		arguments: &[]ASTExpression{},
	}
	// 这里很重要。
	// 这里的想法是，让tNode 当前的node的Context 数组和接下来的_expression.arguments 绑定
	_expresssion.arguments = tNode.Context
	if parent.Type != "CallExpression" {
		newExpresssion := ASTExpression{
			Type:  "ExpressionStatement",
			expression: &ASTExpression{
				Type:       _expresssion.Type,
				callee:     _expresssion.callee,
				arguments: 	_expresssion.arguments,
			},
		}
		_expresssion = newExpresssion
	}
	*tParent.Context = append(*tParent.Context, _expresssion)
}

type VisitorLiteral struct {
	enter func(*ASTExpression, *ASTExpression)
	exit func(*ASTExpression, *ASTExpression)
}

////////////////////////////
////////////////////////////
////////////////////////////

func CodeGenerator(node *ASTExpression) (res string, err error)  {
	err = nil
	var tNode =*node
	switch tNode.Type {
	case	"Program":
		var stringArray []string
		for _, body := range tNode.Body {
			resultStr, err := CodeGenerator(&body)
			if err != nil {
				log.Fatal("[compiler.codeGenerator] program error ", err)
			}
			stringArray = append(stringArray, resultStr)
		}
		res = strings.Join(stringArray, "\n")
		return res, nil
	case "ExpressionStatement":
		resultStr, err := CodeGenerator(tNode.expression)
		if err != nil {
			log.Fatal("[compiler.codeGenerator] ExpressionStatement error ", err)
		}
		res = resultStr + ";"
		return res, nil
	case "CallExpression":
		resultStr, err := CodeGenerator(tNode.callee)
		if err != nil {
			log.Fatal("[compiler.codeGenerator] ExpressionStatement error ", err)
		}
		var stringArray []string
		for _, argument := range *tNode.arguments {
			resultStr, err := CodeGenerator(&argument)
			if err != nil {
				log.Fatal("[compiler.codeGenerator] program error ", err)
			}
			stringArray = append(stringArray, resultStr)
		}
		res = resultStr + "(" + strings.Join(stringArray, ", ") + ")"
		return res, nil
	case "Identifier":
		res = tNode.name
		return res, nil
	case "NumberLiteral":
		res = tNode.Value
		return res, nil
	case "StringLiteral":
		res = "\"" + tNode.Value + "\""
		return res, nil
	default:
		log.Fatal("[compiler.codeGenerator] default error not such type: ", tNode.Type)
	}
	return
}

func Compiler(input string) string {
	tokens, err := Tokenizer(input)
	if err != nil {
		log.Fatal("Compiler tokenizer error: ", err)
	}
	ast := Parser(tokens)
	newAst := Transformer(&ast)
	output, err := CodeGenerator(&newAst)
	if err != nil {
		log.Fatal("Compiler codeGenerator error: ", err)
	}
	return output
}

type Visitor struct {
	NumberLiteral VisitorLiteral
	StringLiteral VisitorLiteral
	CallExpression VisitorLiteral
}
