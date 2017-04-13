// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffpool

import (
	"sync"
)

var bufTypes = [...]int{
	//16
	16, 32, 48, 64, 80, 96, 112, 128,
	//32
	160, 192, 224, 256,
	//64
	320, 384, 448, 512,
	//128
	640, 768, 896, 1024,
	//256
	1280, 1536, 1792, 2048,
	//512
	2560, 3072, 3584, 4096,
}

const bufTypeNum = len(bufTypes)

var bufPools [bufTypeNum]sync.Pool

func init() {
	for i := 0; i < bufTypeNum; i++ {
		l := bufTypes[i]
		bufPools[i].New = func() interface{} {
			return make([]byte, l, l)
		}
	}
}

func BufGet(size int) []byte {
	if size <= 0 {
		return make([]byte, 0, 0)
	} else if size > bufTypes[bufTypeNum-1] {
		return make([]byte, size, size)
	}
	for i := 0; i < bufTypeNum; i++ {
		if size <= bufTypes[i] {
			return bufPools[i].Get().([]byte)[0:size]
		}
	}

	return make([]byte, size, size)
}

func BufPut(b []byte) {
	size := cap(b)
	if size < bufTypes[0] || size > bufTypes[bufTypeNum-1] {
		return
	}
	for i := 1; i < bufTypeNum; i++ {
		if size <= bufTypes[i] {
			if size == bufTypes[i] {
				bufPools[i].Put(b[0:size])
			} else {
				bufPools[i-1].Put(b[0:bufTypes[i-1]])
			}
			return
		}
	}
}
