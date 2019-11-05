package dbservice

type accountUtil struct {
	uid uint64
}

const (
	// AccountPrefix 帐号表前缀
	AccountPrefix = "account"

	// AccountOpenID 用户名表前缀, 存储帐号和UID的对应关系
	AccountOpenID = "accountopenid"

	// UIDField UID字段
	UIDField = "uid"
)
