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

	evolution "github.com/damoon/lab-evolution/pkg"
	"github.com/llgcode/draw2d/draw2dimg"
)

const triangleCount = 20
const populationCount = 500

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

	genSize := 3 + triangleCount*(3+6*16)

	r := rand.New(rand.NewSource(0))

	population := evolution.NewPopulation(r, genSize, populationCount, eva(img), fit)

	champion := population.Fittest()
	log.Printf("champion #0 eva=%f fit=%f", champion.Evaluation, champion.Fitness)
	championImage := genReaderToImage(evolution.GenomeReader{Gens: champion.Genes}, img)
	draw2dimg.SaveToPngFile("0.png", championImage)

	//	for j, champion := range population {
	//		championImage := genReaderToImage(evolution.GenomeReader{Gens: champion.Genes}, img)
	//		draw2dimg.SaveToPngFile(fmt.Sprintf("_0_%d.png", j), championImage)
	//	}

	for i := 0; i < 30; i++ {
		log.Printf("generation %d", i)
		population = population.Evolve(
			r,
			eva(img),
			fit,
			evolution.SimpleParentSelector,
			evolution.SimpleGenExchange,
			evolution.SimpleMutation,
		)

		champion := population.Fittest()
		log.Printf("champion #0 eva=%f fit=%f", champion.Evaluation, champion.Fitness)
		championImage := genReaderToImage(evolution.GenomeReader{Gens: champion.Genes}, img)
		draw2dimg.SaveToPngFile(fmt.Sprintf("%d.png", i+1), championImage)

		//		for j, champion := range population {
		//			championImage := genReaderToImage(evolution.GenomeReader{Gens: champion.Genes}, img)
		//			draw2dimg.SaveToPngFile(fmt.Sprintf("_%d_%d.png", i+1, j), championImage)
		//		}
		var eval float32 = 0.0
		for _, champion := range population {
			eval += champion.Evaluation
		}
		eval /= float32(len(population))
		log.Printf("avg eval=%f", eval)
	}

	return nil
}

func fit(v, min, max float32) float32 {
	return -v + max
}

func eva(target image.Image) evolution.EvaluationFunc {
	return func(genReader evolution.GenomeReader) float32 {
		dest := genReaderToImage(genReader, target)

		d, err := difference(target, dest)
		if err != nil {
			panic(err)
		}

		return d
	}
}

func genReaderToImage(genReader evolution.GenomeReader, target image.Image) image.Image {
	dest := image.NewRGBA(target.Bounds())

	background := color.RGBA{
		genReader.Uint8(),
		genReader.Uint8(),
		genReader.Uint8(),
		0xff,
	}
	draw.Draw(dest, dest.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)

	gc := draw2dimg.NewGraphicContext(dest)

	for i := 0; i < triangleCount; i++ {
		c := color.RGBA{
			genReader.Uint8(),
			genReader.Uint8(),
			genReader.Uint8(),
			0x88,
		}

		// Set some properties
		gc.SetFillColor(c)
		gc.SetStrokeColor(c)
		gc.SetLineWidth(1)

		xOffset := -0.1 * float64(dest.Bounds().Dx())
		yOffset := -0.1 * float64(dest.Bounds().Dy())

		p1x := xOffset + 1.2*float64(dest.Bounds().Dx()*int(genReader.Uint16()))/math.Pow(2, 16)
		p1y := yOffset + 1.2*float64(dest.Bounds().Dy()*int(genReader.Uint16()))/math.Pow(2, 16)
		p2x := xOffset + 1.2*float64(dest.Bounds().Dx()*int(genReader.Uint16()))/math.Pow(2, 16)
		p2y := yOffset + 1.2*float64(dest.Bounds().Dy()*int(genReader.Uint16()))/math.Pow(2, 16)
		p3x := xOffset + 1.2*float64(dest.Bounds().Dx()*int(genReader.Uint16()))/math.Pow(2, 16)
		p3y := yOffset + 1.2*float64(dest.Bounds().Dy()*int(genReader.Uint16()))/math.Pow(2, 16)

		// Draw a closed shape
		gc.BeginPath()      // Initialize a new path
		gc.MoveTo(p1x, p1y) // Move to a position to start the new path
		gc.LineTo(p2x, p2y)
		gc.LineTo(p3x, p3y)
		//		gc.QuadCurveTo(100, 200, 100, 100)
		gc.Close()
		gc.FillStroke()
	}
	return dest
}

func difference(img1, img2 image.Image) (float32, error) {
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

	return float32(diff), nil
}
