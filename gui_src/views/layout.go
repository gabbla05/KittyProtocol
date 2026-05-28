package views

import (
	"fyne.io/fyne/v2"
)

// formLayout wymusza szerokość 350px, ale wysokość liczy dynamicznie
type formLayout struct{}

func (f *formLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for _, child := range objects {
		childHeight := child.MinSize().Height
		child.Resize(fyne.NewSize(350, childHeight))
		child.Move(pos)
		pos.Y += childHeight + 5 // 5px odstępu
	}
}

func (f *formLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var totalHeight float32
	for _, child := range objects {
		totalHeight += child.MinSize().Height + 5
	}
	// Zwracamy dokładnie tyle wysokości, ile potrzeba (nie sztywne 600!)
	return fyne.NewSize(350, totalHeight) 
}