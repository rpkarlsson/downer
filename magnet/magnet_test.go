package magnet

import (
	"testing"
)

const expected = "d10:magnet-uri60:magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88ae"
const magnetLink = "magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a"

func TestIsMagnetLink(t *testing.T) {
	if !IsMagnetLink(magnetLink) {
		t.Fatal("Should recognize strings starting with \"magnet:\"")
	}

	if IsMagnetLink("http://gnu.org/magnet:") {
		t.Fatal("Should not recognize string not starting with \"magnet:\"")
	}
}

func TestMakeTorrent(t *testing.T) {
	out := MakeTorrentBody(magnetLink)
	if out != expected {
		t.Fatal("Wrong magnet link was built")
	}
}
