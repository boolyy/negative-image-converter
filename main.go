package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const PixelOffsetforRGB24 int64 = 54

func verifyCorrectFormat(imgBytes []byte) error {
	if !(isIdFieldValid(imgBytes) && isPixelOffsetValid(imgBytes)) {
		return fmt.Errorf("file is not correct format of RGB24 bmp")
	}

	return nil
}

func isIdFieldValid(imgBytes []byte) bool {

	const BmpFirstIdField uint8 = 66
	const BmpSecondIdField uint8 = 77

	if len(imgBytes) < 2 || imgBytes[0] != BmpFirstIdField || imgBytes[1] != BmpSecondIdField {
		return false
	}

	return true
}

func isPixelOffsetValid(imgBytes []byte) bool {

	const LocationOfPixelOffset int = 9
	var pixelArrayOffset int64 = 0

	if len(imgBytes) <= LocationOfPixelOffset+4 {
		return false
	}

	for i := LocationOfPixelOffset; i <= LocationOfPixelOffset+4; i++ {
		pixelArrayOffset += int64(imgBytes[i])
	}

	return pixelArrayOffset == PixelOffsetforRGB24
}

func main() {

	imgFilePath := flag.String("image", "", "path of the image file that you want the negative of")

	flag.Parse()

	imgFile, err := os.Open(*imgFilePath)
	if err != nil {
		fmt.Print(err)
		return
	}

	imgFileInfo, err := imgFile.Stat()
	if err != nil {
		fmt.Print(err)
		return
	}

	numOfBytes := imgFileInfo.Size()

	imgBytes := make([]byte, numOfBytes)

	_, err = imgFile.Read(imgBytes)
	if err != nil {
		fmt.Print(err)
		return
	}

	err = verifyCorrectFormat(imgBytes)
	if err != nil {
		fmt.Print(err)
		return
	}

	const MaxRgbValue uint8 = 255

	//Subtract each pixel from the maximum rgb value to get the negative value
	for byteNum := PixelOffsetforRGB24 - 1; byteNum < numOfBytes; byteNum++ {
		imgBytes[byteNum] = MaxRgbValue - imgBytes[byteNum]
	}

	// Writes new byte list to new file
	permissions := 0644
	err = os.WriteFile("negative-"+filepath.Base(*imgFilePath), imgBytes, fs.FileMode(permissions))
	if err != nil {
		fmt.Println(err)
	}

}
