package routing

type IRoutingTable interface {
	KBuckets() []IKBucket
	KSize() KSize
}
