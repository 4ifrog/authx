package avatar

import (
	"math/rand"
	"time"

	"github.com/stdatiks/jdenticon-go"
)

// avatarSize returns avatar's size. avatarSize = Height = Width.
const avatarSize = 80

//nolint:gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetIdenticonWithSize generates and returns an avatar icon unique to an identity eg. user's email
// or username.
func GetIdenticonWithSize(identity string, size int) ([]byte, error) {
	cfg := jdenticon.Config{
		Hues:       jdenticon.DefaultConfig.Hues,
		Colored:    jdenticon.DefaultConfig.Colored,
		Grayscale:  jdenticon.DefaultConfig.Grayscale,
		Background: jdenticon.DefaultConfig.Background,
		Width:      size,
		Height:     size,
		Padding:    jdenticon.DefaultConfig.Padding,
	}

	icon := jdenticon.NewWithConfig(identity, &cfg)
	return icon.SVG()
}

// GetIdenticon generates an avatar with the standard size.
func GetIdenticon(identity string) ([]byte, error) {
	return GetIdenticonWithSize(identity, avatarSize)
}
