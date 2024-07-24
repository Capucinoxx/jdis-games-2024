package utils

import (
	"crypto/sha256"
	"math"
	"math/rand"
)

func SafeClose[T any](ch chan T) {
	defer func() {
		recover()
	}()

	close(ch)
}

func Shuffle[S ~[]E, E any](r *rand.Rand, s S) {
	r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func NilIf[T any](v *T, b bool) *T {
	if b {
		return nil
	}
	return v
}

func Round(f float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(f*pow) / pow
}

func ToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func hslToRGB(h, s, l float64) (int, int, int) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2

	var r, g, b float64

	if 0 <= h && h < 60 {
		r, g, b = c, x, 0
	} else if 60 <= h && h < 120 {
		r, g, b = x, c, 0
	} else if 120 <= h && h < 180 {
		r, g, b = 0, c, x
	} else if 180 <= h && h < 240 {
		r, g, b = 0, x, c
	} else if 240 <= h && h < 300 {
		r, g, b = x, 0, c
	} else if 300 <= h && h < 360 {
		r, g, b = c, 0, x
	}

	return int((r + m) * 255), int((g + m) * 255), int((b + m) * 255)
}

func NameColor(name string) int32 {
	hasher := sha256.New()
	hasher.Write([]byte(name))

	var hashArray [32]byte
	copy(hashArray[:], hasher.Sum(nil))

	hash := (int(hashArray[0]) + int(hashArray[1])*256 + int(hashArray[2])*256*256)
	hue := hash % 360
	saturation := 0.6
	lightness := 0.7

	r, g, b := hslToRGB(float64(hue), saturation, lightness)
	return int32((r << 16) | (g << 8) | b)
}
