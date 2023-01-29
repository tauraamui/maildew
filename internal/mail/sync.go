package mail

import (
	"github.com/tauraamui/gonp"
)

func ResolveNewAndMissing(a, b []uint32) ([]uint32, []uint32) {
	return diffUIDs(a, b)
}

func diffUIDs(a, b []uint32) (additions, missing []uint32) {
	diff := gonp.New(a, b)
	diff.Compose()
	ses := diff.Ses()

	for _, e := range ses {
		switch e.GetType() {
		case gonp.SesAdd:
			additions = append(additions, e.GetElem())
		case gonp.SesDelete:
			missing = append(missing, e.GetElem())
		}
	}

	return additions, missing
}
