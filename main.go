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

// im
func ErrOut(err error) string {
	return err.Error()
}

func converter(doc string) string {
	d, err := html.Parse(strings.NewReader(doc))
	if err != nil {

		return ""
	}

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

func getAttrMap(t *html.Node) map[string]string {
	mapAttr := make(map[string]string)
	for _, a := range t.Attr {
		mapAttr[a.Key] = a.Val
	}
	return mapAttr
}

func parser(fileOut *os.File, root *html.Node) {
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

func parseNode(root *html.Node) (bool, []byte) {

	tt := root.NextSibling

	var bits []byte
	switch {

	case tt.Type == html.ErrorNode:
		return false, bits

	case tt.Type == html.ElementNode:

		// t := html.Parse
		tmap := getAttrMap(tt)

		class := tmap["class"]
		switch {
		case strings.Contains(class, "report-master-meta"):
			bits = parseMasterConfig(tt, tmap)
		case strings.Contains(class, "page"):
			bits = parsePage(tt, tmap)
		}
	}
	return true, bits
}

func parseMasterConfig(element *html.Node, attrMap AttrMap) []byte {
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

func parsePage(element *html.Node, attrMap AttrMap) []byte {
	// tt := element.Next()
	// var bits []byte
	// switch {

	// case tt == html.ErrorToken:
	// 	return false, bits

	// case tt == html.StartTagToken:
	// 	t := root.Token()
	// 	tmap := getAttrMap(t)

	// 	class := tmap["class"]
	// 	switch {
	// 	case strings.Contains(class, "report-master-meta"):
	// 		bits = parseMasterConfig(t, tmap)
	// 	case strings.Contains(class, "page"):
	// 		bits = parsePage(t, tmap)
	// 	}
	// }
	// return true, bits

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

	// pageHeader := PageHeader{
	// 	height:          0,
	// 	backgroundColor: "",
	// }
	// viz := Viz{
	// 	component: Component{
	// 		Height: 0,
	// 		Width:  0,
	// 		x:      0,
	// 		y:      0,
	// 		id:     "",
	// 	},
	// 	image: Base64Image{
	// 		height: 0,
	// 		width:  0,
	// 		image:  "",
	// 	},
	// }
	// arr := []interface{}{pageHeader, viz}
	// var g = reflect.TypeOf(arr[0])

	// g.

	// get children
	doc, er := html.Parse(r)
	if er != nil {
		panic(er)
	}

	parser(fo, doc)

}
