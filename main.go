package tabulator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type AttrMap map[string]string

//export In
func In(doc string) string {
	return converter(doc)
}

func converter(doc string) string {
	d := html.NewTokenizer(strings.NewReader(doc))
	return parser(d)

}

// func getClass(t html.Token) string {
// 	return t.Attr["class"].Val
// }

// func getAttr(t html.Token, attr string) string {
// 	return t.Attr[attr].Val
// }

func getAttrMap(t html.Token) map[string]string {
	mapAttr := make(map[string]string)
	for _, a := range t.Attr {
		mapAttr[a.Key] = a.Val
	}
	return mapAttr
}

// func hasClass(t html.Token, class string) bool {
// 	s := getClass(t)
// 	return strings.Contains(s, class)
// }

func parser(root *html.Tokenizer) string {
	for {

		tt := root.Next()
		switch {

		case tt == html.ErrorToken:
			return ""

		case tt == html.StartTagToken:
			t := root.Token()
			tmap := getAttrMap(t)

			class := tmap["class"]
			switch {
			case strings.Contains(class, "report-master-meta"):
				return parseMasterConfig(t, tmap)
			case strings.Contains(class, "page"):
				return parsePage(t, tmap)
			}
		}
	}
	// return ""
}

func parseMasterConfig(element html.Token, attrMap AttrMap) string {
	editorConfigStr := attrMap["data-editor-config"]
	var editorConfig EditorConfig
	if editorConfigStr != "" {
		editorConfigStr = html.UnescapeString(editorConfigStr)
		err := json.Unmarshal([]byte(editorConfigStr), &editorConfig)
		if err != nil {
			fmt.Println(err)
		}
	}
	// renderTime := attrMap["data-render-time"]
	// platformWaitTime := attrMap["data-platform-wait-time"]
	return ""
}

func parsePage(element html.Token, attrMap AttrMap) string {
	return ""
}

func main() {
	// load test
	r, _ := os.Open("test.html")
	defer r.Close()

	out := parser(html.NewTokenizer(r))
	println(out)
	// z := html.NewTokenizer(r)
	// fmt.Println("yup")
	// html.EscapeString(os.Stdout, "<script>alert('you have been pwned')</script>")
}
