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

// Excessive but slightly faster then doing a full attr map lookup for class name, use for one offs
func hasClass(t *html.Node, className string) bool {
	for _, a := range t.Attr {
		if a.Key == "class" {
			return strings.Contains(a.Val, className)
		}
	}
	return false
}

func getStyleMap(t *html.Node, attrMap AttrMap) map[string]string {
	styleMap := make(map[string]string)
	styleStr := attrMap["style"]
	splitStr := strings.Split(styleStr, ";")
	for _, s := range splitStr {
		split2 := strings.Split(s, ":")
		if len(split2) == 2 {
			styleMap[strings.TrimSpace(split2[0])] = strings.TrimSpace(split2[1])
		}
	}
	return styleMap
}

func or(styleMap AttrMap, checks ...string) string {
	for _, check := range checks {
		if styleMap[check] != "" {
			return styleMap[check]
		}
	}
	return ""
}

// unwrap multiple return to non error type, ignoring error
func unwrap[T any](x T, e error) T {
	return x
}

func parseInt(s string) int {
	s, _ = strings.CutSuffix(s, "px")
	i, _ := strconv.Atoi(s)
	return i
}

func asString(bytes []byte, e error) string {
	return string(bytes)
}

func getInnerBody(element *html.Node) *html.Node {
	if element.Data == "" || element.Data == "html" {
		return getInnerBody(element.FirstChild)
	}
	if element.Data == "body" {
		return element
	} else if element.Data == "head" {
		return getInnerBody(element.NextSibling)
	}
	return nil
}

func parser(fileOut *os.File, root *html.Node) {
	for {
		running, bits := parseNode(root)
		root = root.NextSibling
		if !running || root == nil {
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

	tt := root

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
	// set root to next sibling
	return true, bits
}

func parseMasterConfig(element *html.Node, attrMap AttrMap) []byte {
	editorConfigStr := attrMap["data-editor-config"]
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

func parsePage(element *html.Node, attrMap AttrMap) []byte {
	tt := element
	id := attrMap["id"]

	styleMap := getStyleMap(tt, attrMap)

	pageNode := Page{
		Height:          parseInt(or(styleMap, "--page-height", "height")),
		Width:           parseInt(or(styleMap, "--page-width", "width")),
		Id:              id,
		BackgroundColor: or(styleMap, "background-color", "background"),
		TableSection:    []TableSection{},
		// Viz:             Viz{},
		Contents: []interface{}{},
	}

	t := TableSection{
		Component: Component{
			Height: 0,
			Width:  0,
			X:      0,
			Y:      0,
			Id:     "",
		},
	}

	pageNode.TableSection = append(pageNode.TableSection, t)

	child := getFirstNonTextChild(tt)
	for child != nil {
		walkContent(child, &pageNode)
		// if item != nil {
		// 	pageNode.Contents = append(pageNode.Contents, item)
		// }
		child = nextNonTextSibling(child)
	}

	out, er := kdl.Marshal(struct {
		Page Page `kdl:"page,child"`
	}{
		Page: pageNode,
	})

	if er != nil {
		fmt.Println("parseMasterConfig:", er)
	}
	return out

}

func walkContent(element *html.Node, page *Page) {
	tt := element
	if tt.Type == html.TextNode {
		return
	}
	if tt.Type == html.ElementNode {
		tmap := getAttrMap(tt)
		compType := tmap["data-component-type"]
		switch compType {
		case "live-viz":
			parseViz(tt, tmap, &page.Contents)
			// if viz != (Viz{}) {
			// 	// vv := struct {
			// 	// 	Viz Viz `kdl:"viz"`
			// 	// }{
			// 	// 	Viz: viz,
			// 	// }
			// 	output := unwrap(kdl.Marshal(struct {
			// 		Viz Viz `kdl:"viz"`
			// 	}{
			// 		Viz: viz,
			// 	}))

			// 	page.Contents = append(page.Contents, output...)

			// }

		case "table-section":
			parseTable(tt, tmap, &page.Contents)

		}
	}
	// return
}

// FirstChild but skips text nodes
func getFirstNonTextChild(element *html.Node) *html.Node {
	child := element.FirstChild
	for child != nil {
		if child.Type != html.TextNode {
			return child
		}
		child = child.NextSibling
	}
	return nil
}

// NextSibling but skips text nodes
func nextNonTextSibling(element *html.Node) *html.Node {
	sibling := element.NextSibling
	for sibling != nil {
		if sibling.Type != html.TextNode {
			return sibling
		}
		sibling = sibling.NextSibling
	}
	return nil
}

func parseViz(tt *html.Node, attrMap AttrMap, contents *[]interface{}) {
	styleMap := getStyleMap(tt, attrMap)
	inner := getFirstNonTextChild(tt)

	if inner == nil {
		return
	}
	if hasClass(inner, "viz-content") {
		inner = getFirstNonTextChild(inner)
		if inner == nil {
			return
		}
		if hasClass(inner, "thumbnail-preview") {
			inner = nextNonTextSibling(inner)
			if inner == nil {
				return
			}
		}
	}

	if inner.Data != "img" {
		return
	}

	innerMap := getAttrMap(inner)
	data := innerMap["src"]

	// splitData := strings.Split(data, ",")
	// v, e := base64.StdEncoding.DecodeString(splitData[1])
	// if e == nil {
	// fmt.Println(v)
	// }
	data = "eldata"
	viz := Viz{
		Component: Component{
			Height: parseInt(styleMap["height"]),
			Width:  parseInt(styleMap["width"]),
			X:      parseInt(styleMap["left"]),
			Y:      parseInt(styleMap["top"]),
			Id:     attrMap["id"],
		},
		Image: Base64Image{
			Height: parseInt(innerMap["height"]),
			Width:  parseInt(innerMap["width"]),
			Data:   data,
		},
	}
	fmt.Printf("%+v\n", viz)

	// *contents = append(*contents, struct {
	// 	Viz Viz `kdl:"viz,child"`
	// }{
	// 	Viz: viz,
	// })

	*contents = append(*contents, struct {
		Viz Viz
	}{
		Viz: viz,
	})

}

func parseTable(tt *html.Node, attrMap AttrMap, tables *[]interface{}) {

	section := TableSection{
		Component: Component{
			Height: 0,
			Width:  0,
			X:      0,
			Y:      0,
			Id:     "",
		},
		// Contents: []interface{}{},
	}

	*tables = append(*tables, struct {
		TableSection TableSection `kdl:"table-section,child"`
	}{
		TableSection: section,
	})

}

func main() {
	// str := base64.StdEncoding.EncodeToString([]byte("Hello, playground"))
	// fmt.Println("hi", str)

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
	doc = getInnerBody(doc)
	if doc == nil {
		panic("no body")
	}

	parser(fo, doc.FirstChild)

}
