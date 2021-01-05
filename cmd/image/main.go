package main

import (
	"image"
	_ "image/jpeg"
	"log"
	"os"

	evolution "github.com/damoon/lab-evolution/pkg"
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

	return nil
}

func genomeToImage(genome evolution.Genome) image.Image {

}

func similar(original, generated image.Image) float32 {

}
