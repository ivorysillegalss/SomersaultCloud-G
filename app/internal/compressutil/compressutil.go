package compressutil

import (
	"SomersaultCloud/app/constant/common"
	"SomersaultCloud/app/constant/sys"
	"SomersaultCloud/app/infrastructure/log"
	__proto "SomersaultCloud/app/proto/.proto"
	"bytes"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/gzip"
	"google.golang.org/protobuf/proto"
	"io"
)

// Compress 序列化压缩数据统一接口
// 提供Gzip+Jsontier & Protobuf实现
type Compress interface {
	// CompressData 序列化数据并进行压缩
	CompressData(data any) ([]byte, error)
	// DecompressData 解压缩后数据直接反序列化回输入的v对象中
	DecompressData(data []byte, v any) error
}

type GzipCompress struct{}

func (g GzipCompress) CompressData(data any) ([]byte, error) {
	marshal, err := jsoniter.Marshal(data)
	if err != nil {
		return common.ZeroByte, err
	}

	// gzip 压缩
	var buf bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	defer gzipWriter.Close()

	if err != nil {
		return common.ZeroByte, err
	}
	_, err = gzipWriter.Write(marshal)
	if err != nil {
		return common.ZeroByte, err
	}
	err = gzipWriter.Close()
	if err != nil {
		return common.ZeroByte, err
	}

	return buf.Bytes(), nil
}

func (g GzipCompress) DecompressData(data []byte, v any) error {
	// gzip 解压缩
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(decompressedData, v)
	return nil
}

type ProtoBufCompress struct{}

func (p ProtoBufCompress) CompressData(data any) ([]byte, error) {
	//message, ok := data.([]*proto.Message)
	//if !ok {
	//	log.GetTextLogger().Fatal("data does not implement proto.Message")
	//	return nil, errors.New("data does not implement proto.Message")
	//}
	//TODO 暂时想不到更好的解决办法 赶时间暂时写死类型了
	message, ok := data.([]*__proto.Record)
	if !ok {
		log.GetTextLogger().Fatal("data does not implement proto.Message")
		return nil, errors.New("data does not implement proto.Message")
	}
	return proto.Marshal(&__proto.RecordsList{Records: message})
}

func (p ProtoBufCompress) DecompressData(data []byte, v any) error {
	//message, ok := v.(proto.Message)
	//if !ok {
	//	log.GetTextLogger().Fatal("v does not implement proto.Message")
	//	return errors.New("v does not implement proto.Message")
	//}

	h := new(__proto.RecordsList)
	err := proto.Unmarshal(data, h)
	return err
}

func NewCompress(args ...string) Compress {
	switch args[0] {
	case sys.GzipCompress:
		return &GzipCompress{}
	case sys.ProtoBufCompress:
		return &ProtoBufCompress{}
	default:
		log.GetTextLogger().Fatal("error Compress sign")
		panic("error Compress sign")
	}
}
