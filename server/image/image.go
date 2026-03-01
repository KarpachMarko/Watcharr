package image

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/buckket/go-blurhash"
	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/database/entity"
	"gorm.io/gorm"
)

// TODO now that this file is in the image package it no longer needs to have "image(s)"
// in the name of all the functions..

// Insert an image into database
func InsertImage(db *gorm.DB, hash string, path string, f io.Reader) (entity.Image, error) {
	bh, _ := GetBlurHash(f)
	img := entity.Image{
		Hash:     hash,
		Path:     path,
		BlurHash: bh,
	}
	r := db.Where(entity.Image{Hash: hash}).FirstOrCreate(&img)
	if r.Error != nil {
		slog.Error("insertImage firstOrCreate failed!", "error", r.Error)
		return entity.Image{}, errors.New("failed to select or create image")
	}
	return img, nil
}

func GetBlurHash(img io.Reader) (string, error) {
	i, _, err := image.Decode(img)
	if err != nil {
		// Handle errors
		slog.Error("getBlurHash decoding image failed", "error", err)
		return "", errors.New("decoding image failed")
	}
	bh, err := blurhash.Encode(6, 5, i)
	if err != nil {
		// Handle errors
		slog.Error("getBlurHash generation failed", "error", err)
		return "", errors.New("blur hash generation failed")
	}
	slog.Debug("getBlurHash", "hash", bh)
	return bh, nil
}

func CleanupImages(db *gorm.DB) {
	slog.Info("cleanupImages running")
	var unusedImgs []entity.Image
	// Select images that are not referenced by at least one other row.
	// Currently only used for user avatars, add new tables when used.
	db.Raw(`SELECT *
FROM images
WHERE NOT EXISTS (
	SELECT 1
	FROM users
	WHERE users.avatar_id = images.id
);`).Scan(&unusedImgs)
	slog.Info("cleanupImages: scanned for unused images", "amount", len(unusedImgs))
	if len(unusedImgs) > 0 {
		for _, v := range unusedImgs {
			slog.Debug("cleanupImages: removing an image", "id", v.ID, "path", v.Path)
			err := db.Transaction(func(tx *gorm.DB) error {
				// Try to delete image from db
				if err := tx.Where("id = ?", v.ID).Delete(&entity.Image{}).Error; err != nil {
					return err
				}
				// hope its ok to do this sorta thing here :skull:
				if err := os.Remove(path.Join(config.DataPath, v.Path)); err != nil {
					return err
				}
				// commit transaction if no errors
				return nil
			})
			if err != nil {
				slog.Error("cleanupImages: failed to remove image - db row and file kept", "img", v, "error", err)
			} else {
				slog.Debug("cleanupImages: successfully removed unused image.", "id", v.ID)
			}
		}
	}
}

func IsValidImageType(f multipart.File) error {
	// Read first 512 bytes, since that is all `DetectContentType` will evaluate on.
	// Reading whole file is a waste.
	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		slog.Error("isValidImageType: failed to read file into buffer", "error", err)
		return errors.New("failed to verify if image is valid")
	}
	t := http.DetectContentType(buff)
	slog.Debug("isValidImageType", "type", t)
	if t != "image/png" && t != "image/jpeg" && t != "image/webp" && t != "image/gif" {
		slog.Debug("isValidImageType: rejecting file as not valid (supported) image type")
		return errors.New("invalid file type")
	}
	return nil
}

func DownloadAndInsertImage(db *gorm.DB, url string, imgSubPath string) (entity.Image, error) {
	slog.Debug("Attempting to download image", "url", url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return entity.Image{}, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return entity.Image{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read body into byte array, then create new reader
	// So we have the ability to seek.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("downloadAndInsertImage failed to read response into byte array", "error", err)
		return entity.Image{}, err
	}
	br := bytes.NewReader(b)

	h := sha256.New()
	if _, err := io.Copy(h, br); err != nil {
		slog.Error("DownloadAndInsertImage: Copy failed!", "error", err)
		return entity.Image{}, errors.New("copy failed")
	}
	hs := hex.EncodeToString(h.Sum(nil))

	// Seek back for file
	_, err = br.Seek(0, 0)
	if err != nil {
		slog.Error("downloadAndInsertImage seeking back to start of br failed", "error", err)
		return entity.Image{}, err
	}

	outp := path.Join("img/", imgSubPath, hs[0:1], hs+filepath.Ext(resp.Request.URL.Path))
	dataOutP := path.Join(config.DataPath, outp)

	// Create the file
	out, err := os.Create(dataOutP)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(dataOutP), 0764)
			if err != nil {
				return entity.Image{}, err
			}
			// If dirs made, try making file again
			out, err = os.Create(dataOutP)
			if err != nil {
				return entity.Image{}, err
			}
		} else {
			return entity.Image{}, err
		}
	}
	defer out.Close()
	_, err = io.Copy(out, br)
	if err != nil {
		return entity.Image{}, err
	}

	// Seek back for insertImage
	_, err = br.Seek(0, 0)
	if err != nil {
		slog.Error("downloadAndInsertImage seeking back to start of br failed", "error", err)
		return entity.Image{}, err
	}

	img, err := InsertImage(db, hs, outp, br)
	if err != nil {
		return entity.Image{}, err
	}

	return img, nil
}
