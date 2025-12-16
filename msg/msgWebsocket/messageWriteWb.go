package msg

import "errors"

func (c *WsConn) MessageWrite(msg []byte) (n int, err error) {
	if !c.GetOnline() {
		return 0, errors.New("节点不在线")
	}
	// if len(msg) != 0 {
	// 	c.CounterSub()
	// }
	return c.Write(msg)
}
