package hotstuff

import (
	"context"

	pb "github.com/theodocius271/hotstuff/proto"
)

type HotStuffService struct {
	hotStuff *BasicHotStuff
	pb.UnimplementedHotStuffServiceServer
}

func (basic *HotStuffService) SendMsg(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.hotStuff.GetMsgEntrance() <- in
	return &pb.Empty{}, nil
}

func (basic *HotStuffService) SendRequest(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.hotStuff.GetMsgEntrance() <- in
	return &pb.Empty{}, nil
}

func (basic *HotStuffService) SendReply(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (basic *HotStuffService) GetImpl() *BasicHotStuff {
	return basic.hotStuff
}

func NewHotStuffService(impl *BasicHotStuff) *HotStuffService {
	return &HotStuffService{hotStuff: impl}
}
