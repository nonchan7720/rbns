// +build !embed

package main

import "flag"

func init() {
	flag.StringVar(&ui.root, "staticRoot", "./static/dist", "web ui static root")
	flag.BoolVar(&ui.indexes, "staticIndexes", false, "web ui static allow indexes")
}
