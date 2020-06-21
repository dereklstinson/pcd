package pcd

import (
	"fmt"
	"os"
	"testing"
)

func Test_getHeader(t *testing.T) {

	file, err := os.Open("example.PCD")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	h, err := getHeader(file)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(h)
}
func Test_GetPoints(t *testing.T) {
	file, err := os.Open("example.PCD")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	h, err := getHeader(file)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(h)

	points, err := h.GetPoints(h.Points)
	if err != nil {
		t.Error(err)
	}
	if len(points) != h.Points {

		t.Error("Len of points is off", len(points), h.Points)
	}
	for _, p := range points {
		fields := p.GetFields()
		for _, f := range fields {
			fmt.Println(f.GetValuesf64())
		}
	}

}
