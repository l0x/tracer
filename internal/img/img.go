package img

type Img struct {
	Width  int
	Height int
	Pixels []uint8
}

func NewImg(width, height int) *Img {
	return &Img{
		Width:  width,
		Height: height,
		Pixels: make([]uint8, width*height*3),
	}
}

func (img *Img) SetPixel(x, y int, r, g, b uint8) {
	index := (y*img.Width + x) * 3
	img.Pixels[index] = r
	img.Pixels[index+1] = g
	img.Pixels[index+2] = b
}

func (img *Img) GetPixel(x, y int) (uint8, uint8, uint8) {
	index := (y*img.Width + x)
	return img.Pixels[index], img.Pixels[index+1], img.Pixels[index+2]
}
