package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	kdl "github.com/sblinch/kdl-go"
	"golang.org/x/net/html"
)

type AttrMap map[string]string

type MastConfig struct {
	EditorConfig     EditorConfig
	RenderTime       int
	PlatformWaitTime int
}

//export In
func In(doc string) string {
	return converter(doc)
}

func converter(doc string) string {
	d := html.NewTokenizer(strings.NewReader(doc))
	str := ""

	for {
		running, bits := parseNode(d)
		if !running {
			return str
		}
		if len(bits) > 0 {
			str += string(bits)
		}
	}
}

func getAttrMap(t html.Token) map[string]string {
	mapAttr := make(map[string]string)
	for _, a := range t.Attr {
		mapAttr[a.Key] = a.Val
	}
	return mapAttr
}

func parser(fileOut *os.File, root *html.Tokenizer) {
	for {
		running, bits := parseNode(root)
		if !running {
			return
		}
		if len(bits) > 0 {
			if _, err := fileOut.Write(bits); err != nil {
				panic(err)
			}
		}
	}
}

func parseNode(root *html.Tokenizer) (bool, []byte) {
	tt := root.Next()
	var bits []byte
	switch {

	case tt == html.ErrorToken:
		return false, bits

	case tt == html.StartTagToken:
		t := root.Token()
		tmap := getAttrMap(t)

		class := tmap["class"]
		switch {
		case strings.Contains(class, "report-master-meta"):
			bits = parseMasterConfig(t, tmap)
		case strings.Contains(class, "page"):
			bits = parsePage(t, tmap)
		}
	}
	return true, bits
}

func parseMasterConfig(element html.Token, attrMap AttrMap) []byte {
	editorConfigStr := attrMap["data-editor-config"]
	fmt.Println(editorConfigStr, "\n")
	var editorConfig EditorConfig

	if editorConfigStr != "" {
		editorConfigStr, err := url.QueryUnescape(editorConfigStr)
		if err != nil {
			fmt.Println("parseMasterConfig:", err)
		}
		fmt.Println(editorConfigStr)

		if err = json.Unmarshal([]byte(editorConfigStr), &editorConfig); err != nil {
			fmt.Println("parseMasterConfig:", err)
		}
	}

	master := MastConfig{
		EditorConfig:     editorConfig,
		RenderTime:       unwrap(strconv.Atoi(attrMap["data-render-time"])),
		PlatformWaitTime: unwrap(strconv.Atoi(attrMap["data-platform-wait-time"])),
	}
	out, er := kdl.Marshal(struct {
		MastConfig MastConfig `kdl:"master"`
	}{
		MastConfig: master,
	})
	if er != nil {
		fmt.Println("parseMasterConfig:", er)
	}
	return out
}

// unwrap multiple return to non error type, ignoring error
func unwrap[T any](x T, e error) T {
	return x
}

func parsePage(element html.Token, attrMap AttrMap) []byte {
	return []byte{}
}

func main() {
	r, _ := os.Open("test.html")
	defer r.Close()

	fo, err := os.Create("output.kdl")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	parser(fo, html.NewTokenizer(r))

}
