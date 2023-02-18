package dochandler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/info"
)

func MdDocHandler(doc *openapi3.T, info *info.Info) http.HandlerFunc {
	var text string
	var retErr error
	var once sync.Once
	return func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			mddoc := docgen.Generate(doc, info)
			mddoc.SkipMetadata = true
			buf := new(strings.Builder)
			if err := docgen.Docgen(buf, mddoc); err != nil {
				retErr = err
				return
			}
			text = strings.ReplaceAll(buf.String(), "```", "~~~")
		})

		if retErr != nil {
			log.Printf("[WARN]  !! %+v", retErr)
			fmt.Fprintf(w, "!! %v", retErr.Error())
			return
		}

		title := fmt.Sprintf("%s (%s)", doc.Info.Title, doc.Info.Version)
		fmt.Fprintf(w, MDDOC_TEMPLATE, title, text)
	}
}

const MDDOC_TEMPLATE = `<!DOCTYPE html>
<html lang="ja">
<meta charset="UTF-8">
<title>%s</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.1.0/github-markdown.min.css">
<style>
	.markdown-body {
		box-sizing: border-box;
		min-width: 200px;
		max-width: 980px;
		margin: 0 auto;
		padding: 45px;
	}
	@media (max-width: 767px) {
		.markdown-body {
			padding: 15px;
		}
	}
</style>
<script defer type="module">
import React from "https://esm.sh/react@18.2.0?bundle";
import ReactMarkdown from "https://esm.sh/react-markdown@8.0.5?bundle";
import remarkGfm from "https://esm.sh/remark-gfm@3.0.1?bundle";
import { createRoot } from "https://esm.sh/react-dom@18.2.0/client?bundle";
import { Prism as SyntaxHighlighter } from "https://esm.sh/react-syntax-highlighter@15.5.0?bundle";

const flatten = (text, child) => {
  return typeof child === "string"
    ? text + child
    : React.Children.toArray(child.props.children).reduce(flatten, text);
};

const HeadingRenderer = (props) => {
  const children = React.Children.toArray(props.children);
  const text = children.reduce(flatten, "");
  const slug = text.toLowerCase().replace(/[{\/\.}]+/g, "").replace(
    /[ \t]+/g,
    "-",
  );

  const a = React.createElement("a", {
    "class": "x-anchor",
    "aria-hidden": "true",
    href: "#" + slug,
  }, props.children);
  return React.createElement("h" + props.level, {
    id: slug,
    tabindex: "-1",
    dir: "auto",
  }, a);
};

const languageRegex = /language-(\w+)/;

const SyntaxHighlightRenderer = (
    { node, inline, className, children, ...props },
  ) => {
    const match = languageRegex.exec(className || "");
    return !inline && match
      ? SyntaxHighlighter(
        {
          language: match[1],
          children: children,
          ...props,
        },
      )
      : React.createElement("code", { className: className, ...props }, children);
  };

  const SyntaxHighlightRendererForPre = (props) => {
    const child = props.children[0].props;
	const match = languageRegex.exec(child.className || "");
	return !props.inline && match ? SyntaxHighlightRenderer({...child}) : React.createElement("pre", {}, props.children);
  };
    
const text = document.getElementById("mdtext").innerText;

const domNode = document.getElementById("mdbody");
const root = createRoot(domNode);
root.render(
  ReactMarkdown({
    children: text,
    components: {
      "h1": HeadingRenderer,
      "h2": HeadingRenderer,
      "h3": HeadingRenderer,
      "pre": SyntaxHighlightRendererForPre,
    },
    remarkPlugins: [remarkGfm],
  }),
);
</script>
<body>
<x-markdown id="mdtext" style="display:none;">%s</x-markdown>
<article id="mdbody" class="markdown-body">loading...</article>
</body>
<html>
`
