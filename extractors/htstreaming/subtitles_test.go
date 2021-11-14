package htstreaming

import (
	"strings"
	"testing"
)

func TestParseFirePlayerParams(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			jsTemplate string
			a          int
			c          int
			keywords   []string
		}
		want string
	}{
		{
			name: "Default",
			in: struct {
				jsTemplate string
				a          int
				c          int
				keywords   []string
			}{
				jsTemplate: `h E(F=g){1Q("15",{"16":{"17":"","18":""},"19":{"i":"","1a":"","1b":"1c-1d","14":2},"1e":[{"C":"b","i":"6:\\/\\/B.A.3\\/1g\\/c%z-%y%x%w%v%u%t%1h%p%1i%1j.s","r":"1k","q":g,"o":4},{"C":"b","i":"6:\\/\\/B.A.3\\/1f\\/c%z%y%x%w%v%u%t%L%p%N.s","r":"G","q":g,"o":2}],"b":{"U":"20","V":"W X"},"Y":"","Z":"\\\\n\\\\l\\\\m\\\\10\\\\8\\\\1p\\\\m\\\\1q\\\\8\\\\j\\\\1r\\\\1S\\\\7\\\\e\\\\1T\\\\f\\\\7\\\\9\\\\1U\\\\k\\\\7\\\\9\\\\1V\\\\f\\\\8\\\\e\\\\d\\\\9\\\\d\\\\l\\\\1X\\\\1Y\\\\d\\\\j\\\\21\\\\f\\\\7\\\\9\\\\8\\\\k\\\\n\\\\e\\\\29\\\\2a","2b":4,"1R":4,"1P":4,"1D":"c 1t: 1u 1v 1w 1x 1y 1z 1A 1","1B":2,"1s":2,"1C":{"1E":"1F","1G":[{"1H":"1I","1J":"6:\\/\\/a.1K.3\\/1L\\/1M","1N":"5"}],"12":"1O"},"28":2,"27":["26:\\/\\/25.24.3\\/"],"23":2,"22":"0","1Z":"6:\\/\\/1W.3\\/","13":4,"T":"1","S":"Q","P":"O+M=","K":"6:\\/\\/J.I.3\\/H\\/11.R"},2,F)}$(h(){$(1o).1n(h(){1m D=4;1l(D){E()}})});`,
				a:          62,
				c:          136,
				keywords:   strings.Split(`||false|com|true||https|x4e|x4d|x6d||captions|Onaho|x5a|x44|x78|null|function|file|x6a|x79|x47|x46|x4f|default|201|language|label|vtt|20Animation|20The|20Keikaku|20Ninshin|20Zenin|20Joshi|20Kyoushitsu|htstreaming|cdn|kind|fireplay|fireload|source|Spanish|libraries|jwplatform|content|jwPlayerURL|20Episodio|rZFhulEcXvUQMbyWAmIQyyjPjZAQPLw|20SubESP|ksaKvjlJRbnrPXSGpuPVqfscYS9|jwPlayerKey|jwplayer|js|videoPlayer|downloadType|fontSize|fontfamily|Trebuchet|MS|defaultImage|ck|x69|hDZaZjnc|admessage|downloadFile|active|f874601f91ab162335b5856b05987b7d|skin|name|url|logo|link|position|top|right|tracks|spanish|english|20Episode|20English|20Subbed|English|if|var|ready|document|x7a|x6c|x63|rememberPosition|Kyoushitsu|Joshi|Zenin|Ninshin|Keikaku|The|Animation|Episode|displaytitle|advertising|title|client|vast|schedule|offset|pre|tag|adtng|get|10012948|skipoffset|Reklam|jwplayer8quality|FirePlayer|jwplayer8button1|x34|x45|x49|x59|firevideoplayer|x52|x68|popurl||x41|poplimit|popactive|openwebtorrent|tracker|wss|p2pTrackers|p2p|x51|x3d|SubtitleManager`, "|"),
			},
			want: `function fireload(source=null){FirePlayer("f874601f91ab162335b5856b05987b7d",{"skin":{"name":"","url":""},"logo":{"file":"","link":"","position":"top-right","active":false},"tracks":[{"kind":"captions","file":"https:\/\/cdn.htstreaming.com\/english\/Onaho%20Kyoushitsu-%20Joshi%20Zenin%20Ninshin%20Keikaku%20The%20Animation%20Episode%201%20English%20Subbed.vtt","label":"English","language":null,"default":true},{"kind":"captions","file":"https:\/\/cdn.htstreaming.com\/spanish\/Onaho%20Kyoushitsu%20Joshi%20Zenin%20Ninshin%20Keikaku%20The%20Animation%20Episodio%201%20SubESP.vtt","label":"Spanish","language":null,"default":false}],"captions":{"fontSize":"20","fontfamily":"Trebuchet MS"},"defaultImage":"","ck":"\\x4f\\x47\\x46\\x69\\x4d\\x7a\\x46\\x6c\\x4d\\x6a\\x63\\x34\\x4e\\x44\\x45\\x78\\x4e\\x6d\\x49\\x79\\x4e\\x6d\\x59\\x78\\x4d\\x44\\x5a\\x6d\\x5a\\x47\\x52\\x68\\x5a\\x6a\\x41\\x78\\x4e\\x6d\\x4d\\x79\\x4f\\x44\\x51\\x3d","SubtitleManager":true,"jwplayer8button1":true,"jwplayer8quality":true,"title":"Onaho Kyoushitsu: Joshi Zenin Ninshin Keikaku The Animation Episode 1","displaytitle":false,"rememberPosition":false,"advertising":{"client":"vast","schedule":[{"offset":"pre","tag":"https:\/\/a.adtng.com\/get\/10012948","skipoffset":"5"}],"admessage":"Reklam"},"p2p":false,"p2pTrackers":["wss:\/\/tracker.openwebtorrent.com\/"],"popactive":false,"poplimit":"0","popurl":"https:\/\/firevideoplayer.com\/","downloadFile":true,"downloadType":"1","videoPlayer":"jwplayer","jwPlayerKey":"ksaKvjlJRbnrPXSGpuPVqfscYS9+rZFhulEcXvUQMbyWAmIQyyjPjZAQPLw=","jwPlayerURL":"https:\/\/content.jwplatform.com\/libraries\/hDZaZjnc.js"},false,source)}$(function(){$(document).ready(function(){var fireplay=true;if(fireplay){fireload()}})});`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := parseFirePlayerParams(tt.in.jsTemplate, tt.in.a, tt.in.c, tt.in.keywords)
			if out != tt.want {
				t.Errorf("Got: %v \n want: %v", out, tt.want)
			}
		})
	}
}

func TestParseSubtitles(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			name: "Default",
			in:   `function fireload(source=null){FirePlayer("f874601f91ab162335b5856b05987b7d",{"skin":{"name":"","url":""},"logo":{"file":"","link":"","position":"top-right","active":false},"tracks":[{"kind":"captions","file":"https:\/\/cdn.htstreaming.com\/english\/Onaho%20Kyoushitsu-%20Joshi%20Zenin%20Ninshin%20Keikaku%20The%20Animation%20Episode%201%20English%20Subbed.vtt","label":"English","language":null,"default":true},{"kind":"captions","file":"https:\/\/cdn.htstreaming.com\/spanish\/Onaho%20Kyoushitsu%20Joshi%20Zenin%20Ninshin%20Keikaku%20The%20Animation%20Episodio%201%20SubESP.vtt","label":"Spanish","language":null,"default":false}],"captions":{"fontSize":"20","fontfamily":"Trebuchet MS"},"defaultImage":"","ck":"\\x4f\\x47\\x46\\x69\\x4d\\x7a\\x46\\x6c\\x4d\\x6a\\x63\\x34\\x4e\\x44\\x45\\x78\\x4e\\x6d\\x49\\x79\\x4e\\x6d\\x59\\x78\\x4d\\x44\\x5a\\x6d\\x5a\\x47\\x52\\x68\\x5a\\x6a\\x41\\x78\\x4e\\x6d\\x4d\\x79\\x4f\\x44\\x51\\x3d","SubtitleManager":true,"jwplayer8button1":true,"jwplayer8quality":true,"title":"Onaho Kyoushitsu: Joshi Zenin Ninshin Keikaku The Animation Episode 1","displaytitle":false,"rememberPosition":false,"advertising":{"client":"vast","schedule":[{"offset":"pre","tag":"https:\/\/a.adtng.com\/get\/10012948","skipoffset":"5"}],"admessage":"Reklam"},"p2p":false,"p2pTrackers":["wss:\/\/tracker.openwebtorrent.com\/"],"popactive":false,"poplimit":"0","popurl":"https:\/\/firevideoplayer.com\/","downloadFile":true,"downloadType":"1","videoPlayer":"jwplayer","jwPlayerKey":"ksaKvjlJRbnrPXSGpuPVqfscYS9+rZFhulEcXvUQMbyWAmIQyyjPjZAQPLw=","jwPlayerURL":"https:\/\/content.jwplatform.com\/libraries\/hDZaZjnc.js"},false,source)}$(function(){$(document).ready(function(){var fireplay=true;if(fireplay){fireload()}})});`,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := parseCaptions(tt.in)
			if len(out) != tt.want {
				t.Errorf("Got: %v \n want: %v", len(out), tt.want)
			}
		})
	}
}
