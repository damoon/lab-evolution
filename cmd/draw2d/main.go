package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		return err
	}

	rand.Seed(0)

	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(img.Bounds())

	background := color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		0xff,
	}
	draw.Draw(dest, dest.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)

	gc := draw2dimg.NewGraphicContext(dest)

	for i := 0; i < 100; i++ {
		c := color.RGBA{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			0xff,
		}

		// Set some properties
		gc.SetFillColor(c)
		gc.SetStrokeColor(c)
		gc.SetLineWidth(1)

		xOffset := -0.1 * float64(dest.Bounds().Dx())
		yOffset := -0.1 * float64(dest.Bounds().Dy())

		p1x := xOffset + 1.2*float64(rand.Intn(dest.Bounds().Dx()))
		p1y := yOffset + 1.2*float64(rand.Intn(dest.Bounds().Dy()))
		p2x := xOffset + 1.2*float64(rand.Intn(dest.Bounds().Dx()))
		p2y := yOffset + 1.2*float64(rand.Intn(dest.Bounds().Dy()))
		p3x := xOffset + 1.2*float64(rand.Intn(dest.Bounds().Dx()))
		p3y := yOffset + 1.2*float64(rand.Intn(dest.Bounds().Dy()))

		// Draw a closed shape
		gc.BeginPath()      // Initialize a new path
		gc.MoveTo(p1x, p1y) // Move to a position to start the new path
		gc.LineTo(p2x, p2y)
		gc.LineTo(p3x, p3y)
		//		gc.QuadCurveTo(100, 200, 100, 100)
		gc.Close()
		gc.FillStroke()

		// Save to file
		// draw2dimg.SaveToPngFile(fmt.Sprintf("%d.png", i), dest)

		d, err := difference(img, dest)
		if err != nil {
			return err
		}

		log.Printf("image=%d difference=%.2f", i, d)
	}

	draw2dimg.SaveToPngFile("out.png", dest)

	return nil
}

func difference(img1, img2 image.Image) (float64, error) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx() != bounds2.Dx() {
		return 0, fmt.Errorf("images need to have the same dimentions")
	}

	diff := 0.0

	for x := 0; x < bounds1.Dx(); x++ {
		for y := 0; y < bounds1.Dx(); y++ {
			px1 := img1.At(x, y)
			px2 := img2.At(x, y)

			r1, g1, b1, _ := px1.RGBA()
			r2, g2, b2, _ := px2.RGBA()

			diff += math.Sqrt(float64(r1 * r2))
			diff += math.Sqrt(float64(g1 * g2))
			diff += math.Sqrt(float64(b1 * b2))
		}
	}

	return diff, nil
}
