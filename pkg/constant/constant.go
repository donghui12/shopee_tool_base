package constant

const WorkerPoolSize = 20
const RetryTimes = 200
const FailedMaxCount = 10
const DefaultPageSize = "100"
const MaxProductNum = 30000

const (
	ShopeeSessionKey = "SPC_CNSC_SESSION"
)

const (
	TopicShop    string = "shop"
	TopicProduct string = "product"
)

var Topics = []string{
	TopicShop,
	TopicProduct,
}

const (
	PartnerId  = "2010576"
	PartnerKey = "4a657554534646714f79684d77586f624d624f52617458784865697561454579"

	TestPartnerId  = "1257803"
	TestPartnerKey = "54466b6d594645716175747a7362665172586966584f45486f45456e616e5861"

	LivePartnerId  = "2010576"
	LivePartnerKey = "4a657554534646714f79684d77586f624d624f52617458784865697561454579"
)

const IsVip = true
const IsOpen = true

const (
	ActionDeleted  = "delete" // 删除不活跃账户
	ActionListed   = "listed" // 上架全部商品
	ActionUnlisted = "unlist" // 下架全部商品

	ActionStopDiscount       = "stop"   // 停止折扣
	ActionCopyCreateDiscount = "create" // 停止折扣
)

const (
	ProxyHost     = "https://share.proxy.qg.net"
	ProxyAuthKey  = "80C9D1B4"
	ProxyPassword = "E0043D53BF8B"
)
