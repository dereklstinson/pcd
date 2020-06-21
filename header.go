package pcd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//PCDVERSION is the version that this package supports
const PCDVERSION = float64(0.7)

//Header is the header for PCD version .7
//Fields,Size,Dtype, and Count should be the same length
type Header struct {
	Comment   []string  `json:"comment,omitempty"`
	Version   float64   `json:"version,omitempty"`
	Fields    []string  `json:"fields,omitempty"`
	Size      []int     `json:"size,omitempty"`
	Dtype     []string  `json:"dtype,omitempty"`
	Count     []int     `json:"count,omitempty"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	Viewpoint []float64 `json:"viewpoint,omitempty"`
	Points    int       `json:"points,omitempty"`
	Data      string    `json:"data,omitempty"`
	rd        *bufio.Reader
}

func (h Header) String() string {
	return fmt.Sprintf("Header{ \n"+
		"Comment: %v\n"+
		"Version: %v\n"+
		"Fields: %v\n"+
		"Size: %v\n"+
		"Dtype: %v\n"+
		"Count: %v\n"+
		"Width: %v\n"+
		"Height: %v\n"+
		"Viewpoint: %v\n"+
		"Points: %v\n"+
		"Data: %v\n"+
		"}", h.Comment,
		h.Version,
		h.Fields,
		h.Size,
		h.Dtype,
		h.Count,
		h.Width,
		h.Height,
		h.Viewpoint,
		h.Points,
		h.Data)

}
func getHeader(r io.Reader) (h Header, err error) {
	rd := bufio.NewReader(r)
	h.rd = rd
	//check if it has a comment string
	for {
		var line string
		line, err = rd.ReadString('\n')
		if err != nil {
			return
		}
		var done bool
		done, err = h.fillheader(line)
		if err != nil {
			return
		}
		if done {

			return
		}
	}

}

func (h *Header) fillheader(line string) (done bool, err error) {

	if strings.Contains(line, "#") {
		x := strings.TrimPrefix(line, "# ")
		if x == line {
			return false, errors.New("ParseError couldn't parse Comment")
		}

		h.Comment = append(h.Comment, x)
		return

	}

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, errors.New("Space filled fields")
	}
	flds := fields[1:]
	switch fields[0] {
	case "VERSION":
		h.Version, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return false, err
		}
		if h.Version != PCDVERSION {
			return false, fmt.Errorf("Version Error: Supported Version is %f, File Header says its version %f", PCDVERSION, h.Version)
		}
		return false, err
	case "FIELDS":

		for _, field := range flds {
			h.Fields = append(h.Fields, field)
		}

	case "SIZE":
		for _, field := range flds {
			size, err := strconv.Atoi(field)
			if err != nil {
				return false, err
			}
			h.Size = append(h.Size, size)
		}
	case "TYPE":
		for _, field := range flds {
			h.Dtype = append(h.Dtype, field)
		}
	case "COUNT":
		for _, field := range flds {
			countint, err := strconv.Atoi(field)
			if err != nil {
				return false, err
			}
			h.Count = append(h.Count, countint)
		}

	case "WIDTH":
		widthint, err := strconv.Atoi(flds[0])
		if err != nil {
			return false, err
		}
		h.Width = widthint
		return false, nil
	case "HEIGHT":
		heightint, err := strconv.Atoi(flds[0])
		if err != nil {
			return false, err
		}
		h.Height = heightint
		return false, nil
	case "VIEWPOINT":
		for _, field := range flds {
			viewfloat, err := strconv.ParseFloat(field, 64)
			if err != nil {
				return false, err
			}
			h.Viewpoint = append(h.Viewpoint, viewfloat)
		}
	case "POINTS":
		pointint, err := strconv.Atoi(flds[0])
		if err != nil {
			return false, err
		}
		h.Points = pointint
	case "DATA":
		h.Data = strings.ToLower(flds[0])
		if h.Data != "binary" && h.Data != "ascii" {
			return false, errors.New("DATA != ascii or binary")
		}
		return true, nil

	default:
		return false, errors.New("Unsupported Header Line")
	}
	return false, nil
}
