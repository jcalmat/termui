// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"

	. "github.com/jcalmat/termui/v3"
)

type Gauge struct {
	Block
	Percent    int
	BarColor   Color
	Label      string
	LabelStyle Style
	Type       GaugeType
	FillType   GaugeFillType
}

type GaugeType int
type GaugeFillType int

const (
	GaugeNatural GaugeFillType = iota
	GaugeReverse
)

const (
	GaugeHorizontal GaugeType = iota
	GaugeVertical
)

func NewGauge() *Gauge {
	return &Gauge{
		Block:      *NewBlock(),
		BarColor:   Theme.Gauge.Bar,
		LabelStyle: Theme.Gauge.Label,
	}
}

func (self *Gauge) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	label := self.Label
	if label == "" {
		label = fmt.Sprintf("%d%%", self.Percent)
	}

	barWidth := int((float64(self.Percent) / 100) * float64(self.Inner.Dx()))
	barHeight := int((float64(self.Percent) / 100) * float64(self.Inner.Dy()))

	// plot bar
	switch self.Type {
	case GaugeVertical:
		switch self.FillType {
		case GaugeReverse:
			buf.Fill(
				NewCell(' ', NewStyle(ColorClear, self.BarColor)),
				image.Rect(self.Inner.Min.X, self.Inner.Min.Y, self.Inner.Max.X, self.Inner.Min.Y+barHeight),
			)
		case GaugeNatural:
			buf.Fill(
				NewCell(' ', NewStyle(ColorClear, self.BarColor)),
				image.Rect(self.Inner.Min.X, self.Inner.Max.Y, self.Inner.Max.X, self.Inner.Max.Y-barHeight),
			)
		}
	case GaugeHorizontal:
		switch self.FillType {
		case GaugeReverse:
			buf.Fill(
				NewCell(' ', NewStyle(ColorClear, self.BarColor)),
				image.Rect(self.Inner.Max.X-barWidth, self.Inner.Min.Y, self.Inner.Max.X, self.Inner.Max.Y),
			)
		case GaugeNatural:
			buf.Fill(
				NewCell(' ', NewStyle(ColorClear, self.BarColor)),
				image.Rect(self.Inner.Min.X, self.Inner.Min.Y, self.Inner.Min.X+barWidth, self.Inner.Max.Y),
			)
		}
	}

	// plot label
	labelXCoordinate := self.Inner.Min.X + (self.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := self.Inner.Min.Y + ((self.Inner.Dy() - 1) / 2)
	if labelYCoordinate < self.Inner.Max.Y {
		for i, char := range label {
			style := self.LabelStyle
			switch self.Type {
			case GaugeHorizontal:
				switch self.FillType {
				case GaugeNatural:
					if labelXCoordinate+i+1 <= self.Inner.Min.X+barWidth {
						style = NewStyle(self.BarColor, ColorClear, ModifierReverse)
					}
				case GaugeReverse:
					if labelXCoordinate+i+1 >= self.Inner.Max.X-barWidth {
						style = NewStyle(self.BarColor, ColorClear, ModifierReverse)
					}
				}
			case GaugeVertical:
				switch self.FillType {
				case GaugeNatural:
					if labelYCoordinate >= self.Inner.Max.Y-barHeight {
						style = NewStyle(self.BarColor, ColorClear, ModifierReverse)
					}
				case GaugeReverse:
					if labelYCoordinate < self.Inner.Min.Y+barHeight {
						style = NewStyle(self.BarColor, ColorClear, ModifierReverse)
					}
				}
			}
			buf.SetCell(NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}
