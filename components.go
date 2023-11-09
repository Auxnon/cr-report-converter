package main

type Component struct {
	Height int    `kdl:"height"`
	Width  int    `kdl:"width"`
	X      int    `kdl:"x"`
	Y      int    `kdl:"y"`
	Id     string `kdl:"id"`
}

type Base64Image struct {
	Height int
	Width  int
	Data   string
}

type Viz struct {
	Component Component   `kdl:"component"`
	Image     Base64Image `kdl:"image"`
}

const TABLE_SECTION_STYLES = "position: absolute; overflow: visible;"
const CELL_STYLES = "padding: 2px 4px;"

type TableSection struct {
	Component Component `kdl:",props"`
	// Contents  []interface{}
}

type Table struct {
	component      Component
	borderColor    string
	fontSize       int
	primaryColor   string
	secondaryColor string
	headerColor    string
	widths         []int
	heights        []int
	schema         int
}

type TableSchema struct {
	id         int
	cellStyles []string
	cellValues []string
}

// type Row struct {
