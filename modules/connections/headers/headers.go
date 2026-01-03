package headers

import (
	"fmt"
	"net/http"
)

func GetLivesearchHeaders(poesessid string) http.Header {
	head := http.Header{}
	head.Add("Accept", "*/*")
	head.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	head.Add("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	head.Add("Cache-Control", "no-cache")
	head.Add("Cookie", fmt.Sprintf("POESESSID=%v", poesessid))
	head.Add("Host", "www.pathofexile.com")
	head.Add("Origin", "https://www.pathofexile.com")
	head.Add("Pragma", "no-cache")
	head.Add("Sec-Fetch-Dest", "empty")
	head.Add("Sec-Fetch-Mode", "websocket")
	head.Add("Sec-Fetch-Site", "same-origin")
	head.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0")

	return head
}

func GetFetchitemHeaders(poesessid string) http.Header {
	head := http.Header{}
	head.Add("Cookie", fmt.Sprintf("POESESSID=%v", poesessid))
	head.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	head.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	head.Add("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	head.Add("Connection", "keep-alive")
	head.Add("Host", "www.pathofexile.com")
	head.Add("Priority", "u=0, i")
	head.Add("Sec-Fetch-Dest", "document")
	head.Add("Sec-Fetch-Mode", "navigate")
	head.Add("Sec-Fetch-Site", "none")
	head.Add("Sec-Fetch-User", "?1")
	head.Add("TE", "trailers")
	head.Add("Upgrade-Insecure-Requests", "1")
	head.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0")

	return head
}

func GetWhisperHeaders(poesessid string) http.Header {
	head := http.Header{}
	head.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0")
	head.Add("Accept", "*/*")
	head.Add("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	head.Add("Content-Type", "application/json")
	head.Add("X-Requested-With", "XMLHttpRequest")
	head.Add("Sec-Fetch-Dest", "empty")
	head.Add("Sec-Fetch-Mode", "cors")
	head.Add("Sec-Fetch-Site", "same-origin")
	head.Add("Cookie", fmt.Sprintf("POESESSID=%v", poesessid))
	return head
}
