package qsv

import (
	"math/rand"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

func pivotPickerFirst[T constraints.Ordered](s []T, l, h int) T {
	return s[l]
}

func swap[T constraints.Ordered](s []T, i, j int, delay time.Duration) {
	time.Sleep(delay)
	s[i], s[j] = s[j], s[i]
}

func partition[T constraints.Ordered](s []T, l, h int, pivotPicker func(s []T, l, h int) T, swapDelay time.Duration) int {
	piv := pivotPicker(s, l, h)

	lI := l - 1
	rI := h + 1

	for {
		for lI += 1; s[lI] < piv; lI += 1 {
		}
		for rI -= 1; s[rI] > piv; rI -= 1 {
		}

		if lI >= rI {
			return rI
		}

		swap(s, lI, rI, swapDelay)
	}
}

func quickSort[T constraints.Ordered](s []T, l, h int, swapDelay time.Duration) {
	if l >= 0 && h >= 0 && l < h {
		p := partition(s, l, h, pivotPickerFirst, swapDelay)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			quickSort(s, l, p, swapDelay)
		}()
		go func() {
			defer wg.Done()
			quickSort(s, p+1, h, swapDelay)
		}()
		wg.Wait()
	}
}

//func quickSortSteped[T constraints.Ordered](s []T, l, h int) {
//	stepCh := make(chan struct{})
//	step := func() { stepCh <- struct{}{} }
//
//	var qs func([]T, int, int)
//	qs = func(s []T, l, h int) {
//		if l >= 0 && h >= 0 && l < h {
//			p := partition(s, l, h, pivotPickerFirst)
//			qs(s, l, p)
//			qs(s, p+1, h)
//		}
//	}
//
//}

func sliceVBox(w, h, ww, wh, py int32) rl.Rectangle {
	box := rl.Rectangle{
		Width:  float32(w),
		Height: float32(h),
	}
	box.X = float32((ww - box.ToInt32().Width) / 2)
	box.Y = float32(wh - box.ToInt32().Height - py)

	return box
}

func drawSlice[T constraints.Float | constraints.Integer](box rl.Rectangle, boxLineThick float32, s []T, sliceMax T) {
	ehm := box.Height / 100
	ehMin := box.Height / 100
	ehMax := box.Height - (boxLineThick * 2) - ehm

	ew := box.Width / float32(len(s))
	ewm := ew / 4
	ew = ew - ewm - (ewm / float32(len(s)))
	ex := box.X + ewm

	for _, v := range s {
		eh := (ehMax-ehMin)*float32(v)/float32(sliceMax) + ehMin
		ey := box.Y - boxLineThick + box.Height - eh
		rl.DrawRectangleLinesEx(rl.Rectangle{Width: ew, Height: eh, X: ex, Y: ey}, 1, rl.Black)
		ex += ew + ewm
	}
}

func Run() error {
	var ww int32 = 1240
	var wh int32 = 720
	sLen := 10_000
	sMaxVal := 1000
	sBoxLineThick := float32(3)

	swapDelay := 1 * time.Millisecond

	s := make([]int, sLen)
	for i := 0; i < sLen; i++ {
		s[i] = rand.Intn(sMaxVal + 1)
	}

	go func() {
		time.Sleep(3 * time.Second)
		quickSort(s, 0, len(s)-1, swapDelay)
	}()

	rl.InitWindow(ww, wh, "Quick sort visualization")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	sliceBox := sliceVBox(ww-40, wh/2, ww, wh, 10)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawRectangleLinesEx(sliceBox, sBoxLineThick, rl.LightGray)
		drawSlice(sliceBox, sBoxLineThick, s, sMaxVal)

		rl.EndDrawing()
	}

	return nil
}
