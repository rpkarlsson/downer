package magnet

import (
	"fmt"
	"regexp"
)

const torrentTemplate = "d10:magnet-uri%d:%se"

func MakeTorrentBody(URI string) string {
	return fmt.Sprintf(torrentTemplate, len(URI), URI)
}

func IsMagnetLink(URI string) bool {
	matched, _ := regexp.MatchString(`^magnet:*`, URI)
	return matched
}
