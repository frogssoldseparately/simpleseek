package sreader

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func OpenArchive(name string) (*SimpleZipReader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return OpenArchiveFromBytes(name, &buf)
}

func OpenArchiveFromBytes(name string, buf *[]byte) (*SimpleZipReader, error) {
	archive, err := zip.NewReader(bytes.NewReader(*buf), int64(len(*buf)))
	if err != nil {
		return nil, err
	}
	entries := []*zip.File{}
	byName := map[string]*zip.File{}
	byExtension := map[string][]*zip.File{}
	for _, file := range archive.File {
		name := file.Name
		ext := filepath.Ext(name)
		if _, ok := byExtension[ext]; !ok {
			byExtension[ext] = []*zip.File{}
		}
		entries = append(entries, file)
		byName[name] = file
		byExtension[ext] = append(byExtension[ext], file)
	}
	return &SimpleZipReader{name, archive, entries, byName, byExtension}, nil
}

func (r *SimpleZipReader) GetFile(name string) (*zip.File, bool) {
	entry, ok := r.byName[name]
	return entry, ok
}

func (r *SimpleZipReader) GetFiles() []*zip.File {
	return r.entries
}

func (r *SimpleZipReader) GetFirstByExt(ext string) (*zip.File, bool) {
	list, ok := r.byExtension[ext]
	if !ok || len(list) == 0 {
		return nil, false
	}
	return list[0], true
}

func (r *SimpleZipReader) GetFirstByAnyExt(exts []string) (*zip.File, bool) {
	for _, ext := range exts {
		if f, ok := r.GetFirstByExt(ext); ok {
			return f, true
		}
	}
	return nil, false
}

func (r *SimpleZipReader) GetAllByExt(ext string) ([]*zip.File, bool) {
	list, ok := r.byExtension[ext]
	return list, ok
}

func (r *SimpleZipReader) GetAllByAnyExt(exts []string) []*zip.File {
	out := []*zip.File{}
	for _, ext := range exts {
		if list, ok := r.GetAllByExt(ext); ok {
			out = append(out, list...)
		}
	}
	return out
}
