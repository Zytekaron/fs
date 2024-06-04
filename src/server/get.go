package server

import (
	"bytes"
	"errors"
	"fs/src/cache"
	"fs/src/config"
	"fs/src/server/response"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var mdParser = parser.New()
var htmlRenderer = html.NewRenderer(html.RendererOptions{
	Flags: html.CommonFlags | html.HrefTargetBlank,
})

func readFileHandler(cfg *config.Config, c *cache.Cache) http.Handler {
	baseDir := cfg.Files.Directory

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Clean(r.URL.Path)
		if strings.Contains(filePath, "..") {
			response.WriteErrorReason(w, http.StatusBadRequest, "path must not contain ..")
			return
		}

		q := r.URL.Query()
		viewas := q.Get("viewas")

		// select browser rendering mode as a default in browsers
		if len(viewas) == 0 && isBrowser(r.UserAgent()) {
			viewas = "browser"
		}

		// check for a cached entry

		if data, ok := c.Get(filePath); ok {
			if data.IsDir {
				readFileHandlerSendDir(w, *data.BaseDir, *data.DirEntries, viewas)
			} else {
				size := int64(len(*data.FileBytes))
				reader := bytes.NewReader(*data.FileBytes)
				readFileHandlerSendFile(w, reader, size, viewas)
			}
			return
		}

		// open the file/directory from fs and gather info about it

		fullPath := path.Join(baseDir, filePath)
		file, err := os.Open(fullPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				response.WriteError(w, http.StatusNotFound)
				return
			}
			log.Println("error reading file:", err)
			response.WriteErrorReason(w, http.StatusInternalServerError, "error reading file")
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Println("error statting file:", err)
			response.WriteErrorReason(w, http.StatusInternalServerError, "error statting file")
		}

		// conditionally cache, then serve the file/directory

		if stat.IsDir() {
			dir, err := file.ReadDir(-1)
			if err != nil {
				log.Println("error reading directory:", err)
				response.WriteErrorReason(w, http.StatusInternalServerError, "error statting file")
				return
			}
			slices.SortFunc(dir, func(a, b os.DirEntry) int {
				if a.IsDir() && !b.IsDir() {
					return -1
				}
				if !a.IsDir() && b.IsDir() {
					return 1
				}
				return strings.Compare(a.Name(), b.Name())
			})

			c.Put(filePath, &cache.FileData{
				DirEntries: &dir,
				BaseDir:    &filePath,
				IsDir:      true,
			})

			readFileHandlerSendDir(w, filePath, dir, viewas)
		} else {
			size := stat.Size()
			reader := io.Reader(file)

			// if the file is small enough to cache, read
			// it into memory now and swap out the reader
			if size <= cfg.Cache.MaxEntry {
				data, err := io.ReadAll(reader)
				if err != nil {
					log.Println("error reading file:", err)
					response.WriteErrorReason(w, http.StatusInternalServerError, "error reading file")
					return
				}

				c.Put(filePath, &cache.FileData{
					FileBytes: &data,
				})

				reader = bytes.NewReader(data)
			}

			readFileHandlerSendFile(w, reader, size, viewas)
		}
	})
}

func readFileHandlerSendDir(w http.ResponseWriter, baseDir string, entries []fs.DirEntry, viewas string) error {
	switch viewas {
	case "json_obj":
		return response.WriteJSON(w, http.StatusOK, map[string]any{
			"entries": makeList(entries),
		})
	case "json":
		return response.WriteJSON(w, http.StatusOK, makeList(entries))
	case "hl", "browser":
		list := makeListHyperlinks(baseDir, entries)
		return response.WriteData(w, http.StatusOK, []byte(list))
	default:
		list := strings.Join(makeList(entries), "\n")
		return response.WriteData(w, http.StatusOK, []byte(list))
	}
}

func makeList(entries []fs.DirEntry) []string {
	list := make([]string, len(entries))
	var buf strings.Builder
	for i, entry := range entries {
		buf.WriteString(entry.Name())
		if entry.IsDir() {
			buf.WriteByte('/')
		}
		list[i] = buf.String()
		buf.Reset()
	}
	return list
}

func makeListHyperlinks(baseDir string, entries []fs.DirEntry) string {
	var buf strings.Builder
	for _, entry := range entries {
		buf.WriteString(`<a href="`)
		buf.WriteString(path.Join(baseDir, entry.Name()))
		buf.WriteString(`">`)
		buf.WriteString(entry.Name())
		if entry.IsDir() {
			buf.WriteByte('/')
		}
		buf.WriteString("</a><br>")
	}
	return buf.String()
}

func readFileHandlerSendFile(w http.ResponseWriter, rd io.Reader, size int64, viewas string) error {
	switch viewas {
	case "md", "markdown":
		// impose a maximum file size of 1 MB
		if size > 1024*1024 {
			// "i aint reading allat"
			return response.WriteErrorReason(w, http.StatusBadRequest, "file too large to render")
		}

		data, err := io.ReadAll(rd)
		if err != nil {
			log.Println("error reading from file:", err)
			return response.WriteErrorReason(w, http.StatusInternalServerError, "could not read from file")
		}

		doc := mdParser.Parse(data)
		htmlStr := markdown.Render(doc, htmlRenderer)
		return response.WriteData(w, http.StatusOK, htmlStr)

	default:
		return response.WriteStream(w, http.StatusOK, rd)
	}
}

func isBrowser(userAgent string) bool {
	browserPatterns := []string{
		`(?i)chrome`,    // Google Chrome
		`(?i)firefox`,   // Mozilla Firefox
		`(?i)safari`,    // Safari
		`(?i)msie`,      // Internet Explorer
		`(?i)trident`,   // Internet Explorer 11
		`(?i)edge`,      // Microsoft Edge
		`(?i)opera`,     // Opera
		`(?i)opr`,       // Opera
		`(?i)vivaldi`,   // Vivaldi
		`(?i)brave`,     // Brave
		`(?i)ucbrowser`, // UC Browser
	}

	for _, pattern := range browserPatterns {
		matched, _ := regexp.MatchString(pattern, userAgent)
		if matched {
			return true
		}
	}
	return false
}
