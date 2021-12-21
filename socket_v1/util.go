package socket_v1

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"sync"
	"time"
)

// 全局变量：生成不重复客户端id用
var (
	lastAutoId     int64        // 转成string格式前的客户端id
	lastAutoIdLock sync.RWMutex // 同步锁
)

// 初始化函数：初始化不重复客户端id
func init() {
	lastAutoId = time.Now().UnixNano()
}

// 生成11位单机不重复字串 (耗时短 100000 条需要40毫秒)
func utilUuidShort() string {
	dem := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	retBuf := make([]byte, 0)
	// int64位时间
	lastAutoIdLock.Lock()
	lastAutoId++
	T1 := lastAutoId
	lastAutoIdLock.Unlock()
	Mask := int64(63)
	for i := 0; i < 64; i += 6 {
		retBuf = append(retBuf, dem[T1>>i&Mask])
	}
	return string(retBuf)
}

// 获取当前时间字串 返回："2021-08-25 11:16:20"
func utilDateTime(T ...time.Time) string {
	timeObj := time.Now()
	if len(T) > 0 {
		timeObj = T[0]
	}
	return timeObj.Format("2006-01-02 15:04:05")
}

// 进行zlib压缩
func utilZLibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, _ = w.Write(src)
	_ = w.Close()
	return in.Bytes()
}

// 进行zlib解压缩
func utilZLibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	_, _ = io.Copy(&out, r)
	return out.Bytes()
}

// 整形转换成字节
func utilInt2Bytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 字节转换成整形
func utilBytes2Int(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func utilInt64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func utilBytes2Int64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
