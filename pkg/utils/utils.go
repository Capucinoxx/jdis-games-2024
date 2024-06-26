package utils

import (
  "crypto/sha256"
	"math/rand"
)

func Shuffle[S ~[]E, E any](r *rand.Rand, s S) {
  r.Shuffle(len(s), func(i, j int) {
    s[i], s[j] =  s[j], s[i]
  })
}

func ToInt(b bool) int {
  if b {
    return 1
  }
  return 0
}


func hueToRGB(p, q, t float64) float64 {
  if t < 0 {
    t += 1
  }
  if t > 1 {
    t -= 1
  }
  if t < 1/6 {
    return p + (q - p) * 6 * t
  }
  if t < 1/2 {
    return q
  }
  if t < 2/3 {
    return p + (q - p) * (2/3 - t) * 6
  }
  return p
}


func HslToRGB(h, s, l float64) (int, int, int) {
  var r, g, b float64

  if s == 0 {
    r = l
    g = l
    b = l
  } else {
    var q float64
    if l < 0.5 {
      q = l * (1 + s)
    } else {
      q = l + s - l * s
    }
    p := 2 * l - q
    r = hueToRGB(p, q, h + 1/3)
    g = hueToRGB(p, q, h)
    b = hueToRGB(p, q, h - 1/3)
  }

  return int(r * 255), int(g * 255), int(b * 255)
}


func NameColor(name string) int32 {
  hasher := sha256.New()
  hasher.Write([]byte(name))
  hash := hasher.Sum(nil)

  var hashBytes [32]byte
  copy(hashBytes[:], hash)
  
  hue := (int(hashBytes[0]) + (int(hashBytes[1]) * 0xff) + (int(hashBytes[2]) * 256 * 256)) % 360
  s := 0.7
  l := 0.5

  r, g, b := HslToRGB(float64(hue), s, l)
  return int32(r << 16 + g << 8 + b)
}
