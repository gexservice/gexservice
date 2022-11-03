package baseapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/web"
	"github.com/codingeasygo/web/httptest"
	"github.com/gexservice/gexservice/base/define"
)

func TestFile(t *testing.T) {
	ts := httptest.NewMuxServer()
	ts.Mux.Handle("/", NewUploadH("/tmp", "/upload"))
	res, err := ts.UploadMap(nil, "file", "file.go", "/")
	if err != nil || res.IntDef(-1, "code") != define.Success {
		t.Errorf("%v,%v", err, converter.JSON(res))
		return
	}
	fmt.Printf("Upload--->%v\n", converter.JSON(res))
	res, err = ts.UploadMap(nil, "filex", "file.go", "/")
	if err != nil || res.IntDef(-1, "code") != define.ServerError {
		t.Errorf("%v,%v", err, converter.JSON(res))
		return
	}

	UploadFileAccess = func(s *web.Session) bool { return false }
	res, err = ts.UploadMap(nil, "file", "file.go", "/")
	if err != nil || res.IntDef(-1, "code") != define.NotAccess {
		t.Errorf("%v,%v", err, converter.JSON(res))
		return
	}
}
