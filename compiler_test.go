package compiler

import (
	"fmt"
	"reflect"
	"testing"
)

const testInput = "(add 200 (subtract 60 3))"
const output = "add(200, subtract(60, 3));"

var testTokens = []Token{
	{
		Type:  "paren",
		Value: "(",
	},
	{
		Type:  "name",
		Value: "add",
	},
	{
		Type:  "number",
		Value: "200",
	},
	{
		Type:  "paren",
		Value: "(",
	},
	{
		Type:  "name",
		Value: "subtract",
	},
	{
		Type:  "number",
		Value: "60",
	},
	{
		Type:  "number",
		Value: "3",
	},
	{
		Type:  "paren",
		Value: ")",
	},
	{
		Type:  "paren",
		Value: ")",
	},
}

var testAst = ASTExpression{
	Type: "Program",
	Body: []ASTExpression{
		ASTExpression{
			Type: "CallExpression",
			name: "add",
			params: []ASTExpression{
				ASTExpression{
					Type:  "NumberLiteral",
					Value: "200",
				},
				ASTExpression{
					Type: "CallExpression",
					name: "subtract",
					params: []ASTExpression{
						ASTExpression{
							Type:  "NumberLiteral",
							Value: "60",
						},
						ASTExpression{
							Type:  "NumberLiteral",
							Value: "3",
						},
					},
				},
			},
		},
	},
}

func TestTokenizer(t *testing.T) {
	fmt.Println("going to have a test!!!")
	results, err := Tokenizer(testInput)
	if err != nil {
		t.Error("tokenizer is not working")
	}
	if !reflect.DeepEqual(results, testTokens) {
		t.Error("\nWan: ", testTokens, "\nGot: ", results)
	}
}

func TestParser(t *testing.T) {
	result := Parser(testTokens)
	if result.Type != "Program" {
		t.Error("err program")
	}
	if len(result.Body) < 1 {
		t.Error("program.body is empty")
	}
	ast := result.Body[0]
	if ast.Type != "CallExpression" {
		t.Error("program got not CallExpression")
	}
}

func TestCodeGenerator(t *testing.T) {
	result := Compiler(testInput)
	if result != output {
		t.Error("\nWan: ", output, "\nGot: ", result)
	}
}

func BenchmarkTokenizer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Tokenizer(testInput)
	}
}