package baseapi

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/codingeasygo/util/uuid"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
)

type UploadH struct {
	SaveDIR string
	Prefix  string
}

func NewUploadH(saveDIR, prefix string) (upload *UploadH) {
	upload = &UploadH{
		SaveDIR: saveDIR,
		Prefix:  prefix,
	}
	return
}

//UploadH is http handler
/**
 *
 * @api {POST} /usr/upload Upload
 * @apiName Upload
 * @apiGroup File
 *
 *
 * @apiParam  {String} file the multipart key
 *
 * @apiSuccess (200) {Number} code the response code, see common define
 * @apiSuccess (200) {String} path the file access path
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "path": "/upload/2021-11-04/6183a586285c6681b0000002.go"
 * }
 */
func (u *UploadH) SrvHTTP(s *web.Session) web.Result {
	if !UploadFileAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess,
		})
	}
	savePath := fmt.Sprintf("%v/%v", time.Now().Format("2006-01-02"), uuid.New())
	_, err := s.RecvMultipart(false, false, func(p *multipart.Part) (filename string, mode os.FileMode, external []io.Writer, err error) {
		if p.FormName() != "file" {
			err = fmt.Errorf("not file")
			return
		}
		name := p.FileName()
		savePath += filepath.Ext(name)
		filename = filepath.Join(u.SaveDIR, savePath)
		mode = os.ModePerm
		return
	})
	if err != nil {
		return s.SendJSON(xmap.M{
			"code":    define.ServerError,
			"message": err.Error(),
		})
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
		"path": filepath.Join(u.Prefix, savePath),
	})
}
