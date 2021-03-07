package gui

import (
	"errors"
	"fmt"
	"github.com/lxn/walk"
	"graphics-lab-2/imaging"
	"sort"
)

const (
	Name = iota
	Size
	XResolution
	YResolution
	ColorDepth
	CompressionType
)

type ImageMetadataModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*imaging.ImageMetadata
}

func NewImageMetadataModel(items []*imaging.ImageMetadata) *ImageMetadataModel {
	return &ImageMetadataModel {
		items: items,
	}
}

func (m *ImageMetadataModel) RowCount() int {
	return len(m.items)
}

func (m *ImageMetadataModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case Name:
		return item.Name
	case Size:
		return item.Size
	case XResolution:
		if item.XResolution == imaging.NoResolution {
			return "-"
		}
		return fmt.Sprintf("%d dpi", item.XResolution)
	case YResolution:
		if item.YResolution == imaging.NoResolution {
			return "-"
		}
		return fmt.Sprintf("%d dpi", item.YResolution)
	case ColorDepth:
		return fmt.Sprintf("%d bits", item.ColorDepth)
	case CompressionType:
		return item.CompressionType
	default:
		panic("incorrect column")
	}
}

func (m *ImageMetadataModel) Sort(col int, order walk.SortOrder) (err error) {
	if col == Size {
		return errors.New("can't sort on size")
	}

	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		alpha, beta := m.items[i], m.items[j]
		orderCond := func (ascCond bool) (res bool) {
			if m.sortOrder == walk.SortAscending {
				return ascCond
			}
			return !ascCond
		}
		switch m.sortColumn {
		case Name:
			return orderCond(alpha.Name < beta.Name)
		case XResolution:
			return orderCond(alpha.XResolution < beta.XResolution)
		case YResolution:
			return orderCond(alpha.YResolution < beta.YResolution)
		case ColorDepth:
			return orderCond(alpha.ColorDepth < beta.ColorDepth)
		case CompressionType:
			return orderCond(alpha.CompressionType < beta.CompressionType)
		default:
			panic("unreachable code")
		}
	})

	return m.SorterBase.Sort(col, order)
}

func (m *ImageMetadataModel) SetMetadata(metadata []*imaging.ImageMetadata) {
	m.items = metadata
	m.PublishRowsReset()

	m.Sort(Name, walk.SortAscending)
}