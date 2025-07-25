package proxy

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sammy007/open-ethereum-pool/util"
)

func (s *ProxyServer) handleWSClient(cs *Session) error {
	for {
		var req StratumReq
		err := cs.wsconn.ReadJSON(&req)
		if err != nil {
			log.Printf("Websocket read %v: %v", cs.ip, err)
			return err
		}
		s.setWSDeadline(cs.wsconn)
		err = cs.handleWSMessage(s, &req)
		if err != nil {
			log.Printf("Websocket write %v: %v", cs.ip, err)
			return err
		}
	}
}

func (cs *Session) handleWSMessage(s *ProxyServer, req *StratumReq) error {
	// Handle RPC methods
	switch req.Method {
	case "eth_submitLogin":
		var params []string
		err := json.Unmarshal(req.Params, &params)
		if err != nil {
			log.Println("Malformed websocket request params from", cs.ip)
			return err
		}
		reply, errReply := s.handleLoginRPC(cs, params, req.Worker)
		if errReply != nil {
			return cs.sendWSError(req.Id, errReply)
		}
		return cs.sendWSResult(req.Id, reply)
	case "eth_getWork":
		reply, errReply := s.handleGetWorkRPC(cs)
		if errReply != nil {
			return cs.sendWSError(req.Id, errReply)
		}
		return cs.sendWSResult(req.Id, &reply)
	case "eth_submitWork":
		var params []string
		err := json.Unmarshal(req.Params, &params)
		if err != nil {
			log.Println("Malformed websocket request params from", cs.ip)
			return err
		}
		reply, errReply := s.handleTCPSubmitRPC(cs, req.Worker, params)
		if errReply != nil {
			return cs.sendWSError(req.Id, errReply)
		}
		return cs.sendWSResult(req.Id, &reply)
	case "eth_submitHashrate":
		return cs.sendWSResult(req.Id, true)
	default:
		errReply := s.handleUnknownRPC(cs, req.Method)
		return cs.sendWSError(req.Id, errReply)
	}
}

func (cs *Session) sendWSResult(id json.RawMessage, result interface{}) error {
	cs.Lock()
	defer cs.Unlock()

	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
	return cs.wsconn.WriteJSON(&message)
}

func (cs *Session) pushWSNewJob(result interface{}) error {
	cs.Lock()
	defer cs.Unlock()
	// FIXME: Temporarily add ID for Claymore compliance
	// ID is now essential though
	message := JSONPushMessage{Id: 0, Version: "2.0", Result: result}
	return cs.wsconn.WriteJSON(&message)
}

func (cs *Session) sendWSError(id json.RawMessage, reply *ErrorReply) error {
	cs.Lock()
	defer cs.Unlock()

	message := JSONRpcResp{Id: id, Version: "2.0", Error: reply}
	err := cs.wsconn.WriteJSON(&message)
	if err != nil {
		return err
	}
	return errors.New(reply.Message)
}

func (s *ProxyServer) setWSDeadline(wsconn *websocket.Conn) {
	wsconn.SetWriteDeadline(time.Now().Add(s.timeout))
}

func (s *ProxyServer) broadcastWSNewJobs() {
	t := s.currentBlockTemplate()
	if t == nil || len(t.Header) == 0 || s.isSick() {
		return
	}
	height := util.ToHex(int64(t.Height), false)
	reply := []string{t.Header, t.Seed, s.diff, height}

	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	start := time.Now()
	bcast := make(chan int, 1024)
	n := 0
	count := 0

	for m, _ := range s.sessions {
		if m.wsconn == nil {
			continue
		}

		n++
		count++
		bcast <- n

		go func(cs *Session) {
			err := cs.pushWSNewJob(&reply)
			<-bcast
			if err != nil {
				log.Printf("Job transmit error to %v@%v: %v", cs.login, cs.ip, err)
				s.removeSession(cs)
			} else {
				s.setWSDeadline(cs.wsconn)
			}
		}(m)
	}

	log.Printf("Jobs broadcast finished to %v websocket miners: %s", count, time.Since(start))
}
