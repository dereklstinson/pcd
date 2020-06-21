package pcd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	lzf "github.com/dereklstinson/golzf"
)

//Line is the per line bytes.
type Line []byte

//GetLines gets up to nlines from r up intil it
func getlines(rd *bufio.Reader, nlines int) (l []Line, err error) {

	for i := 0; i < nlines; i++ {

		nl, err := rd.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				l = append(l, nl)
			} else {
				panic(err)
			}

			return l, err

		}
		l = append(l, nl)
	}
	return l, nil
}

//GetPoints will get points up until an error is found or end of file.  If no data is found p will be nil.  Otherwise p will always returns something.
func (h *Header) GetPoints(npoints int) (p []Point, err error) {
	l, err := getlines(h.rd, npoints)
	if err != nil {
		if len(l) == 0 {
			return nil, err
		}
		if err == io.EOF {
			err = nil
		}
	}
	if len(l) != npoints {
		fmt.Println("Len(l) is != npoints")
	}
	p = make([]Point, len(l))
	for i := range p {

		p[i].fields = h.extractFields(l[i])
	}
	return p, err
}

//File contains both the Header and the Data.  This could be used to save to Json format if needed.
type File struct {
	Header Header `json:"header,omitempty"`
	Data   []Line `json:"data,omitempty"`
}

//A Point is a (semi) decoded Line seperated out into fields
type Point struct {
	fields []Field
}

//GetFields returns the fields of the point
func (p *Point) GetFields() []Field {
	return p.fields
}

//GetValuesf64 is a helper/standin function to quickly get values
func (f *Field) GetValuesf64() (dim string, vals []float64) {

	switch f.h.Data {
	case "ascii":
		var dtype string
		var size int
		var count int
		dim, dtype, size, count = f.GetFieldInfo()
		vals = make([]float64, count)
		var err error
		for i := 0; i < count; i++ {
			fmt.Println(len(f.data))
			valstring := string(f.data[i])
			switch dtype {
			case "F":
				if size == 4 {
					vals[i], err = strconv.ParseFloat(valstring, 32)
				} else if size == 8 {
					vals[i], err = strconv.ParseFloat(valstring, 64)
				} else {
					return "", nil
				}

			default:
				var vint int
				vint, err = strconv.Atoi(valstring)
				vals[i] = float64(vint)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

	case "binary":
		panic("Bindary Not supported right now")
	case "binary_compressed":
		panic("compressedbinary not working")
	default:
		return "", nil
	}

	return
}

//GetFieldInfo gets the info for the field.
func (f *Field) GetFieldInfo() (dim, dtype string, size, count int) {
	i := f.fieldindex
	dim = f.h.Fields[i]
	dtype = f.h.Dtype[i]
	size = f.h.Size[i]
	count = f.h.Count[i]
	return
}

//Field is data contained in a line.
type Field struct {
	fieldindex int
	data       [][]byte // some data has many fields
	h          *Header
}

func (h *Header) extractFields(l Line) []Field {
	nfields := len(h.Fields)
	f := make([]Field, nfields)
	switch h.Data {
	case "ascii":
		flds := bytes.Fields(l)
		if nfields != len(flds) {
			return nil
		}
		var offsetindex int
		for i := range flds {
			f[i].fieldindex = i
			f[i].h = h
			f[i].data = make([][]byte, h.Count[i])
			for j := range f[i].data {
				f[i].data[j] = make([]byte, len(flds[offsetindex]))
				copy(f[i].data[j], flds[offsetindex])
				offsetindex++
			}

		}
	default:
		panic("Unsupported format")
	}

	return f
}

func (h *Header) compressdata(compressdata []byte) (compressed []byte, err error) {
	holder := make([]byte, len(compressdata))
	var n int
	n, err = lzf.Compress(compressdata, holder)
	if err != nil {
		return nil, err
	}
	compressed = make([]byte, n)
	copy(compressed, holder[:n])
	return compressed, err

}

//decompressdata decompresses lzf
func (h *Header) decompressdata(compressed []byte) (decompressed []byte, err error) {

	decompressed = make([]byte, int(compressed[0])<<24|int(compressed[1])<<16|int(compressed[2])<<8|int(compressed[3])<<3)

	var n int
	n, err = lzf.Decompress(compressed, decompressed)
	return decompressed[:n], err

}
