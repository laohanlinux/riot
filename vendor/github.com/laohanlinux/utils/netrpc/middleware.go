package netrpc

import (
	"fmt"

	"golang.org/x/net/context"
)

func BeforeExcuteFuncMiddleware() func(*Request) context.Context {
	return func(req *Request) context.Context {
		if req == nil || req.MetaData == nil {
			fmt.Printf("req:%v\n", req)
			return nil
		}
		return NewContext(context.Background(), req.MetaData)
	}
}
