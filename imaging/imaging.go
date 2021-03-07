package imaging

import (
	"bufio"
	"fmt"
	"github.com/barasher/go-exiftool"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/samuel/go-pcx/pcx"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	NoResolution = -1
)

type ImageSize struct {
	Height, Width int
}

func (is ImageSize) String() string {
	return fmt.Sprintf("%dx%d", is.Width, is.Height)
}

type ImageMetadata struct {
	Name            string
	XResolution     int
	YResolution     int
	ColorDepth      int
	CompressionType string
	Size 			ImageSize
}

func getImageMetadataFromFields(fields map[string]interface{}) (metadata *ImageMetadata) {
	var (
		compressionType string
		comps float64
		colorDepth float64
		ok bool
		xRes, yRes float64
	)

	if comps, ok = fields["ColorComponents"].(float64); !ok {
		comps = 1
	}
	if compressionType, ok = fields["Compression"].(string); !ok {
		compressionType = "-"
	}
	if colorDepth, ok = fields["BitsPerPixel"].(float64); !ok {
		if colorDepth, ok = fields["BitDepth"].(float64); !ok {
			colorDepth = 0
			if _, ok = fields["BitsPerSample"]; ok {
				bitsPerSample := fmt.Sprint(fields["BitsPerSample"])
				scanner := bufio.NewScanner(strings.NewReader(bitsPerSample))
				scanner.Split(bufio.ScanWords)
				for scanner.Scan() {
					depth, _ := strconv.ParseFloat(scanner.Text(), 64)
					colorDepth += depth
				}
				colorDepth *= comps
			}
		}
	}
	if xRes, ok = fields["XResolution"].(float64); !ok {
		xRes = NoResolution
	}
	if yRes, ok = fields["YResolution"].(float64); !ok {
		yRes = NoResolution
	}

	return &ImageMetadata {
		Name:            fields["FileName"].(string),
		XResolution:     int(xRes),
		YResolution:     int(yRes),
		ColorDepth: 	 int(colorDepth),
		CompressionType: compressionType,
		Size: ImageSize {
			Height: int(fields["ImageHeight"].(float64)),
			Width:  int(fields["ImageWidth"].(float64)),
		},
	}
}

func GetImagesFromDir(dir string) (images []string, err error) {
	var (
		allFiles []os.FileInfo
		file *os.File
	)

	allFiles, err = ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(allFiles))
	for idx := range allFiles {
		name := filepath.Join(dir, allFiles[idx].Name())
		file, err = os.Open(name)
		if err != nil {
			goto cleanup
		}
		if _, _, err = image.DecodeConfig(file); err != nil {
			goto cleanup
		}
		fileNames = append(fileNames, name)
	cleanup:
		if file != nil {
			file.Close()
		}
	}

	return fileNames, nil
}

func FetchMetadata(files ...string) (data []*ImageMetadata, err error) {
	tools, e := exiftool.NewExiftool(exiftool.Charset("filename=utf8"))
	if e != nil {
		return nil, e
	}
	defer tools.Close()

	innerData := tools.ExtractMetadata(files...)
	result := make([]*ImageMetadata, 0, len(innerData))

	for idx := range innerData {
		if innerData[idx].Err != nil {
			return nil, innerData[idx].Err
		}
		result = append(result, getImageMetadataFromFields(innerData[idx].Fields))
	}

	return result, nil
}