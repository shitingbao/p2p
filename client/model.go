package client

const DefaultUri = "stun:stun.l.google.com:19302"

type Client struct {
	clientId string // 在所有连接中自己定义唯一的值，比如依托于另一个权限平台的用户id
	StunRaw  string
	FlagHost string // 注册地址，也是用户交互的地址
}

type option struct {
	StunRaw  string
	FlagHost string // 注册地址，也是用户交互的地址
}

type Option func(o *option)
