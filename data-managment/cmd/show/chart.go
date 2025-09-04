package show

import (
	"data-managment/util/repo"
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

type barChart struct {
	widget.BaseWidget
	data []genePair
}

func draw(patientData repo.PatientData, geneData genesExpressions) error {
	a := app.New()
	w := a.NewWindow("Gene Expression Viewer")
	w.Resize(fyne.NewSize(800, 600))

	patientInfo := widget.NewForm(
		widget.NewFormItem("Patient Barcode", widget.NewLabel(patientData.BCRPatientBarcode)),
		widget.NewFormItem("Clinical Stage", widget.NewLabel(patientData.ClinicalStage)),
		widget.NewFormItem("DSS", widget.NewLabel(fmt.Sprintf("%t", patientData.DSS))),
		widget.NewFormItem("OS", widget.NewLabel(fmt.Sprintf("%t", patientData.OS))),
	)

	chart := newBarChart(geneData.Genes)

	content := container.NewBorder(patientInfo, nil, nil, nil, chart)

	w.SetContent(content)

	zap.S().Info("Press q to quit")
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		zap.S().Debugf("Event %+v", *keyEvent)
		if keyEvent.Name == fyne.KeyQ {
			a.Quit()
		}
	})

	w.ShowAndRun()
	return nil
}

func newBarChart(data []genePair) *barChart {
	b := &barChart{data: data}
	b.ExtendBaseWidget(b)
	return b
}

func (b *barChart) CreateRenderer() fyne.WidgetRenderer {
	r := &barChartRenderer{
		chart:    b,
		bars:     []fyne.CanvasObject{},
		labels:   []fyne.CanvasObject{},
		dividers: []fyne.CanvasObject{},
	}
	r.Refresh()
	return r
}

type barChartRenderer struct {
	chart    *barChart
	bars     []fyne.CanvasObject
	labels   []fyne.CanvasObject
	dividers []fyne.CanvasObject
}

func (r *barChartRenderer) Layout(size fyne.Size) {
	if len(r.chart.data) == 0 {
		return
	}

	minVal, maxVal := 0.0, 0.0
	for _, d := range r.chart.data {
		if d.Expression < minVal {
			minVal = d.Expression
		}
		if d.Expression > maxVal {
			maxVal = d.Expression
		}
	}

	totalRange := maxVal - minVal
	if totalRange == 0 {
		totalRange = 1
	}

	barWidth := (size.Width - 20) / float32(len(r.chart.data))
	zeroY := size.Height * (float32(maxVal) / float32(totalRange))
	labelHeight := float32(20)

	for i, d := range r.chart.data {
		bar := r.bars[i].(*canvas.Rectangle)
		label := r.labels[i].(*canvas.Text)

		x := 10 + float32(i)*barWidth
		barHeight := (size.Height - labelHeight) * float32(math.Abs(d.Expression)) / float32(totalRange)

		if d.Expression >= 0 {
			bar.Move(fyne.NewPos(x, zeroY-barHeight))
			bar.Resize(fyne.NewSize(barWidth-5, barHeight))
		} else {
			bar.Move(fyne.NewPos(x, zeroY))
			bar.Resize(fyne.NewSize(barWidth-5, barHeight))
		}

		label.Move(fyne.NewPos(x, size.Height-labelHeight))
		label.Resize(fyne.NewSize(barWidth-5, labelHeight))
	}

	if len(r.dividers) > 0 {
		r.dividers[0].Move(fyne.NewPos(0, zeroY))
		r.dividers[0].Resize(fyne.NewSize(size.Width, 1))
	}
}

func (r *barChartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

func (r *barChartRenderer) Refresh() {
	r.bars = []fyne.CanvasObject{}
	r.labels = []fyne.CanvasObject{}
	r.dividers = []fyne.CanvasObject{}

	for _, d := range r.chart.data {
		bar := canvas.NewRectangle(color.NRGBA{B: 200, A: 255})
		if d.Expression < 0 {
			bar.FillColor = color.NRGBA{R: 200, A: 255}
		}
		r.bars = append(r.bars, bar)

		label := canvas.NewText(d.Gene, color.White)
		label.Alignment = fyne.TextAlignCenter
		r.labels = append(r.labels, label)
	}

	zeroLine := canvas.NewRectangle(color.White)
	r.dividers = append(r.dividers, zeroLine)

	r.Layout(r.chart.Size())
	canvas.Refresh(r.chart)
}

// Objects returns all objects that should be drawn.
func (r *barChartRenderer) Objects() []fyne.CanvasObject {
	return append(append(r.bars, r.labels...), r.dividers...)
}

// Destroy is for cleanup.
func (r *barChartRenderer) Destroy() {}
