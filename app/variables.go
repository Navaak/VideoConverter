package app

const (
	smilHead = `<?xml version="1.0" encoding="UTF-8"?><smil title=""><body> <switch>`

	smilFooter = `</switch></body></smil>`
)

var (
	smilQualities = map[int]string{
		240: `<video height="240" src="%s" systemLanguage="eng" width="426">
                <param name="videoBitrate" value="300000" valuetype="data"></param>
                <param name="audioBitrate" value="44100" valuetype="data"></param>
            </video>`,
		360: `<video height="360" src="%s" systemLanguage="eng" width="640">
                <param name="videoBitrate" value="400000" valuetype="data"></param>
                <param name="audioBitrate" value="44100" valuetype="data"></param>
            </video>`,
		480: ` <video height="480" src="%s" systemLanguage="eng" width="854">
                <param name="videoBitrate" value="500000" valuetype="data"></param>
                <param name="audioBitrate" value="44100" valuetype="data"></param>
            </video>`,
		720: `<video height="720" src="%s" systemLanguage="eng" width="1280">
                <param name="videoBitrate" value="1500000" valuetype="data"></param>
                <param name="audioBitrate" value="44100" valuetype="data"></param></video>`,
		1080: `<video height="1080" src="%s" systemLanguage="eng" width="1920">
<param name="videoBitrate" value="3000000" valuetype="data"></param>
<param name="audioBitrate" value="44100" valuetype="data"></param></video>`,
	}

	bitRate = map[int]int{
		240:  300000,
		360:  400000,
		480:  500000,
		720:  1500000,
		1080: 3000000,
	}
)
