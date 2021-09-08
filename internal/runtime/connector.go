package runtime

type Connector interface {
	ConnectSubnet([]Connection)
	ConnectOperator(string, IO) error
}
