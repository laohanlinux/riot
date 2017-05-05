package clientrpc

import (
	"fmt"
	"testing"
	"time"

	"github.com/laohanlinux/assert"
	"github.com/laohanlinux/riot/api"
)

func TestClientRPCKV(t *testing.T) {
	var (
		addrs      = []string{"127.0.0.1:32124", "127.0.0.1:32125", "127.0.0.1:32123"}
		arg        api.NotArg
		reply      api.NodeStateReply
		bucketName = "student"
		key        = "lusi"
		value      = []byte("{\"Age\":100}")
	)
	fmt.Println(InitRPC(addrs, 10))

	time.Sleep(time.Second)
	assert.Nil(t, DefaultLeaderRPC.Call(APIServiceState, &arg, &reply))
	fmt.Println(reply)
	{
		var (
			arg   api.NotArg
			reply api.PeersReply
		)
		assert.Nil(t, DefaultRaftRPC.Call(APIServicePeers, &arg, &reply))
		fmt.Println(reply)
	}

	{
		var (
			arg   api.NotArg
			reply api.LeaderReply
		)
		assert.Nil(t, DefaultRaftRPC.Call(APIServiceLeader, &arg, &reply))
		fmt.Println(reply)
	}

	{
		//var (
		//	arg   api.NotArg
		//	reply api.SnapshotReply
		//)
		//assert.Nil(t, DefaultRaftRPC.Call(APIServiceSnapshot, &arg, &reply))
		//fmt.Println(reply)
	}

	{
		var (
			arg   = api.CreateBucketArg{BucketName: bucketName}
			reply api.NotReply
		)
		assert.Nil(t, DefaultLeaderRPC.Call(APIServiceCreateBucket, &arg, &reply))

	}
	{
		var (
			arg   = api.BucketInfoArg{BucketName: bucketName}
			reply api.BucketInfoReply
		)
		assert.Nil(t, DefaultLeaderRPC.Call(APIServiceBucketInfo, &arg, &reply))
		assert.Equal(t, true, reply.Has)
	}

	{
		var (
			arg   = api.DelBucketArg{BucketName: bucketName}
			reply api.NotReply
		)
		assert.Nil(t, DefaultLeaderRPC.Call(APIServiceDelBucket, &arg, &reply))
	}

	{
		var (
			arg   = api.BucketInfoArg{BucketName: bucketName}
			reply api.BucketInfoReply
		)
		assert.Nil(t, DefaultLeaderRPC.Call(APIServiceBucketInfo, &arg, &reply))
		assert.Equal(t, false, reply.Has)
	}
	////////////////////////////////////////////////
	// set value
	{
		var (
			arg   = api.SetKVArg{BucketName: bucketName, Key: key, Value: value}
			reply api.NotReply
		)
		assert.Nil(t, DefaultLeaderRPC.Call(APIServiceSetKV, &arg, &reply))
	}

	// get value
	{
		var (
			arg   = api.GetKVArg{BucketName: bucketName, Key: key}
			reply api.GetKVReply
		)
		assert.Nil(t, DefaultLeaderRPC.Call(APIServiceKV, &arg, &reply))
		fmt.Println(reply.Has)
	}
	// del value
	{
		{
			var (
				arg   = api.DelKVArg{BucketName: bucketName, Key: key}
				reply api.NotReply
			)
			assert.Nil(t, DefaultLeaderRPC.Call(APIServiceDelKey, &arg, &reply))
		}

		{
			var (
				arg   = api.GetKVArg{BucketName: bucketName, Key: key}
				reply api.GetKVReply
			)
			assert.Nil(t, DefaultRaftRPC.Call(APIServiceKV, &arg, &reply))
			assert.Equal(t, false, reply.Has)
		}
	}
}
