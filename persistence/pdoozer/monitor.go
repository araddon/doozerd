package main

import (
	"fmt"
	"github.com/ha/doozer"
	st "github.com/ha/doozerd/store"
	"regexp"
)

var exclude = regexp.MustCompile("^/ctl")

func encode(ev doozer.Event) (mut string) {
	switch {
	case ev.IsSet():
		mut = st.MustEncodeSet(ev.Path, string(ev.Body), ev.Rev)
	case ev.IsDel():
		mut = st.MustEncodeDel(ev.Path, ev.Rev)
	}
	return
}

func monitor() {
	var (
		ev  doozer.Event
		err error
		rev int64 = -1
	)
	for {
		ev, err = conn.Wait("/**", rev)
		if err != nil {
			exit(err)
		}
		rev = ev.Rev + 1

		if exclude.MatchString(ev.Path) {
			continue
		}

		mut := encode(ev)
		fmt.Println(mut) // TODO(aram): check -v flag.

		store <- mutation{ev: ev, mut: mut}
	}
}
