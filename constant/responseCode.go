package constant

const (
	// CreateNewChat 开启聊天
	StartChatSuccess = 100
	StartChatError   = 101

	//InitChat 主页面渲染
	ShowChatHistorySuccess = 110
	ShowChatHistoryError   = 111

	//CallBot 调用机器人
	CallBotSuccess = 120
	CallBotError   = 121

	//AdminGetBot 管理员获取机器人信息
	AdminGetBotSuccess = 900
	AdminGetBotError   = 901

	//AdminSaveNewBot 管理员创建新机器人
	AdminSaveNewBotSuccess = 910
	AdminSaveNewBotError   = 911

	//AdminModifyBot 管理员修改机器人
	AdminModifyBotSuccess = 920
	AdminModifyBotError   = 921

	//用户注册
	RegisterSuccess = 1000
	RegisterError   = 1001

	//用户登录
	LoginSuccess = 1010
	LoginError   = 1011

	//用户名和密码解析
	UserGetError = 1021
	//用户名或密码为空
	UserExistNull = 1021
)
