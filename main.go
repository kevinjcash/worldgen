package main

import (
	"fmt"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"time"
)

const (
	imageWidth  int = 1900
	imageHeight int = 1080
	nSites      int = 256
)

func main() {
	sx, sy := randomSites()
	diagram := generateVoronoi(sx, sy)
	relaxed_points := utils.LloydRelaxation(diagram.Cells)
	for i, point := range relaxed_points {
		sx[i] = point.X
		sy[i] = point.Y
	}
	diagram = generateVoronoi(sx, sy)
	writePngFile(drawVoronoi(diagram))
}

func generateVoronoi(sx []float64, sy []float64) (diagram *voronoi.Diagram) {
	// Sites of voronoi diagram
	var sites []voronoi.Vertex
	for i := 0; i < nSites; i++ {
		sites = append(sites, voronoi.Vertex{sx[i], sy[i]})
	}

	// Create bounding box of [0, 20] in X axis
	// and [0, 10] in Y axis
	bbox := voronoi.NewBBox(0, float64(imageWidth), 0, float64(imageHeight))

	// Compute diagram and close cells (add half edges from bounding box)
	// diagram := NewVoronoi().Compute(sites, bbox, true)
	diagram = voronoi.ComputeDiagram(sites, bbox, true)
	return diagram
}

func drawVoronoi(diagram *voronoi.Diagram) image.Image {
	// generate a random color for each site
	sc := make([]color.NRGBA, nSites)
	for i := 0; i < len(diagram.Cells); i++ {
		sc[i] = color.NRGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)),
			uint8(rand.Intn(256)), 255}
	}

	// generate diagram by coloring each pixel with color of nearest site
	img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			dMin := dot(imageWidth, imageHeight)
			var sMin int
			for s, cell := range diagram.Cells {
				if d := dot(int(cell.Site.X)-x, int(cell.Site.Y)-y); d < dMin {
					sMin = s
					dMin = d
				}
			}
			img.SetNRGBA(x, y, sc[sMin])
		}
	}

	// mark each site with a black box
	black := image.NewUniform(color.Black)
	for _, cell := range diagram.Cells {
		draw.Draw(img, image.Rect(int(cell.Site.X)-2, int(cell.Site.Y)-2, int(cell.Site.X)+2, int(cell.Site.Y)+2),
			black, image.ZP, draw.Src)
	}
	return img
}

func dot(x, y int) int {
	return x*x + y*y
}

func randomSites() (sx, sy []float64) {
	rand.Seed(time.Now().Unix())
	sx = make([]float64, nSites)
	sy = make([]float64, nSites)
	for i := range sx {
		sx[i] = float64(rand.Intn(imageWidth))
		sy[i] = float64(rand.Intn(imageHeight))
	}
	return
}

func writePngFile(img image.Image) {
	f, err := os.Create("voronoi.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = png.Encode(f, img); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}
