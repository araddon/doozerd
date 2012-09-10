package server

import (
	"github.com/4ad/doozerd/consensus"
	"github.com/4ad/doozerd/store"
	. "github.com/4ad/doozerd/logging"
	"net"
	"syscall"
)

// ListenAndServe listens on l, accepts network connections, and
// handles requests according to the doozer protocol.
func ListenAndServe(l net.Listener, canWrite chan bool, st *store.Store, p consensus.Proposer, rwsk, rosk string) {
	var w bool
	for {
		c, err := l.Accept()
		if err != nil {
			if err == syscall.EINVAL {
				break
			}
			if e, ok := err.(*net.OpError); ok && e.Err == syscall.EINVAL {
				break
			}
			Log(ERROR, err)
			continue
		}

		// has this server become writable?
		select {
		case w = <-canWrite:
			canWrite = nil
		default:
		}

		go serve(c, st, p, w, rwsk, rosk)
	}
}

func serve(nc net.Conn, st *store.Store, p consensus.Proposer, w bool, rwsk, rosk string) {
	c := &conn{
		c:        nc,
		addr:     nc.RemoteAddr().String(),
		st:       st,
		p:        p,
		canWrite: w,
		rwsk:     rwsk,
		rosk:     rosk,
	}
	Logf(INFO, "New Client Conn: %s", c.addr)
	c.grant("") // start as if the client supplied a blank password
	c.serve()
	nc.Close()
}
