package server

import (
	"bufio"
	"elide/internal/api/request"
	"elide/internal/api/response"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/anonyindian/gotgproto/ext"
)

type Server struct {
	protocol int
	port     int
	handler  map[string]HandlerFunc
	Context  *ext.Context
}

func New(protocol, port int) *Server {
	return &Server{
		protocol: protocol,
		port:     port,
		handler:  make(map[string]HandlerFunc),
	}
}

func (s *Server) AddHandler(method string, handler HandlerFunc) {
	s.handler[method] = handler
}

func (s *Server) handlerWrapper(w io.Writer, b []byte) bool {
	req, err := request.Parse(b)
	if err != nil {
		w.Write(response.InitError(err))
		return false
	}
	rHandler, ok := s.handler[req.Method]
	if !ok {
		w.Write(response.Error("unknown method: " + req.Method))
		return false
	}
	res, err := rHandler(s.Context, req.Data)
	if err != nil {
		w.Write(response.InitError(err))
		return false
	}
	w.Write(response.Result(res))
	return true
}

func parseFormData(url *url.URL) map[string]any {
	sRaw := strings.Split(url.RawQuery, "&")
	paramKeys := make([]string, len(sRaw))
	for i, s := range sRaw {
		paramKeys[i] = strings.Split(s, "=")[0]
	}
	v := url.Query()
	params := make(map[string]any)
	for _, key := range paramKeys {
		val := v.Get(key)
		nVal, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			params[key] = nVal
			continue
		}
		bVal, err := strconv.ParseBool(val)
		if err == nil {
			params[key] = bVal
			continue
		}
		if val[0] == '[' && val[len(val)-1] == ']' {
			if len(val) == 2 {
				continue
			}
			arrElems := strings.Split(val[1:(len(val)-1)], ",")
			var arr = make([]any, len(arrElems))
			for i, arrValue := range arrElems {
				arrValue = strings.TrimSpace(arrValue)
				nArrValue, err := strconv.ParseInt(arrValue, 10, 64)
				if err == nil {
					arr[i] = nArrValue
				} else {
					arr[i] = arrValue
				}
			}
			params[key] = arr
			continue
		}
		params[key] = val
	}
	return params
}

func (s *Server) send(conn net.Conn, buf []byte) {
	conn.Write(buf)
	s.transmitEnd(conn)
}

func (s *Server) transmitEnd(conn net.Conn) {
	conn.Write([]byte{0})
}

func (s *Server) Start() {
	// go e.establishTelegramConn()
	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	switch s.protocol {
	case 0:
		// HTTP
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				method := r.URL.Path[1:]
				rHandler, ok := s.handler[method]
				if ok {
					params := parseFormData(r.URL)
					data, err := json.Marshal(params)
					if err != nil {
						w.Write(response.InitError(err))
						return
					}
					res, err := rHandler(s.Context, data)
					if err != nil {
						w.Write(response.InitError(err))
						return
					}
					w.Write(response.Result(res))
					return
				}
			}
			body := r.Body
			b, err := io.ReadAll(body)
			if err != nil {
				w.Write(response.InitError(err))
				return
			}
			if len(b) == 0 {
				method, data := r.FormValue("method"), r.FormValue("data")
				if method == "" || data == "" {
					w.Write(response.Error("insufficient info provided"))
					return
				}
				b = []byte(fmt.Sprintf(`{"method":"%s","data":%s}`, method, data))
			}
			_ = s.handlerWrapper(w, b)
		})
		server := http.Server{
			Addr:    addr,
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	case 1:
		// TCP
		l, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("failed to listen on", "tcp", ":", err.Error())
			os.Exit(1)
		}
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			// Handle connections in a new goroutine.
			go func(conn net.Conn) {
				for {
					b, err := bufio.NewReader(conn).ReadBytes(0)
					if err != nil {
						s.send(conn, response.InitError(err))
						return
					}
					_ = s.handlerWrapper(conn, b[:len(b)-1])
					s.transmitEnd(conn)
				}
				// err = conn.Close()
				// if err != nil {
				// 	fmt.Println("failed to close connection:", err.Error())
				// }
			}(conn)
		}
	case 2:
		fmt.Println("UDP is not implemented yet, please use HTTP or TCP instead!")
		os.Exit(1)
	}
}
