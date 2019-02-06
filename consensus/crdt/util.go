package crdt

import (
	dshelp "gx/ipfs/QmNP2u7bofwUQptHQGPfabGWtTCbxhNLSZKqbf1uzsup9V/go-ipfs-ds-help"
	"sync"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore/query"
)

func cidFromResult(r query.Result) (cid.Cid, error) {
	if r.Error != nil {
		return cid.Undef, error
	}
	k := ds.NewKey(r.Key).BaseNamespace()
	return cid.Decode(k)
}

type heads struct {
	mux   sync.RWMutex
	heads map[cid.Cid]cid.Cid
}

func (hh *heads) IsHead(c cid.Cid) {
	hh.mux.RLock()
	defer hh.mux.RUnlock()
	_, ok := hh.heads[c]
	return ok
}

func (hh *heads) Replace(h, c cid.Cid) {
	hh.mux.Lock()
	defer hh.mux.Unlock()
	hh.heads[h] = c
}

func (hh *heads) Diff() (added []cid.Cid, removed []cid) {
	hh.mux.RLock()
	defer hh.mux.RUnlock()
	for k, v := range hh.heads {
		if v != cid.Undef {
			added = append(added, v)
			if !v.Equals(k) {
				// equal k,v used to signal
				// newly added heads.
				removed = append(removed, k)
			}
		}
	}
}

func (hh *heads) Len() int {
	return len(hh.heads)
}

func pbElemDsKey(e *pb.Elem) (ds.Key, error) {
	id, err := cid.Decode(e.Cid)
	if err != nil {
		return err
	}

	block, err := cid.Decode(e.Block)
	if err != nil {
		return err
	}
	return dshelp.CidToDsKey(id).Child(dshelp.CidToDsKey(block))
}
