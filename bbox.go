package geojson

import (
	"strconv"

	"github.com/tidwall/gjson"
)

// BBox is a GeoJSON bbox member
type BBox interface {
	Defined() bool
	Rect() Rect
	AppendJSON(dst []byte) []byte
}

type xyBBox struct {
	rect Rect
}

func (bbox xyBBox) Defined() bool {
	return true
}

func (bbox xyBBox) Rect() Rect {
	return bbox.rect
}

func (bbox xyBBox) AppendJSON(dst []byte) []byte {
	dst = append(dst, '[')
	dst = strconv.AppendFloat(dst, bbox.rect.Min.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Min.Y, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Max.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Max.Y, 'f', -1, 64)
	dst = append(dst, ']')
	return dst
}

type xyzBBox struct {
	rect       Rect
	minZ, maxZ float64
}

func (bbox xyzBBox) Defined() bool {
	return true
}

func (bbox xyzBBox) Rect() Rect {
	return bbox.rect
}

func (bbox xyzBBox) AppendJSON(dst []byte) []byte {
	dst = append(dst, '[')
	dst = strconv.AppendFloat(dst, bbox.rect.Min.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Min.Y, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.minZ, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Max.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Max.Y, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.maxZ, 'f', -1, 64)
	dst = append(dst, ']')
	return dst
}

type xyzmBBox struct {
	rect       Rect
	minZ, maxZ float64
	minM, maxM float64
}

func (bbox xyzmBBox) Defined() bool {
	return true
}

func (bbox xyzmBBox) Rect() Rect {
	return bbox.rect
}

func (bbox xyzmBBox) AppendJSON(dst []byte) []byte {
	dst = append(dst, '[')
	dst = strconv.AppendFloat(dst, bbox.rect.Min.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Min.Y, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.minZ, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.minM, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Max.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.rect.Max.Y, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.maxZ, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, bbox.maxM, 'f', -1, 64)
	dst = append(dst, ']')
	return dst
}

// must conform to BBox
var _ = []BBox{
	xyBBox{},
	xyzBBox{},
	xyzmBBox{},
}

func bboxWeight(bbox BBox) int {
	return bboxPositionCount(bbox) * 16
}

func bboxPositionCount(bbox BBox) int {
	if bbox == nil {
		return 0
	}
	return 2
}

func loadBBox(data string) (BBox, error) {
	rbbox := gjson.Get(data, "bbox")
	if !rbbox.Exists() {
		return nil, nil
	}
	if !rbbox.IsArray() {
		return nil, errBBoxInvalid
	}
	var err error
	var count int
	var nums [8]float64
	rbbox.ForEach(func(key, value gjson.Result) bool {
		if count == 8 {
			return false
		}
		if value.Type != gjson.Number {
			err = errBBoxInvalid
			return false
		}
		nums[count] = value.Float()
		count++
		return true
	})
	if err != nil {
		return nil, err
	}
	if count < 4 || count%2 == 1 {
		return nil, errBBoxInvalid
	}
	var rect Rect
	rect.Min.X = nums[0]
	rect.Min.Y = nums[1]
	rect.Max.X = nums[count/2]
	rect.Max.Y = nums[count/2+1]
	if count == 4 {
		bbox := xyBBox{rect: rect}
		return bbox, nil
	}
	if count == 6 {
		bbox := xyzBBox{rect: rect}
		bbox.minZ = nums[2]
		bbox.maxZ = nums[count/2+2]
		return bbox, nil
	}
	bbox := xyzmBBox{rect: rect}
	bbox.minZ = nums[2]
	bbox.minM = nums[3]
	bbox.maxZ = nums[count/2+2]
	bbox.maxM = nums[count/2+3]
	return bbox, nil
}
