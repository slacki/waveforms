package spectogram

import (
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/slacki/waveforms/wavreader"
)

// Config is a set of options for Spectogram
type Config struct {
	// TODO: ?
	Ratio float64

	Width  uint
	Height uint
	// TODO: ?
	Bins uint
	// TODO: wtf, even square is rectangle after all?
	Rectangle bool

	// TODO: no idea, but probably used when actually generating spect
	HideAverage bool
	// TODO: also tbd
	HideRulers bool

	// TODO: poke around with this setting and check out what does it mean
	PreEmphasis float64

	// TODO: the followint 3 options don't make any sense for me yet
	DFT   bool
	LOG10 bool
	MAG   bool

	BG0, BG1, FG0, FG1 color.Color
	RulerColor         color.Color
}

// Spectogram represents a spectogram to be generated
type Spectogram struct {
	Config  *Config
	samples []float64
	File    string
}

// NewSpectogram creates new Spectogram instance and configures it with Config
func NewSpectogram(c *Config, file string) (*Spectogram, error) {
	if c.PreEmphasis == 0 {
		c.PreEmphasis = 0.95
	}

	s := &Spectogram{
		Config: c,
		File:   file,
	}

	err := s.sampleWav()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Spectogram) sampleWav() error {
	file, err := os.Open(s.File)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wr, err := wavreader.New(file)
	if err != nil {
		return err
	}

	length := wr.Len()
	s.samples = make([]float64, length)
	for i := uint64(0); i < length; i++ {
		samp, err := wr.At(0, i)
		if err != nil {
			return err
		}
		s.samples[i] = float64(samp)
	}

	return nil
}

func (s *Spectogram) preEmphasis() {
	for i := len(s.samples) - 1; i > 0; i-- {
		s.samples[i] = s.samples[i] - s.Config.PreEmphasis*s.samples[i-1]
	}
}

func (s *Spectogram) imageBounds() image.Rectangle {
	w := int(s.Config.Width)
	h := int(s.Config.Height)
	b := int(s.Config.Bins)

	return image.Rect(-20, -20, w+20, h+40+b)
}

// Generate generates a spectogram
func (s *Spectogram) Generate() (*Image128, error) {
	s.preEmphasis()

	img := NewImage128(s.imageBounds())
	draw.Draw(img, img.Bounds(), image.NewUniform(s.Config.BG0), image.ZP, draw.Src)

	return img, nil
}
