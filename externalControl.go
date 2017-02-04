package aurora

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func (c *Client) ActivateExtranalControl() (chan<- *PanelColorCommand, error) {

	//just gonna hard code this request
	r := strings.NewReader(`{ "write": {"command":"display","version":"1.0","animType":"extControl"} }`)
	req, err := http.NewRequest("PUT", c.url(effectsURL), r)
	if err != nil {
		return nil, err
	}
	dat := &struct {
		IP    string `json:"streamControlIpAddr"`
		Port  int    `json:"streamControlPort"`
		Proto string `json:"streamControlProtocol"`
	}{}
	err = c.makeReq(req, dat)
	if err != nil {
		return nil, err
	}
	fmt.Println(dat.IP, dat.Port)
	addr := fmt.Sprintf("%s:%d", dat.IP, dat.Port)
	ch := make(chan *PanelColorCommand)
	go externalControl(addr, ch)
	return ch, nil
}

type PanelColorCommand struct {
	ID, R, G, B byte
}

func externalControl(addr string, ch <-chan *PanelColorCommand) {

	updates := map[byte][]byte{} //data to eventually send
	var nextPossibleSend time.Time
	var timeToSend <-chan time.Time

	for {
		select {
		case <-timeToSend:
			conn, err := net.Dial("udp", addr)
			if err != nil {
				log.Println("ERROR DIALING:", err)
				continue
			}
			conn.Write(buildUdpPacket(updates))
			conn.Close()
			updates = map[byte][]byte{}
			nextPossibleSend = time.Now().Add(100 * time.Millisecond)
			timeToSend = nil
			fmt.Println("SENT")
		case dp := <-ch:
			//start ticker if not going
			if timeToSend == nil {
				if time.Now().After(nextPossibleSend) {
					//have not sent recently. Collect for another 1ms and send
					nextPossibleSend = time.Now().Add(5 * time.Millisecond)
				}
				d := nextPossibleSend.Sub(time.Now())
				timeToSend = time.After(d)
			}
			updates[dp.ID] = []byte{1, dp.R, dp.G, dp.B, 0, 1}
		}
	}
}

func buildUdpPacket(dat map[byte][]byte) []byte {

	buf := make([]byte, 0, len(dat)*6+1)
	buf = append(buf, byte(len(dat)))
	for id, pt := range dat {
		buf = append(buf, id)
		buf = append(buf, pt...)
	}
	fmt.Println(buf)
	return buf
}
