package baseapi

import (
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
)

var SrvAddr = func() string {
	panic("SrvAddr is not initial")
}

func Handle(pre string, mux *web.SessionMux) {
	//config
	mux.HandleFunc("^"+pre+"/usr/updateSysConfig(\\?.*)?$", UpdateSysConfigH)
	mux.HandleFunc("^"+pre+"/usr/loadSysConfig(\\?.*)?$", LoadSysConfigH)
	//version object
	mux.HandleFunc("^"+pre+"/usr/upsertVersionObject(\\?.*)?$", UpsertVersionObjectH)
	mux.HandleFunc("^"+pre+"/pub/findVersionObject(\\?.*)?$", FindVersionObjectH)
	mux.HandleFunc("^"+pre+"/pub/searchVersionObject(\\?.*)?$", SearchVersionObjectH)
	mux.HandleFunc("^"+pre+"/pub/vobject/.*$", LoadLatestVersionObjectH)
	//announce
	mux.HandleFunc("^"+pre+"/usr/upsertAnnounce(\\?.*)?$", UpsertAnnounceH)
	mux.HandleFunc("^"+pre+"/pub/loadAnnounce(\\?.*)?$", LoadAnnounceH)
	mux.HandleFunc("^"+pre+"/pub/searchAnnounce(\\?.*)?$", SearchAnnounceH)
}

var EditSysConfigAccess = func(s *web.Session) bool {
	return true
}

var EditVersionObjectAccess = func(s *web.Session) bool {
	return true
}

var EditAnnounceAccess = func(s *web.Session) bool {
	return true
}

var UploadFileAccess = func(s *web.Session) bool {
	return true
}

func RecvValidJSON(s *web.Session, valider basedb.Validable) (err error) {
	_, err = s.RecvJSON(interface{}(valider))
	if err == nil {
		err = valider.Valid()
	}
	return
}
