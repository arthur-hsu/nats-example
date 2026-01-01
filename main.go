package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/nats-io/nats.go" // 專門處理 proto 與 json 的轉換
	pb "google.golang.org/protobuf/proto"

	// 專門處理 proto 與 json 的轉換
	// 這裡請改成你本地的 proto package 路徑
	"nats-example/usp"
	"nats-example/usp/usp_msg"
)

func SetMsg(ObjPath string, Param string, Value string, Required bool) usp_msg.Set {
	return usp_msg.Set{
		AllowPartial: false,
		UpdateObjs: []*usp_msg.Set_UpdateObject{
			{
				ObjPath: ObjPath,
				ParamSettings: []*usp_msg.Set_UpdateParamSetting{
					{
						Param:    Param,
						Value:    Value,
						Required: Required,
					},
				},
			},
		},
	}
}
func getMsg(ObjPath any) usp_msg.Get {
	var ParamPaths []string
	switch v := ObjPath.(type) {
	case string:
		ParamPaths = []string{v}
	case []string:
		ParamPaths = v
	default:
		return usp_msg.Get{}
	}
	return usp_msg.Get{
		ParamPaths: ParamPaths,
	}
}

func main() {
	// --- 1. 連線設定 ---
	natsURL := "nats://100.81.97.99:4222"
	user := "oktopususer"
	pw := "oktopuspw"

	// 根據你監聽到的 ID
	agentID := "oktopus-0-ws"
	controllerID := "oktopusController"

	nc, err := nats.Connect(natsURL,
		nats.UserInfo(user, pw),
		nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// --- 2. 構建 USP Message ---
	method := "Get"

	var uspInstance usp_msg.Msg
	switch method {
	case "Set":
		Content := SetMsg("Device.LocalAgent.MTP.1.", "Alias", "Locus-Protobuf-Final", true)
		uspInstance = usp.NewSetMsg(Content)
	case "Get":
		// Content := getMsg("Device.LocalAgent.MTP.1.")
		Content := getMsg([]string{"Device.LocalAgent.MTP.1.Alias", "Device.LocalAgent.MTP.1.Protocol"})
		uspInstance = usp.NewGetMsg(Content)
	}

	msgBinary, err := pb.Marshal(&uspInstance)
	if err != nil {
		log.Fatalf("USP Message 打包失敗: %v", err)
	}

	uspRecord := usp.NewUspRecord(msgBinary, agentID, controllerID)

	recordBinary, err := pb.Marshal(&uspRecord)
	if err != nil {
		log.Fatalf("USP Record 打包失敗: %v", err)
	}

	subject := fmt.Sprintf("ws-adapter.usp.v1.%s.api", agentID)
	err = nc.Publish(subject, recordBinary)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 官方格式封包已發送至: %s:%s\n", method, subject)
}
