package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

//
const (
	URLInvoke   = "/api/invoke"
	URLRegister = "/api/chaincode"
)

// Service ...
var Service = &Server{
	ccpool: make(map[string]*Chaincode),
}

// Server ...
type Server struct {
	ccpool map[string]*Chaincode
}

// Register 链码注册
func (s *Server) Register(req *ReqChaincode) error {
	cc, err := NewChaincode(req.Addr)
	if err != nil {
		return err
	}
	s.ccpool[req.Name] = cc
	return nil
}

// Invoke ..
func (s *Server) Invoke(req *ReqTransaction) (*peer.Response, error) {
	cc, ok := s.ccpool[req.Chaincode]
	if !ok {
		return nil, errors.Errorf("链码 %s 不存在，请先注册", req.Chaincode)
	}
	response, err := cc.Invoke(req.Func, req.Args)
	if err != nil {
		return nil, errors.WithMessage(err, "调用失败")
	}
	return response, nil
}

// Router ..
func Router(root *gin.Engine) {
	root.POST(URLRegister, registerCC)
	root.POST(URLInvoke, invoke)
}

// RegisterCC .
// @Summary 链码注册
// @Description 注册登记新的链码
// @Produce json
// @Param req body ReqChaincode true "body参数"
// @Success 200 {object} Response "注册成功"
// @Router /api/chaincode [post]
func registerCC(c *gin.Context) {
	req := &ReqChaincode{}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 1, Message: err.Error()})
		return
	}
	fmt.Println(req)
	err = Service.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 1, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 0, Message: "OK"})
}

// @Summary 交易
// @Description 交易调用
// @Produce json
// @Param req body ReqTransaction true "body参数"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/invoke [post]
func invoke(c *gin.Context) {
	req := &ReqTransaction{}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 1, Message: err.Error()})
		return
	}
	fmt.Println(req)
	rsp, err := Service.Invoke(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 1, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Code: int(rsp.Status), Data: rsp.Payload, Message: rsp.Message})
}
