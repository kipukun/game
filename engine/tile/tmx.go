package tile

import (
	"encoding/csv"
	"encoding/xml"
	"image"
	_ "image/png"
	"io"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kipukun/game/engine/errors"
	"github.com/kipukun/game/engine/object"
)

type tsx struct {
	XMLName      xml.Name `xml:"tileset"`
	Text         string   `xml:",chardata"`
	Version      string   `xml:"version,attr"`
	Tiledversion string   `xml:"tiledversion,attr"`
	Name         string   `xml:"name,attr"`
	Tilewidth    string   `xml:"tilewidth,attr"`
	Tileheight   string   `xml:"tileheight,attr"`
	Tilecount    string   `xml:"tilecount,attr"`
	Columns      string   `xml:"columns,attr"`
	Image        struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source,attr"`
		Width  string `xml:"width,attr"`
		Height string `xml:"height,attr"`
	} `xml:"image"`
}

type tmx struct {
	XMLName      xml.Name `xml:"map"`
	Text         string   `xml:",chardata"`
	Version      string   `xml:"version,attr"`
	Tiledversion string   `xml:"tiledversion,attr"`
	Orientation  string   `xml:"orientation,attr"`
	Renderorder  string   `xml:"renderorder,attr"`
	Width        string   `xml:"width,attr"`
	Height       string   `xml:"height,attr"`
	Tilewidth    string   `xml:"tilewidth,attr"`
	Tileheight   string   `xml:"tileheight,attr"`
	Infinite     string   `xml:"infinite,attr"`
	Nextlayerid  string   `xml:"nextlayerid,attr"`
	Nextobjectid string   `xml:"nextobjectid,attr"`
	Tileset      struct {
		Text     string `xml:",chardata"`
		Firstgid string `xml:"firstgid,attr"`
		Source   string `xml:"source,attr"`
	} `xml:"tileset"`
	Layer []struct {
		Text   string `xml:",chardata"`
		ID     string `xml:"id,attr"`
		Name   string `xml:"name,attr"`
		Width  string `xml:"width,attr"`
		Height string `xml:"height,attr"`
		Data   struct {
			Text     string `xml:",chardata"`
			Encoding string `xml:"encoding,attr"`
		} `xml:"data"`
	} `xml:"layer"`
}

func NewTileSheetFromTSX(sheet, file io.Reader) (*TileSheet, error) {
	var op errors.Op = "NewTileSheetFromTSX"

	xdec := xml.NewDecoder(file)
	var tsx *tsx
	err := xdec.Decode(&tsx)
	if err != nil {
		return nil, errors.Error(op, "error decoding XML", err)
	}
	decimg, _, err := image.Decode(sheet)
	if err != nil {
		return nil, errors.Error(op, "error decoding image", err)
	}
	img := ebiten.NewImageFromImage(decimg)
	dx, err := strconv.Atoi(tsx.Tilewidth)
	if err != nil {
		return nil, errors.Error(op, "error converting string to int", err)
	}
	dy, err := strconv.Atoi(tsx.Tileheight)
	if err != nil {
		return nil, errors.Error(op, "error converting string to int", err)
	}
	return NewTileSheet(img, dx, dy), nil
}

func NewTileMapFromTMX(s *TileSheet, file io.Reader) (*ebiten.Image, map[string]object.Collection, error) {
	var op errors.Op = "NewTileMapFromTMX"

	om := make(map[string]object.Collection)
	xdec := xml.NewDecoder(file)
	var tmx *tmx
	err := xdec.Decode(&tmx)
	if err != nil {
		return nil, nil, errors.Error(op, "error decoding XML", err)
	}

	dx, err := strconv.Atoi(tmx.Tilewidth)
	if err != nil {
		return nil, nil, errors.Error(op, "error converting string to int", err)
	}
	dy, err := strconv.Atoi(tmx.Tileheight)
	if err != nil {
		return nil, nil, errors.Error(op, "error converting string to int", err)
	}
	width, err := strconv.Atoi(tmx.Width)
	if err != nil {
		return nil, nil, errors.Error(op, "error converting string to int", err)
	}
	height, err := strconv.Atoi(tmx.Height)
	if err != nil {
		return nil, nil, errors.Error(op, "error converting string to int", err)
	}

	img := ebiten.NewImage(dx*width, dy*height)
	for _, layer := range tmx.Layer {
		var objectLayer bool
		var prefix string
		if strings.HasPrefix(layer.Name, "objects:") {
			prefix = strings.TrimPrefix(layer.Name, "objects:")
			objectLayer = true
		}
		cdec := csv.NewReader(strings.NewReader(layer.Data.Text))
		cdec.FieldsPerRecord = -1
		records, err := cdec.ReadAll()
		if err != nil {
			return nil, nil, errors.Error(op, "error reading records", err)
		}
		for i, row := range records {
			for j, pt := range row {
				if j+1 == len(row) && i+1 != len(records) {
					continue
				}
				coord, err := strconv.Atoi(pt)
				if err != nil {
					return nil, nil, errors.Error(op, "error converting string to int", err)
				}
				if coord == 0 {
					continue
				}
				t, o := s.Tile(coord, 0)
				posx, posy := float64(j%width*dx), float64(i*dy)
				if objectLayer {
					obj, _ := object.FromEbitenImage(t)
					obj.SetPosition(posx, posy)
					if om[prefix] == nil {
						om[prefix] = object.NewCollection()
					}
					om[prefix] = append(om[prefix], obj)
				}
				o.GeoM.Translate(posx, posy)
				img.DrawImage(t, o)
			}
		}
	}
	return img, om, nil
}
