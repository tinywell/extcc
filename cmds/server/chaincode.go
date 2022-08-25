package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"google.golang.org/grpc/keepalive"
)

// Channel .
const Channel = "mychannel"

// Chaincode ..
type Chaincode struct {
	addr        string
	stream      peer.Chaincode_ConnectClient
	handler     *Handler
	events      map[string]chan *invokeRsp
	db          StateDB
	keepaliveCh <-chan time.Time
}

type invokeRsp struct {
	rsp *peer.Response
	err error
}

// NewChaincode 生成新的 Chaincode 服务实例
func NewChaincode(addr string) (*Chaincode, error) {
	cc := &Chaincode{
		addr:   addr,
		events: make(map[string]chan *invokeRsp),
		db:     NewDB(),
	}
	err := cc.register()
	if err != nil {
		return nil, err
	}

	go cc.serve()
	return cc, nil
}

func (cc *Chaincode) register() error {
	cli, err := cc.getClient()
	if err != nil {
		return errors.WithMessage(err, "")
	}
	stream, err := cli.Connect(context.Background())
	if err != nil {
		return errors.WithMessage(err, "建立到 chaincode 的连接出错")
	}
	cc.stream = stream
	h := NewHandler(stream, cc, cc)
	cc.handler = h
	return nil
}

func (cc *Chaincode) serve() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	cc.keepaliveCh = ticker.C

	type recvMsg struct {
		msg *peer.ChaincodeMessage
		err error
	}
	msgAvail := make(chan *recvMsg, 1)
	receiveMessage := func() {
		in, err := cc.stream.Recv()
		msgAvail <- &recvMsg{in, err}
	}
	go receiveMessage()
	for {
		select {
		case rmsg := <-msgAvail:
			switch {
			case rmsg.err == io.EOF:
				fmt.Println(rmsg.err)
				return
			case rmsg.err == nil:
				fmt.Println(rmsg.msg)
				err := cc.handler.HandleMessage(rmsg.msg)
				if err != nil {
					cc.SendError(rmsg.msg.Txid, err.Error())
				}
				go receiveMessage()
			}

		case <-cc.keepaliveCh:
			err := cc.handler.Keepalive()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (cc *Chaincode) getClient() (peer.ChaincodeClient, error) {
	conn, err := NewClientConn(cc.addr, nil, keepalive.ClientParameters{})
	if err != nil {
		return nil, errors.WithMessage(err, "生成 conn 出错")
	}
	client, err := NewChaincodeClient(conn)
	if err != nil {
		return nil, errors.WithMessage(err, "生成 registerclient 出错")
	}
	return client, nil
}

// Invoke ...
func (cc *Chaincode) Invoke(fn string, args []string, transients map[string][]byte, mspid, cert string) (*peer.Response, error) {
	rc := make(chan *invokeRsp)
	txid := GenerateUUID()
	cc.events[txid] = rc
	var opts []Opt
	if len(transients) > 0 {
		opts = append(opts, WithTransient(transients))
	}
	if len(mspid) > 0 {
		opts = append(opts, WithSigner(mspid, "", cert))
	}
	cc.handler.Invoke(txid, Channel, fn, args, opts...)
	res := <-rc
	if res.err != nil {
		return nil, res.err
	}
	return res.rsp, nil
}

// SendCompleted ..
func (cc *Chaincode) SendCompleted(txid string, res *peer.Response) {
	if rc, ok := cc.events[txid]; ok {
		rc <- &invokeRsp{rsp: res}
		close(rc)
	}
}

// SendError ..
func (cc *Chaincode) SendError(txid string, errMsg string) {
	if rc, ok := cc.events[txid]; ok {
		rc <- &invokeRsp{rsp: nil, err: errors.New(errMsg)}
		close(rc)
	}
}

// GetDB ..
func (cc *Chaincode) GetDB(channel string) StateDB {
	return cc.db
}
