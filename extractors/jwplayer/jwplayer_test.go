package jwplayer

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		Args test.Args
	}{
		{
			Name: "Single Episode hentaihaven.com",
			Args: test.Args{
				URL:     "https://hentaihaven.com/wp-content/plugins/player-logic/player.php?data=RWsra254MXRSelpiUFRnNmNESjJ5ZXg2QWlHQ0xGOE9TVHRrcU9LeTUxQXN3NmJ5NVNZaXNlb015OVQ5THBjb1M0eGc2dFVOR2hHdDNNRjJFVnd6KzFrUVJGTlBkOG9hWHNlUU1nTStJRFI3QkF4ZFNxaFY1QVFTandLWGo2ZXVZeXpkR1JPMTMxcFdVRkp2TWIwTGMwZTBCZW55M053aGtkTEtUcVNFWmtCK25DYmRIc0hmcHF6ZTQvRXQ5WHNIRkQ5RGlMOEdVdS85cXBKK3ZkQVZ6VkZDbE1Vd0o2bU9pZmJhVzY1UHdJVDdjL2I4ZThWdVhNUWllUEgya29IY3c0aGQvWkpYYlhMQXdPTlduRjJ3UFE9PTp8Ojp8OjJrTG42T2ZSbEM5L3FiUnp6Mmphb2c9PQ==",
				Title:   "jwplayer video",
				Quality: "1920x1080",
				Size:    465565080,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}

func TestFindJWPlayerURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want string
	}{
		{
			Name: "HTML string",
			URL:  `<iframe src="https://hentaihaven.com/wp-content/plugins/player-logic/player.php?data=RWsra254MXRSelpiUFRnNmNESjJ5ZXg2QWlHQ0xGOE9TVHRrcU9LeTUxQXN3NmJ5NVNZaXNlb015OVQ5THBjb1M0eGc2dFVOR2hHdDNNRjJFVnd6KzFrUVJGTlBkOG9hWHNlUU1nTStJRFI3QkF4ZFNxaFY1QVFTandLWGo2ZXVZeXpkR1JPMTMxcFdVRkp2TWIwTGMwZTBCZW55M053aGtkTEtUcVNFWmtCK25DYmRIc0hmcHF6ZTQvRXQ5WHNIRkQ5RGlMOEdVdS85cXBKK3ZkQVZ6VkZDbE1Vd0o2bU9pZmJhVzY1UHdJVDdjL2I4ZThWdVhNUWllUEgya29IY3c0aGQvWkpYYlhMQXdPTlduRjJ3UFE9PTp8Ojp8OjJrTG42T2ZSbEM5L3FiUnp6Mmphb2c9PQ==" frameborder="0" scrolling="no" allowfullscreen=""></iframe>`,
			Want: "https://hentaihaven.com/wp-content/plugins/player-logic/player.php?data=RWsra254MXRSelpiUFRnNmNESjJ5ZXg2QWlHQ0xGOE9TVHRrcU9LeTUxQXN3NmJ5NVNZaXNlb015OVQ5THBjb1M0eGc2dFVOR2hHdDNNRjJFVnd6KzFrUVJGTlBkOG9hWHNlUU1nTStJRFI3QkF4ZFNxaFY1QVFTandLWGo2ZXVZeXpkR1JPMTMxcFdVRkp2TWIwTGMwZTBCZW55M053aGtkTEtUcVNFWmtCK25DYmRIc0hmcHF6ZTQvRXQ5WHNIRkQ5RGlMOEdVdS85cXBKK3ZkQVZ6VkZDbE1Vd0o2bU9pZmJhVzY1UHdJVDdjL2I4ZThWdVhNUWllUEgya29IY3c0aGQvWkpYYlhMQXdPTlduRjJ3UFE9PTp8Ojp8OjJrTG42T2ZSbEM5L3FiUnp6Mmphb2c9PQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			u := FindJWPlayerURL(&tt.URL)
			if u == "" {
				t.Errorf("Got: %v - Want: %v", u, tt.Want)
			}
		})
	}
}
