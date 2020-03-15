# Useage

```
func main() {
	var inputString = "(add 2 (subtract 4 2))" // (add 200 (subtract 60 3))
	tokens, err := compiler.Tokenizer(inputString)
	if err != nil {
		log.Fatal("tokens errorrr: ", err)
	}
	ast := compiler.Parser(tokens)
	newAst := compiler.Transformer(&ast)
	output, err := compiler.CodeGenerator(&newAst)
	if err != nil {
		log.Fatal("codeGenerator meet error : ", err)
	}

	//output := compiler.Compiler("(add 2 (subtract 4 2))")
	fmt.Println("$$$$$$$$$$$$$$$$$", output)
}
```

# Testing

        go test .