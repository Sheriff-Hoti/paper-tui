package util

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sheriff-Hoti/paper-tui/cache"
)

func GetWallpapers(dir string) ([]string, error) {

	allowedImageExtensions := map[string]struct{}{
		".jpg":  {},
		".png":  {},
		".jpeg": {},
		".gif":  {},
	}

	// Expand environment variables like $HOME
	dirEnvExpanded := os.ExpandEnv(dir)

	// Get absolute path of directory
	absDir, err := filepath.Abs(dirEnvExpanded)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))

			if _, ok := allowedImageExtensions[ext]; ok {
				fullPath := filepath.Join(absDir, entry.Name())
				absPath, err := filepath.Abs(fullPath)
				if err != nil {
					return nil, err
				}
				fileNames = append(fileNames, absPath)
			}
		}

	}

	return fileNames, nil
}

func ToPng(absInputPath string) (string, error) {
	if absInputPath == "" {
		return "", errors.New("input path is empty")
	}
	if !filepath.IsAbs(absInputPath) {
		return "", fmt.Errorf("input path must be absolute: %q", absInputPath)
	}

	ext := strings.ToLower(filepath.Ext(absInputPath))
	switch ext {
	case ".jpg", ".jpeg", ".gif":
		// convert
	case ".png":
		return absInputPath, nil
	default:
		return "", fmt.Errorf("unsupported extension %q (supported: .jpg/.jpeg/.png/.gif)", ext)
	}

	in, err := os.Open(absInputPath)
	if err != nil {
		return "", fmt.Errorf("open input: %w", err)
	}
	defer in.Close()

	img, format, err := image.Decode(in)
	if err != nil {
		return "", fmt.Errorf("decode image (ext=%s, format=%q): %w", ext, format, err)
	}
	if img == nil {
		return "", fmt.Errorf("decode returned nil image (ext=%s, format=%q)", ext, format)
	}

	outDir, err := cache.GetCacheDir()
	if err != nil {
		return "", fmt.Errorf("get cache dir: %w", err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir cache dir: %w", err)
	}

	base := filepath.Base(absInputPath)
	nameNoExt := strings.TrimSuffix(base, filepath.Ext(base))
	outPath := filepath.Join(outDir, nameNoExt+".png")

	out, err := os.Create(outPath) // truncates if exists
	if err != nil {
		return "", fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", fmt.Errorf("encode png: %w", err)
	}

	return outPath, nil
}
