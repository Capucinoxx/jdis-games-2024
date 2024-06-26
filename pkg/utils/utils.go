package utils

import "math/rand"

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
