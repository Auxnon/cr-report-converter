package main

type Component struct {
	Height int
	Width  int
	x      int
	y      int
	id     string
}

type Base64Image struct {
	height int
	width  int
	image  string
}

type Viz struct {
	component Component
	image     Base64Image
}

const TABLE_SECTION_STYLES = "position: absolute; overflow: visible;"
const CELL_STYLES = "padding: 2px 4px;"

type TableSection struct {
	component Component
	contents  []interface{}
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
