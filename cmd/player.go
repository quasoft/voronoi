// player is a web server, hosting a visualization of the Fortune's voronoi generation algorithm,
// implemented by the Voronoi structure.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"

	"github.com/quasoft/voronoi"
)

// D3Node is the data required to draw a node with D3.js library.
type D3Node struct {
	Name     string   `json:"name"`
	Children []D3Node `json:"children"`
}

// btreeToGraphNode converts a parabola binary tree to a tree of D3Nodes.
func btreeToGraphNode(node *voronoi.Node) *D3Node {
	if node == nil {
		return nil
	}

	left := btreeToGraphNode(node.Left)
	right := btreeToGraphNode(node.Right)

	site := node.Site
	var label string

	if site.X == 0 && site.Y == 0 {
		label = "Internal"
	} else {
		label = fmt.Sprintf("Site %v", site)
	}

	var children []D3Node
	if left != nil && right != nil {
		children = []D3Node{*left, *right}
	} else if left != nil {
		children = []D3Node{*left, D3Node{"-", nil}}
	} else if right != nil {
		children = []D3Node{D3Node{"-", nil}, *right}
	}
	return &D3Node{label, children}
}

// btreeToJSON converts a parabola binary tree to a JSON tree of D3 nodes.
func btreeToJSON(node *voronoi.Node) []byte {
	graphNode := btreeToGraphNode(node)
	jsonTree, err := json.Marshal(graphNode)
	if err != nil {
		return []byte{}
	}
	return jsonTree
}

func main() {
	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	// Start web server
	width, height := 600, 480
	rect := image.Rect(0, 0, width, height)

	sites := []image.Point{
		{X: 110, Y: 20},
		{X: 140, Y: 40},
		{X: 155, Y: 80},
		{X: 350, Y: 120},
		{X: 200, Y: 240},
	}

	v := voronoi.NewFromPoints(sites, rect)
	var img *image.RGBA

	// Index page of visualization
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		img = voronoi.Plot(v)

		t, err := template.ParseFiles("index.html")
		if err != nil {
			panic(err)
		}

		data := struct {
			SweepLine  int
			EventsLeft int
			Log        string
		}{
			v.SweepLine,
			v.EventQueue.Len(),
			logBuf.String(),
		}

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	// Reset the state of the voronoi generator
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		logBuf.Reset()
		v.Reset()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// Handle next event from the queue and update the visualization
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		v.HandleNextEvent()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// Plot the voronoi diagram state into an image
	http.HandleFunc("/diagram.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
	})

	// Transform the internal binary tree to json, suitable for visualization in D3
	http.HandleFunc("/tree.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		jsonTree := btreeToJSON(v.ParabolaTree)
		w.Write([]byte{byte('[')})
		w.Write(jsonTree)
		w.Write([]byte{byte(']')})
	})

	fmt.Printf("Listening at 127.0.0.1:3000\r\n")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}
