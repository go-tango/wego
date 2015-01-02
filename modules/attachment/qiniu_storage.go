package attachment

import (
	"fmt"
	"github.com/go-tango/wetalk/modules/models"
	"github.com/go-tango/wetalk/modules/utils"
	"github.com/qiniu/api/io"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	goio "io"
	"path/filepath"
	"time"
)

func SaveImageToQiniu(m *models.Image, r goio.ReadSeeker, mime string, filename string, created time.Time, bucketName string) error {
	var ext string

	// check image mime type
	switch mime {
	case "image/jpeg":
		ext = ".jpg"

	case "image/png":
		ext = ".png"

	case "image/gif":
		ext = ".gif"

	default:
		ext = filepath.Ext(filename)
		switch ext {
		case ".jpg", ".png", ".gif":
		default:
			return fmt.Errorf("unsupport image format `%s`", filename)
		}
	}

	// decode image
	var img image.Image
	var err error
	switch ext {
	case ".jpg":
		m.Ext = 1
		img, err = jpeg.Decode(r)
	case ".png":
		m.Ext = 2
		img, err = png.Decode(r)
	case ".gif":
		m.Ext = 3
		img, err = gif.Decode(r)
	}

	if err != nil {
		return err
	}

	m.Width = img.Bounds().Dx()
	m.Height = img.Bounds().Dy()
	m.Created = created

	//save to database
	if err := m.Insert(); err != nil || m.Id <= 0 {
		return err
	}

	m.Token = m.GetToken()
	if err := m.Update(); err != nil {
		return err
	}

	var key = m.Token

	//reset reader pointer
	if _, err := r.Seek(0, 0); err != nil {
		return err
	}

	//save to qiniu
	var uptoken = utils.GetQiniuUptoken(bucketName)
	var putRet io.PutRet
	var putExtra = &io.PutExtra{
		MimeType: mime,
	}

	err = io.Put(nil, &putRet, uptoken, key, r, putExtra)
	if err != nil {
		return err
	}

	return nil
}
