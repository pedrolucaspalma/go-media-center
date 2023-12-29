package handlers

import (
	"fmt"
	"html/template"
	"log"
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

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	filePtr, err := os.Open(selectedFile)
	if err != nil {
		log.Fatal(w, "N leu o arquivo")
	}

	defer filePtr.Close()

	fileInfo, _ := os.Stat(selectedFile)

	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Content-type", "video/mp4")
	w.Header().Set("Accept-ranges", "bytes")
	w.Header().Set("Content-length", fmt.Sprint(fileInfo.Size()))

	buffer := make([]byte, 1024)

	for {
		n, err := filePtr.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		if n == 0 {
			break
		}

		w.Write(buffer[:n])

		w.(http.Flusher).Flush()
	}
}

func HandleVideo(w http.ResponseWriter, r *http.Request) {
	// Open the video file
	file, err := os.Open(selectedFile)
	if err != nil {
		http.Error(w, "Error opening video file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get the file's information
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting video file information", http.StatusInternalServerError)
		return
	}

	// Set headers for video streaming
	// w.Header().Set("Content-Type", "video/mp4")
	// w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Check for range request
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Parse range header
		byteRanges, err := parseRangeHeader(rangeHeader, fileInfo.Size())
		if err != nil {
			http.Error(w, "Invalid Range Header", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Handle the first range only (ignoring multiple ranges for simplicity)
		firstRange := byteRanges[0]

		fmt.Printf("%d \n", firstRange)

		// Set headers for partial content
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", firstRange.start, firstRange.end, fileInfo.Size()))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", firstRange.length))
		w.WriteHeader(http.StatusPartialContent)

		// Seek to the start position in the file
		file.Seek(firstRange.start, 0)

		// Copy the specified range to the response writer
		http.ServeContent(w, r, "", fileInfo.ModTime(), file)
	} else {
		// Serve the entire video file
		fmt.Print("Else")
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
