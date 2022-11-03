package basedb

import (
	"context"
	"testing"

	"github.com/codingeasygo/util/xsql"
)

func TestAnnounce(t *testing.T) {
	announce := &Announce{
		Type:    AnnounceTypeNormal,
		Marked:  10,
		Title:   "abc",
		Content: xsql.M{"abc": 123},
		Status:  AnnounceStatusNormal,
	}
	err := UpsertAnnounce(context.Background(), announce)
	if err != nil {
		t.Error(err)
		return
	}
	err = UpsertAnnounce(context.Background(), announce)
	if err != nil {
		t.Error(err)
		return
	}
	searcher := AnnounceUnifySearcher{}
	searcher.Where.Key = announce.Title
	searcher.Where.Type = AnnounceTypeAll
	searcher.Where.Status = AnnounceStatusAll
	searcher.Where.Marked = []int{10}
	searcher.Return.Content = 0
	err = searcher.Apply(context.Background())
	if err != nil || len(searcher.Query.Announces) < 1 || len(searcher.Query.Announces[0].Content) > 0 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}
	searcher = AnnounceUnifySearcher{}
	searcher.Where.Key = announce.Title
	searcher.Where.Type = AnnounceTypeAll
	searcher.Where.Status = AnnounceStatusAll
	searcher.Where.Marked = []int{10}
	searcher.Return.Content = 1
	err = searcher.Apply(context.Background())
	if err != nil || len(searcher.Query.Announces) < 1 || len(searcher.Query.Announces[0].Content) < 1 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}
	findAnnounce, err := FindAnnounce(context.Background(), announce.TID)
	if err != nil || findAnnounce.TID != announce.TID {
		t.Error(err)
		return
	}
}
