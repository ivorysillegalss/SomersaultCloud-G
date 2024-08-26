package compressutil

import (
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/sys"
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/gzip"
	"github.com/thoas/go-funk"
	"io"
)

type Compress interface {
	// CompressData 序列化数据并进行压缩
	CompressData(data any) ([]byte, error)
	// DecompressData 解压缩后数据直接反序列化回输入的v对象中
	DecompressData(data []byte, v any) error
}

type GzipCompress struct {
}

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

// NewCompress TODO 可补充不同的解压方案
func NewCompress(args ...int) Compress {
	if funk.Equal(args[0], sys.GzipCompress) {
	}
	return &GzipCompress{}
}
