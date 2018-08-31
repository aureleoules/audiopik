package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	filePath := flag.String("infile", "audio.wav", "wav file to read")
	flag.Parse()

	infos, _ := os.Stat(*filePath)
	file, _ := os.Open(*filePath)

	samples, fileMeta := GetSamples(file, infos.Size())

	SKIP := 6     // The more means lower quality
	CHANNELS := 1 // Mono = 1; Stereo = 2

	newSamples := MakeSamples(samples, SKIP, CHANNELS)

	log.Println("Length", len(newSamples))

	WriteSamplesToFile("output/sampled.wav", newSamples, fileMeta.SampleRate/uint32(SKIP*CHANNELS), fileMeta.SignificantBits)
	CreateImage(newSamples)

	pixels := ReadImage("output/image.png")

	var pixelSamples [][]byte
	for _, pixel := range pixels {

		leftByte := byte(pixel.G)
		rightByte := byte(pixel.B)

		sample := []byte{leftByte, rightByte}
		pixelSamples = append(pixelSamples, sample)
	}

	WriteSamplesToFile("output/from_image.wav", pixelSamples, fileMeta.SampleRate/uint32(SKIP), fileMeta.SignificantBits)
}
