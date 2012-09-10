package consensus

import (
	. "github.com/4ad/doozerd/logging"
	"github.com/4ad/doozerd/store"
)

type Proposer interface {
	Propose(v []byte) store.Event
}

// Set an ephemeral node, it will go away once a client disconnects
func SetEph(p Proposer, client, path string, body []byte, rev int64) (e store.Event) {
	// creation of ephemeral node is two step process
	// 1. set   /eph/clientaddr/path
	// 2. set   path 

	Logf(INFO, "SETEPH cl=%s path=%s rev=%d", client, path, rev)
	e.Mut, e.Err = store.EncodeSet("/eph/"+client+path, "", rev)
	if e.Err != nil {
		Logf(ERROR, "ERROR on encode? %s   err=%v", "/eph/"+client+path, e.Err)
		return
	}
	if e = p.Propose([]byte(e.Mut)); e.Err != nil {
		Logf(ERROR, "ERROR on ephemeral? %s   err=%v", "/eph/"+client+path, e.Err)
	}

	// now store the item
	e.Mut, e.Err = store.EncodeSet(path, string(body), rev)
	if e.Err != nil {
		return
	}

	return p.Propose([]byte(e.Mut))
}

func Set(p Proposer, path string, body []byte, rev int64) (e store.Event) {
	e.Mut, e.Err = store.EncodeSet(path, string(body), rev)
	if e.Err != nil {
		return
	}

	return p.Propose([]byte(e.Mut))
}

func Del(p Proposer, path string, rev int64) (e store.Event) {
	e.Mut, e.Err = store.EncodeDel(path, rev)
	if e.Err != nil {
		return
	}

	return p.Propose([]byte(e.Mut))
}
