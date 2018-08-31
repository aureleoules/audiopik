package main

import (
	"image"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"image/color"
	"image/png"

	"github.com/cryptix/wav"
)

//GetSamples : return samples of WAV file
func GetSamples(file io.ReadSeeker, size int64) ([][]byte, wav.File) {
	var samples [][]byte

	reader, err := wav.NewReader(file, size)
	if err != nil {
		log.Fatal(err)
	}

	for {
		s, err := reader.ReadRawSample()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		samples = append(samples, s)
	}
	return samples, reader.GetFile()
}

//MakeSamples : create new array of sample with a given skip integer
func MakeSamples(samples [][]byte, skip int, channels int) [][]byte {
	var newSamples [][]byte

	i := 0
	for i < len(samples) {
		sample := samples[i]
		if sample[0] < 5 {
			sample[0] = 0
		}
		newSamples = append(newSamples, samples[i])
		i = i + skip*channels
	}
	return newSamples
}

//WriteSamplesToFile : create file and write samples
func WriteSamplesToFile(fileName string, samples [][]byte, sampleRate uint32, significantBits uint16) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	meta := wav.File{
		Channels:        1,
		SampleRate:      sampleRate,
		SignificantBits: significantBits,
	}
	writer, err := meta.NewWriter(file)

	for _, sample := range samples {
		err = writer.WriteSample(sample)
		if err != nil {
			log.Println(err)
		}
	}

	err = writer.Close()
	if err != nil {
		log.Println(err)
	}
}

//CreateImage create image
func CreateImage(samples [][]byte) {
	length := float64(len(samples))
	lengthSqrt := math.Sqrt(length)

	width := int(math.Ceil(lengthSqrt))
	height := int(math.Ceil(lengthSqrt))

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	position := 0

	min, max := MinMax(samples)
	log.Println("min", min)
	log.Println("max", max)

	rand.Seed(time.Now().UnixNano())
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			byteValueLeft := samples[position][0]
			byteValueRight := samples[position][1]

			valueLeft := uint8(byteValueLeft)
			valueRight := uint8(byteValueRight)
			r := uint8(rand.Intn(255))
			log.Println(r)
			img.SetRGBA(x, y, color.RGBA{r, valueLeft, valueRight, 255})

			if position+1 == len(samples) {
				break
			}
			position++
		}
	}

	f, _ := os.Create("output/image.png")
	png.Encode(f, img)
}

//MinMax return min value and max value of array of int
func MinMax(array [][]byte) (int, int) {

	var max = int(array[0][0])
	var min = int(array[0][0])
	for _, byteValue := range array {
		value := int(byteValue[0])
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

//ReadImage read pixels of image
func ReadImage(fileName string) []Pixel {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal("Can't read image", err)
	}

	defer file.Close()

	pixels, err := GetPixels(file)
	if err != nil {
		log.Fatal(err)
	}
	return pixels
}

//GetPixels get pixels of file
func GetPixels(file io.Reader) ([]Pixel, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels []Pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels = append(pixels, rgbaToPixel(img.At(x, y).RGBA()))
		}
	}
	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{R: int(r / 257), G: int(g / 257), B: int(b / 257), A: int(a / 257)}
}

//Pixel structure
type Pixel struct {
	R int
	G int
	B int
	A int
}
