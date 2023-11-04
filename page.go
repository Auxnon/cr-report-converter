package main

type Page struct {
	height          int
	width           int
	id              string
	backgroundColor string
	contents        []interface{}
}

type PageHeader struct {
	height          int
	backgroundColor string
}
