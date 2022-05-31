package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"owl/base/net"
	"owl/biz/proto"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

func MergeForm(r *http.Request) *Meta {
	var meta = NewMeta(r)
	var err error
	// 20 MB max for upload files
	ct := r.Header.Get("Content-Type")
	multipart := strings.HasPrefix(ct, "multipart/form-data")
	if multipart {
		err = r.ParseMultipartForm(40 * 1024 * 1024)
		if err != nil {
			log.Println("MergeForm err:", err)
			return meta
		}
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println("MergeForm err:", err)
			return meta
		}
	}

	//将json参数合到Form里, 云信回调和抄送用的是json body
	if !multipart {
		jsonform := make(map[string]interface{})
		d := json.NewDecoder(r.Body)
		d.UseNumber()
		_ = d.Decode(&jsonform)
		meta.jsonform = jsonform
	}

	return meta
}

func toMeta(r *http.Request) (*Meta, error) {
	ip := net.GetIP(r)

	md := MergeForm(r)
	md.IP = ip

	return md, nil
}

func onsuccess(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(http.StatusInternalServerError)
	var fail *proto.Fail
	switch e := err.(type) {
	case *Error:
		fail = &proto.Fail{Ok: false, ErrorCode: e.Code, Description: e.Description, Data: e.Data}
	default:
		fail = &proto.Fail{Ok: false, ErrorCode: 500, Description: "INTERNAL_ERROR"}
	}
	_ = json.NewEncoder(w).Encode(fail)
}

func toSuccess(result interface{}) *proto.Success {
	/* 构造成功json */
	var out = &proto.Success{Ok: true}
	out.Result = result
	return out
}

func toFail(err error) *proto.Fail {
	var code = 500
	var description = "INTERAL_ERROR"
	var data interface{}

	var out = &proto.Fail{Ok: false}

	switch err.(type) {
	case *Error:
		e := err.(*Error)
		code = e.Code
		description = e.Description
		data = e.Data
	default:
		panic(err)
	}
	log.Println("✖", code, description)
	out.ErrorCode = code
	out.Description = description
	out.Data = data
	return out
}

type Decorator func(*Meta) error

type ErrorHandler struct {
}

func (h *ErrorHandler) Handle(ctx context.Context, err error) {
	log.Println("Service Error:", err)
}

// 当前活动请求数
var activeRequests int32

func ActiveRequests() int {
	return int(activeRequests)
}

func WaitQuiescent(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	const pollInterval = 500 * time.Millisecond
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		if activeRequests <= 0 {
			return nil
		}
		select {
		case <-ticker.C:
		case <-timer.C:
			return context.DeadlineExceeded
		}
	}
}

// authentication: 是否需要登录
func RegisterMethod(router *mux.Router, path string, method Method, decorators ...Decorator) {
	endpoint := func(md *Meta) (response interface{}, err error) {
		atomic.AddInt32(&activeRequests, 1)
		defer func() {
			atomic.AddInt32(&activeRequests, -1)
			md.Log()
		}()
		r, err := method(md)
		if err != nil {
			return toFail(err), nil
		} else {
			return toSuccess(r), nil
		}
	}
	dec := func(r *http.Request) (*Meta, error) {
		md, err := toMeta(r)
		if err != nil {
			if md != nil {
				md.Log()
			}
			return nil, err
		}
		for _, dec := range decorators {
			err := dec(md)
			if err != nil {
				return nil, err
			}
		}
		return md, nil
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		md, err := dec(r)
		if err != nil {
			encodeError(err, w)
			return
		}
		rsp, err := endpoint(md)
		if err != nil {
			encodeError(err, w)
			return
		}
		err = onsuccess(w, rsp)
		if err != nil {
			encodeError(err, w)
			return
		}
	}

	//handler := kithttp.NewServer(endpoint, dec, onsuccess, opts...)
	router.HandleFunc(path, handler)
}

// 半原生回调，用于第三方平台与服务器的表单回调, 不适用body格式的请求(比如XML)
func RegisterMethodRaw(router *mux.Router, path string, f func(*Meta, http.ResponseWriter)) {
	dec := func(w http.ResponseWriter, r *http.Request) {
		md, err := toMeta(r)
		if err != nil {
			log.Println("500 MD_ERROR:", err)
			http.Error(w, "MD_ERROR", 500)
			return
		}
		atomic.AddInt32(&activeRequests, 1)
		defer func() {
			atomic.AddInt32(&activeRequests, -1)
			md.LogRaw()
		}()
		f(md, w)
	}
	router.HandleFunc(path, dec)
}

// 原生回调，用于第三方平台与服务器的回调
func RegisterHandler(router *mux.Router, path string, f func(*http.Request, http.ResponseWriter)) {
	dec := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&activeRequests, 1)
		st := time.Now()
		defer func() {
			atomic.AddInt32(&activeRequests, -1)
			ip := net.GetIP(r)
			log.Println("<=", r.RequestURI, ip, time.Since(st))
		}()
		f(r, w)
	}
	router.HandleFunc(path, dec)
}
