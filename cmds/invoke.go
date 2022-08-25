package cmds

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"extcc/cmds/server"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/spf13/cobra"
)

// InvokeURL ..
const InvokeURL = server.URLInvoke

// InvokeCmd 链码调用
var InvokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "调用链码",
	Run: func(cmd *cobra.Command, args []string) {
		invoke()
	},
}

var (
	address   string
	chaincode string
	fn        string
	args      string
	isInit    bool
	// txid    string
	channel   string
	transient string
	mspid     string
	cert      string
)

func init() {
	InvokeCmd.Flags().StringVar(&address, "addr", "localhost:9999", "仿真服务地址")
	InvokeCmd.Flags().StringVar(&chaincode, "cc", "mycc", "链码名称")
	// InvokeCmd.Flags().StringVar(&fn, "fn", "get", "调用方法")
	InvokeCmd.Flags().StringVar(&args, "args", "[]", "链码调用参数(exp:'[\"get\",\"A\"]')")
	// InvokeCmd.Flags().StringVar(&channel, "channel", "mychannel", "通道")
	InvokeCmd.Flags().BoolVarP(&isInit, "init", "i", false, "是否调用初始化函数")
	InvokeCmd.Flags().StringVarP(&mspid, "mspid", "", "Org1MSP", "交易用户 MSPID")
	InvokeCmd.Flags().StringVarP(&transient, "transient", "t", "{}", "秘密参数")
}

func invoke() {
	strArgs := make([]string, 0)
	err := json.Unmarshal([]byte(args), &strArgs)
	if err != nil {
		panic(err)
	}

	transientArgs := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader([]byte(transient)))
	decoder.UseNumber()
	err = decoder.Decode(&transientArgs)
	if err != nil {
		panic(err)
	}
	transienMap := make(map[string][]byte)
	for k, v := range transientArgs {
		var vv []byte
		switch t := v.(type) {
		case json.Number:
			var data interface{}
			data, err = t.Int64()
			if err != nil {
				data, err = t.Float64()
				if err != nil {
					panic(err)
				}
			}
			bytebuf := bytes.NewBuffer([]byte{})
			binary.Write(bytebuf, binary.BigEndian, data)
			vv = bytebuf.Bytes()
		case string:
			vv = []byte(t)
		default:
			vv = v.([]byte)
		}
		transienMap[k] = vv
	}
	fmt.Println(transienMap)

	if len(strArgs) == 0 {
		panic("合约调用参数不能为空")
	}
	fn = strArgs[0]

	if isInit {
		fn = "Init"
	}
	req := &server.ReqTransaction{
		Chaincode:    chaincode,
		Func:         fn,
		Args:         strArgs[1:],
		MSPID:        mspid,
		TransientMap: transienMap,
	}
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("http://%s%s", address, InvokeURL)
	rsp, err := Post(url, body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rsp))
}

func printRsp(rsp *peer.ChaincodeMessage) {
	if rsp.Type != peer.ChaincodeMessage_COMPLETED {
		fmt.Println("错误的结果数据")
		return
	}
	response := &peer.Response{}
	err := proto.Unmarshal(rsp.Payload, response)
	if err != nil {
		fmt.Printf("解析结果数据出错: %s\n", err.Error())
	}
	fmt.Printf("Response:\n\tStatus:%d\n\tMessage:%s\n\tPayload:%s\n",
		response.Status, response.Message, string(response.Payload))
}
