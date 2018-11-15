package sync

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/spacemeshos/go-spacemesh/mesh"
	"github.com/spacemeshos/go-spacemesh/p2p"
	"github.com/spacemeshos/go-spacemesh/sync/pb"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type PeersMocks struct {
	p2p.Service
}

func (ml PeersMocks) Count() int {
	return 10
}

func (ml PeersMocks) LatestLayer() int {
	return 10
}

func (ml PeersMocks) GetLayerHash(peer int) string {
	return "asiodfu45987345" //some random string
}

func (ml PeersMocks) GetBlockByID(peer Peer, id string) (Block, error) {
	return nil, nil
}

func (ml PeersMocks) ChoosePeers(pNum int) []Peer {
	return nil
}

func (ml PeersMocks) GetLayerBlockIDs(peers Peer, i int, hash string) ([]string, error) {

	return nil, nil
}

func (ml PeersMocks) GetPeers() []Peer {

	return nil
}

func NewBlockResponseHandler() (func(msg []byte), chan *pb.FetchBlockResp) {
	ch := make(chan *pb.FetchBlockResp)
	foo := func(msg []byte) {
		data := &pb.FetchBlockResp{}
		err := proto.Unmarshal(msg, data)
		if err != nil {
			fmt.Println("some error")
		}
		ch <- data
	}
	return foo, ch
}

func NewLayerHashHandler() (func(msg []byte), chan *pb.LayerHashResp) {
	ch := make(chan *pb.LayerHashResp)
	foo := func(msg []byte) {
		data := &pb.LayerHashResp{}
		err := proto.Unmarshal(msg, data)
		if err != nil {
			fmt.Println("some error")
		}
		ch <- data
	}
	return foo, ch
}

func TestSyncer_Status(t *testing.T) {
	sync := NewSync(nil, nil, nil, Configuration{1, 1, 100 * time.Millisecond, 1})
	assert.True(t, sync.Status() == IDLE, "status was running")
}

func TestSyncer_Start(t *testing.T) {
	layers := mesh.NewLayers(nil, nil)
	sync := NewSync(&PeersMocks{}, layers, nil, Configuration{1, 1, 1 * time.Millisecond, 1})
	fmt.Println(sync.Status())
	sync.Start()
	for i := 0; i < 5 && sync.Status() == IDLE; i++ {
		time.Sleep(1 * time.Second)
	}
	assert.True(t, sync.Status() == RUNNING, "status was idle")
}

func TestSyncer_Close(t *testing.T) {
	sync := NewSync(nil, nil, nil, Configuration{1, 1, 100 * time.Millisecond, 1})
	sync.Start()
	sync.Close()
	s := sync
	_, ok := <-s.forceSync
	assert.True(t, !ok, "channel 'forceSync' still open")
	_, ok = <-s.exit
	assert.True(t, !ok, "channel 'exit' still open")
}

func TestSyncer_ForceSync(t *testing.T) {
	layers := mesh.NewLayers(nil, nil)
	sync := NewSync(&PeersMocks{}, layers, nil, Configuration{1, 1, 60 * time.Minute, 1})
	sync.Start()

	for i := 0; i < 5 && sync.Status() == RUNNING; i++ {
		time.Sleep(1 * time.Second)
	}

	layers.SetLatestKnownLayer(200)
	sync.ForceSync()
	time.Sleep(5 * time.Second)
	assert.True(t, sync.Status() == RUNNING, "status was idle")
}
