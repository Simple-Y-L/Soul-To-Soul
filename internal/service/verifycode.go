package service

import (
	"context"
	pb "kratosItems/api/verifyCode"
	"math/rand"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
		Code: RandCode2(int(req.Length), req.Type),
	}, nil
}
func RandCode(l int, t pb.TYPE) string {
	switch t {
	case pb.TYPE_DEFAULT:
		fallthrough
	case pb.TYPE_DIGIT:
		return randCode("0123456789", l)
	case pb.TYPE_LETTER:
		return randCode("abcdefghijklmnopqrstuvwxyz", l)
	case pb.TYPE_MIXED:
		return randCode("0123456789abcdefghijklmnopqrstuvwxyz", l)
	default:

	}
	return ""
}

// 随机的核心方法
func randCode(chars string, l int) string {

	charsLen := len(chars)

	result := make([]byte, l)
	for i := 0; i < l; i++ {
		randIndex := rand.Intn(charsLen)
		result[i] = chars[randIndex]
	}
	return string(result)
}

func RandCode2(l int, t pb.TYPE) string {
	switch t {
	case pb.TYPE_DEFAULT:
		fallthrough
	case pb.TYPE_DIGIT:
		return randCode2("0123456789", l, 4)
	case pb.TYPE_LETTER:
		return randCode2("abcdefghijklmnopqrstuvwxyz", l, 5)
	case pb.TYPE_MIXED:
		return randCode2("0123456789abcdefghijklmnopqrstuvwxyz", l, 6)
	default:

	}
	return ""
}
func randCode2(chars string, l, idxBits int) string {

	//idxBits = len(fmt.Sprintf("%b", len(chars)))

	idxMask := 1<<idxBits - 1

	idxMax := 63 / idxBits

	result := make([]byte, l)

	for i, cache, remain := 0, rand.Int63(), idxMax; i < l; {
		if 0 == remain {
			cache, remain = rand.Int63(), idxMax
		}
		if randIndex := int(cache & int64(idxMask)); randIndex < len(chars) {
			result[i] = chars[randIndex]
			i++
		}

		cache >>= idxBits

		remain--
	}

	return string(result)
}
