package heightmap

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"image"
	"image/color"

	"image/png"
	"log/slog"
	"nw-buddy/tools/nwfs"
	"nw-buddy/tools/utils"
	"nw-buddy/tools/utils/env"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/x448/float16"
	"golang.org/x/image/draw"
)

//"github.com/rngoodner/gtiff"

var (
	ErrMapSettingsNotFound = errors.New("mapsettings.json not found")
	ErrRegionSizeNotFound  = errors.New("regionSize not found in mapsettings.json")
	ErrMetdataNotLoaded    = errors.New("metadata not loaded")
)

type Meta struct {
	X     int
	Y     int
	Level string
}

type Region struct {
	Meta
	Size int
	File string
	Data []uint16
}

type MapSettings struct {
	CellResolution int `json:"cellResolution"`
	RegionSize     int `json:"regionSize"`
	RegionType     int `json:"regionType"`
}

func LoadRegion(file nwfs.File) (region Region, err error) {
	settingsFile, ok := file.Archive().Lookup(path.Join(path.Dir(file.Path()), "mapsettings.json"))
	if !ok {
		err = ErrMapSettingsNotFound
		return
	}
	settingsData, err := settingsFile.Read()
	if err != nil {
		return
	}
	var settings MapSettings
	err = json.Unmarshal(settingsData, &settings)
	if err != nil {
		return
	}
	region.Size = settings.RegionSize
	region.Meta, ok = ReadPathMetadata(file.Path())
	if !ok {
		err = ErrMetdataNotLoaded
		return
	}

	region.File = file.Path()
	region.Data, err = Load(file)
	if err != nil {
		return
	}
	expectedSize := region.Size * region.Size
	if len(region.Data) != expectedSize {
		slog.Warn("heightmap data size mismatch", "expected", expectedSize, "actual", len(region.Data), "size", region.Size, "file", file.Path())
	}
	return
}

func Load(file nwfs.File) ([]uint16, error) {
	data, err := file.Read()
	if err != nil {
		return nil, err
	}
	return LoadFromTiff(data)
	// return LoadFieldOld(data)
}

func LoadImage(file nwfs.File) (image.Image, error) {
	data, err := file.Read()
	if err != nil {
		return nil, err
	}
	hm, err := ParseTIFF(data)
	if err != nil {
		return nil, err
	}
	return hm.Image(), nil
	//return convertToImageWithMagick(data)
}

func LoadFromTiff(data []byte) ([]uint16, error) {
	r, err := ParseTIFF(data)
	if err != nil {
		return nil, err
	}
	return r.Samples, nil
}

func LoadFieldOld(data []byte) ([]uint16, error) {
	img, err := convertToImageWithMagick(data)
	if err != nil {
		return nil, err
	}
	return LoadHeightFieldFromImage(img)
}

func convertToImageWithMagick(data []byte) (image.Image, error) {
	f, err := os.CreateTemp(env.TempDir(), "*")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return nil, err
	}
	f.Close()

	tiffNam := f.Name()
	pngName := f.Name() + ".png"
	defer os.Remove(tiffNam)
	defer os.Remove(pngName)

	cmd := utils.Magick.Convert(f.Name(), f.Name()+".png")
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	pngFile, err := os.Open(pngName)
	if err != nil {
		return nil, err
	}
	defer pngFile.Close()
	return png.Decode(pngFile)
}

func LoadHeightFieldFromImage(img image.Image) ([]uint16, error) {
	sizeX := img.Bounds().Size().X
	sizeY := img.Bounds().Size().Y
	out := make([]uint16, sizeX*sizeY)
	index := 0
	for y := range sizeY {
		for x := range sizeX {
			r, _, _, _ := img.At(x, y).RGBA()
			out[index] = uint16(r)
			index++
		}
	}

	return out, nil
}

func LoadTractmapFromFile(file nwfs.File, size int) ([]color.Color, error) {
	data, err := file.Read()
	if err != nil {
		return nil, err
	}
	img, err := convertToImageWithMagick(data)
	if err != nil {
		return nil, err
	}
	return LoadTractmapFromImage(img, size)
}

func LoadTractmapFromImage(img image.Image, size int) ([]color.Color, error) {
	sizeX := img.Bounds().Size().X
	sizeY := img.Bounds().Size().Y
	if sizeX != size || sizeY != size {
		dst := image.NewNRGBA(image.Rect(0, 0, size, size))
		draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
		img = dst
	}
	out := make([]color.Color, size*size)
	index := 0
	for y := range size {
		for x := range size {
			out[index] = img.At(x, y)
			index++
		}
	}

	return out, nil
}

var metaRegexp = regexp.MustCompile(`/([^/]+)/regions/r_\+(\d{2})_\+(\d{2}).*`)

func ReadPathMetadata(filePath string) (out Meta, ok bool) {
	match := metaRegexp.FindStringSubmatch(filePath)
	if len(match) != 4 {
		return out, false
	}
	x, _ := strconv.Atoi(match[2])
	y, _ := strconv.Atoi(match[3])
	out.Level = match[1]
	out.X = x
	out.Y = y
	return out, true
}

// ToUint16Bytes returns the heightmap as little-endian uint16 bytes, ready to upload
// as an r16uint texture (2*Size*Size bytes).
func (r *Region) ToUint16Bytes() []byte {
	out := make([]byte, len(r.Data)*2)
	for i, bits := range r.Data {
		binary.LittleEndian.PutUint16(out[i*2:], bits)
	}
	return out
}

// ToFloat16Bytes returns the heightmap as little-endian uint16 bytes, ready to upload
// as an r16float texture (2*Size*Size bytes).
func (r *Region) ToFloat16Bytes() []byte {
	out := make([]byte, len(r.Data)*2)
	for i, v := range r.Data {
		bits := float16.Fromfloat32(float32(v)).Bits()
		binary.LittleEndian.PutUint16(out[i*2:], bits)
	}
	return out
}
