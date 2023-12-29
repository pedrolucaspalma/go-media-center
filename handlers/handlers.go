package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var base_path = "/home/pedro/ubuntu-nas/media"

var shortVideoPath = base_path + "/Jovem-Nerd/NerdOffice/S01/[00001] - Piloto [XkG1F-wJ4Is].webm"
var longVideoPath = base_path + "/movies/Top Gun Maverick (2022) [IMAX] [REPACK] [1080p] [WEBRip] [5.1] [YTS.MX]/Top.Gun.Maverick.2022.IMAX.REPACK.1080p.WEBRip.x264.AAC5.1-[YTS.MX].mp4"
var longVideoPath2 = base_path + "/movies/John Wick Chapter 4 (2023) [1080p] [WEBRip] [5.1] [YTS.MX]/John.Wick.Chapter.4.2023.1080p.WEBRip.x264.AAC5.1-[YTS.MX].mp4"

var listaOpcoes = [...]string{shortVideoPath, longVideoPath, longVideoPath2}
var selectedFile = listaOpcoes[2]

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Foi a home")
}

func PlayerHandler(res http.ResponseWriter, req *http.Request) {
	template := template.Must(template.ParseFiles("./templates/player.html"))
	template.Execute(res, nil)
}

func HandleVideo(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(selectedFile)
	if err != nil {
		http.Error(w, "Error opening video file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting video file information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		byteRanges, err := parseRangeHeader(rangeHeader, fileInfo.Size())
		if err != nil {
			http.Error(w, "Invalid Range Header", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		firstRange := byteRanges[0]

		fmt.Printf("%d \n", firstRange)

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", firstRange.start, firstRange.end, fileInfo.Size()))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", firstRange.length))
		w.WriteHeader(http.StatusPartialContent)

		file.Seek(firstRange.start, 0)

		http.ServeContent(w, r, "", fileInfo.ModTime(), file)
	} else {
		http.ServeContent(w, r, "", fileInfo.ModTime(), file)
	}
}

// Range represents a byte range
type Range struct {
	start  int64
	end    int64
	length int64
}

// parseRangeHeader parses the Range header and returns a list of byte ranges
func parseRangeHeader(rangeHeader string, fileSize int64) ([]Range, error) {
	const prefix = "bytes="
	if !strings.HasPrefix(rangeHeader, prefix) {
		return nil, fmt.Errorf("Invalid Range Header")
	}

	rangeSpecs := strings.Split(rangeHeader[len(prefix):], ",")
	var byteRanges []Range

	for _, rangeSpec := range rangeSpecs {
		var start, end int64
		var err error

		if strings.HasSuffix(rangeSpec, "-") {
			// Range like "start-"
			start, err = strconv.ParseInt(rangeSpec[:len(rangeSpec)-1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid Range Header")
			}
			end = fileSize - 1
		} else if strings.Contains(rangeSpec, "-") {
			// Range like "start-end"
			rangeParts := strings.Split(rangeSpec, "-")
			start, err = strconv.ParseInt(rangeParts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid Range Header")
			}
			end, err = strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid Range Header")
			}
		} else {
			return nil, fmt.Errorf("Invalid Range Header")
		}

		// Ensure valid ranges
		if start < 0 || end < start || end >= fileSize {
			return nil, fmt.Errorf("Invalid Range Header")
		}

		byteRanges = append(byteRanges, Range{start: start, end: end, length: end - start + 1})
	}

	return byteRanges, nil
}
