package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cagnosolutions/appstats/pkg/appstats"
)

func init() {
	// register and run
	appstats.Serve(":8080")
}

func main() {
	ds := NewDS()
	mux := http.NewServeMux()
	mux.Handle("/data/put", handlePut(ds))
	mux.Handle("/data/get/", handleGetID(ds))
	mux.Handle("/data/get/all", handleGetAll(ds))
	mux.Handle("/data/del/", handleDelID(ds))
	mux.Handle("/data/add/a/ton", handleAddATon(ds))
	// appstats.Register(mux)
	log.Println(http.ListenAndServe(":6060", mux))
}

const (
	pageSize  = 8192
	pageCount = 64
)

type dataPage struct {
	data [pageSize]byte
	used int
}

type DS struct {
	pdata map[int]dataPage
	psize int
}

func NewDS() *DS {
	return &DS{
		pdata: make(map[int]dataPage, pageCount),
		psize: 0,
	}
}

func (ds *DS) Put(id int, rec string) error {
	// check to see if we are adding a new entry,
	// or updating an existing one
	dp, ok := ds.pdata[id]
	// if the id does not exist...
	if !ok {
		// make sure that the id provided is the correct
		// one that should come next in the sequence
		if id != ds.psize+1 {
			// if it is not a correct id, return an error
			errorf := fmt.Sprintf(
				"id is incorrect, next id would be: %d",
				ds.psize+1,
			)
			return errors.New(errorf)
		}
		// otherwise, create the new record
		dp := dataPage{}
		copy(dp.data[dp.used:], rec+":")
		dp.used = len(rec) + 1
		// and add the new data page to the set
		ds.pdata[id] = dp
		// and increment the size and return
		ds.psize++
		return nil
	}
	// and if we are here, this means an existing record
	// with the provided id was found, so we are doing an
	// update (but only if there is room) so check that
	if len(dp.data)-dp.used < len(rec)+1 {
		// we do not have room, return and error
		errorf := fmt.Sprintf(
			"no more room, need: %d, have: %d",
			len(rec)+1, len(dp.data)-dp.used,
		)
		return errors.New(errorf)
	}
	// if we get here, we're doing an update, and we have
	// checked and made sure that we have room to insert
	copy(dp.data[dp.used:], rec+":")
	return nil
}

func (ds *DS) Get(id int) (string, error) {
	// check to see if record can be found
	dp, ok := ds.pdata[id]
	// if the id does not exist...
	if !ok {
		// return an error
		errorf := fmt.Sprintf("record not found (using id=%d)", id)
		return "", errors.New(errorf)
	}
	// otherwise, record was found so return it
	return string(dp.data[:dp.used]), nil
}

func (ds *DS) GetAll() string {
	// check for records
	if len(ds.pdata) < 1 {
		return "no records present"
	}
	// init variable to store data
	var ss string
	// range the records
	for id, dp := range ds.pdata {
		ss += fmt.Sprintf(
			"record.id=%d\nrecord.data=%q\n\n", id,
			dp.data[:dp.used],
		)
	}
	// return records string
	return ss
}

func (ds *DS) Del(id int) {
	// check to see if record can be found
	_, ok := ds.pdata[id]
	// if the id does not exist...
	if !ok {
		// do nothing because nothing needs to be removed
		// because it cannot be found, so it's not there
		return
	}
	// otherwise, record was found so wipe the page data
	ds.pdata[id] = dataPage{
		data: [pageSize]byte{},
		used: 0,
	}
}

func handlePut(ds *DS) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		v, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, http.StatusText(code), code)
			return
		}
		sid := v.Get("id")
		if sid == "" {
			code := http.StatusExpectationFailed
			http.Error(w, "'id' was not found or was empty", code)
			return
		}
		id, err := strconv.Atoi(sid)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, "error converting id string to integer", code)
			return
		}
		rec := v.Get("rec")
		if rec == "" {
			code := http.StatusExpectationFailed
			http.Error(w, "'rec' was not found or was empty", code)
			return
		}
		err = ds.Put(id, rec)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, err.Error(), code)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func handleGetID(ds *DS) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		v, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, http.StatusText(code), code)
			return
		}
		sid := v.Get("id")
		if sid == "" {
			code := http.StatusExpectationFailed
			http.Error(w, "'id' was not found or was empty", code)
			return
		}
		id, err := strconv.Atoi(sid)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, "error converting id string to integer", code)
			return
		}
		rec, err := ds.Get(id)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, err.Error(), code)
			return
		}
		_, err = fmt.Fprintf(
			w, "page_data: %q\n\nhext_dump: %s", rec,
			hex.Dump([]byte(rec)),
		)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, "error writing data", code)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func handleGetAll(ds *DS) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		recs := ds.GetAll()
		_, err := fmt.Fprintf(
			w, "%s", recs,
		)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, "error writing data", code)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func handleAddATon(ds *DS) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		nextid := ds.psize + 1
		for i := 0; i < 64; i++ {
			rec := fmt.Sprintf(
				"adding more data, this is record: %d",
				nextid,
			)
			ds.Put(nextid, rec)
			nextid++
		}
	}
	return http.HandlerFunc(fn)
}

func handleDelID(ds *DS) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		v, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, http.StatusText(code), code)
			return
		}
		sid := v.Get("id")
		if sid == "" {
			code := http.StatusExpectationFailed
			http.Error(w, "'id' was not found or was empty", code)
			return
		}
		id, err := strconv.Atoi(sid)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, "error converting id string to integer", code)
			return
		}
		ds.Del(id)
		_, err = fmt.Fprintf(
			w, "successfully recycled page (with id: %d):", id,
		)
		if err != nil {
			code := http.StatusExpectationFailed
			http.Error(w, "error writing data", code)
			return
		}
	}
	return http.HandlerFunc(fn)
}
