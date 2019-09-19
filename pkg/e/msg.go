package e

var MsgFlags = map[int]string{
	SUCCESS:                  "ok",
	ERROR:                    "fail",
	INVALID_PARAMS:           "请求参数错误",
	JSON_MARSHAL_ERROR:       "Json序列化失败",
	CREATE_SHORTLINK_ERROR:   "创建短链接失败",
	SHORT_LINK_IS_MUST:       "未找到短链接",
	SHORT_TO_LONG_LINK_ERROR: "短链接转长链接失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
