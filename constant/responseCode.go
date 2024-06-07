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

	//UserGetHistory 用户获取历史记录
	UserGetHistorySuccess = 130
	UserGetHistoryError   = 131

	//UserShareHistory 用户获取历史记录
	UserShareHistorySuccess = 140
	UserShareHistoryError   = 141
	UserShareHistoryNil     = 142

	//UserGetSharedHistory 用户获取分享时的历史记录
	UserGetSharedHistorySuccess = 150
	UserGetSharedHistoryError   = 151

	//UpdateTitle 更新标题
	UpdateTitleSuccess = 160
	UpdateTitleError   = 161

	//ShowRecycled 展示回收站
	ShowRecycledSuccess = 170
	ShowRecycledError   = 171

	//LogicalDelete 逻辑删除
	LogicalDeleteSuccess = 170
	LogicalDeleteError   = 171

	//RemoveRecycled 取消逻辑删除
	RemoveRecycledSuccess = 170
	RemoveRecycledError   = 171

	//AdminGetBot 管理员获取机器人信息
	AdminGetBotSuccess = 900
	AdminGetBotError   = 901

	//AdminSaveNewBot 管理员创建新机器人
	AdminSaveNewBotSuccess = 910
	AdminSaveNewBotError   = 911

	//AdminModifyBot 管理员修改机器人
	AdminModifyBotSuccess = 920
	AdminModifyBotError   = 921

	//Rpc4Title rpc获取标题
	Rpc4TitleSuccess = 930
	Rpc4TitleError   = 931
)
