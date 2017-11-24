package res

import (
	_ "image/png"
)

func Load(name string) ([]byte, error) {
	return Asset(name)
}
