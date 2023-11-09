package main

type Page struct {
	Height          int
	Width           int
	Id              string
	BackgroundColor string
	TableSection    []TableSection `kdl:",omitempty,multiple,child"`
	// Viz             Viz   `kdl:",omitempty,child,multiple"`
	Contents []interface{} `kdl:",omitempty"`
	// children []interface{}
}

type PageHeader struct {
	height          int
	backgroundColor string
}
