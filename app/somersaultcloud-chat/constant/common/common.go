package common

const (
	ZeroInt    int     = 0     // 整型的零值
	ZeroFloat  float64 = 0.0   // 浮点型的零值
	ZeroBool   bool    = false // 布尔型的零值
	ZeroString string  = ""    // 字符串的零值

	FalseInt int = -1 // 表示无需此参数

	Infix = ":" // 中缀

	LogicalDelete   = 1
	UnLogicalDelete = 0

	True = 1
	False
)

var (
	ZeroSlice []int          = nil // 切片的零值
	ZeroMap   map[string]int = nil // 映射的零值
	ZeroPtr   *int           = nil // 指针的零值
	ZeroByte  []byte         = nil
)

const (
	RecordNotFoundError = "record not found"
	// SystemRole 系统角色
	SystemRole = "system"
	// UserRole 用户角色
	UserRole = "user"
	// GPTRole GPT角色
	GPTRole = "assistant"
	// TextType 文字类型
	TextType = "text"
	// ImageURLType 图片链接类型
	ImageURLType = "image_url"
	// HighDetail 图片分辨率细节
	HighDetail = "high"
	// ApiServerOpenAI OpenAI-API服务器默认网址
	ApiServerOpenAI = "https://api.openai.com/v1/chat/completions"
	// Info 日志级别 info
	Info = "info"
	// Error 日志级别 Error
	Error = "Error"
)
