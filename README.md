# Useage

```
func main() {
	var inputString = "(add 2 (subtract 4 2))" // (add 200 (subtract 60 3))
	tokens, err := compiler.Tokenizer(inputString)
	if err != nil {
		log.Fatal("tokens error: ", err)
	}
	ast := compiler.Parser(tokens)
	newAst := compiler.Transformer(&ast)
	output, err := compiler.CodeGenerator(&newAst)
	if err != nil {
		log.Fatal("codeGenerator meet error : ", err)
	}

	//output := compiler.Compiler("(add 2 (subtract 4 2))")
	fmt.Println("output: ", output)
}
```

# Testing

        go test .
        
---

---

---
        
另有一篇Go 的compiler 翻译[repository](https://github.com/hazbo/the-super-tiny-compiler)

我在全局中，使用指针式上下传递；

那篇在Transformer 中，使用直观易懂的方式做了个取巧，不过我觉得还是俺这种更承接点吧。
