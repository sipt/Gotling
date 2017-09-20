package utils

import (
	"github.com/zheng-ji/goSnowFlake"
)

var iw *goSnowFlake.IdWorker

func init() {
	var err error
	iw, err = goSnowFlake.NewIdWorker(1)
	if err != nil {
		panic(err)
	}
}

//UniqueID 唯一ID生成
func UniqueID() int64 {
	id, _ := iw.NextId()
	return id
}
