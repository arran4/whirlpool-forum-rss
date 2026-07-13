package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	md, err := os.ReadFile("readme.md")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	htmlOutput := markdown.Render(doc, renderer)

	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<title>Whirlpool Forum RSS</title>
<meta charset="utf-8">
</head>
<body>
%s
</body>
</html>`, htmlOutput)

	err = os.WriteFile("public/index.html", []byte(fullHTML), 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
	fmt.Println("Created public/index.html")
}
