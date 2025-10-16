package pkg

import (
	pb "kratosItems/api/user"

	"google.golang.org/protobuf/types/known/structpb"
)

func RequestReturnDynamic(code int, msg string, data interface{}) (*pb.PleaseReturn, error) {
	var structData *structpb.Struct

	if data != nil {
		var err error
		if mapData, ok := data.(map[string]interface{}); ok {
			structData, err = structpb.NewStruct(mapData)
			if err != nil {
				return nil, err
			}
		}
	}
	return &pb.PleaseReturn{
		Code: int64(code),
		Msg:  msg,
		Data: structData,
	}, nil
}
