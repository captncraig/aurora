package aurora

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	addr  string
	token string
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func NewWithToken(addr string, token string) *Client {
	return &Client{
		addr:  addr,
		token: token,
	}
}

var (
	baseURL    = "/api/beta"
	authURL    = "/new"
	infoURL    = "/%s"
	effectsURL = "/%s/effects"
)

// Authorize will attempt to get an auth token from the device. Device must be put into pairing mode (hold down power button 5-7 seconds) for this to work.
// On success will return a valid auth token.
func (c *Client) Authorize() error {
	resp, err := http.Post(c.url(authURL), "application/json", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad status code from nanoleaf device: %d", resp.StatusCode)
	}

	dat := &struct {
		Token string `json:"auth_token"`
	}{}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = dec.Decode(dat); err != nil {
		return err
	}
	c.token = dat.Token
	return nil
}

type HardwareInfo struct {
	Name            string `json:"name"`
	SerialNo        string `json:"serialNo"`
	Manufacturer    string `json:"manufacturer"`
	FirmwareVersion string `json:"firmwareVersion"`
	Model           string `json:"model"`
	State           struct {
		On struct {
			Value bool `json:"value"`
		} `json:"on"`
		Brightness struct {
			Value int `json:"value"`
			Max   int `json:"max"`
			Min   int `json:"min"`
		} `json:"brightness"`
		Hue struct {
			Value int `json:"value"`
			Max   int `json:"max"`
			Min   int `json:"min"`
		} `json:"hue"`
		Sat struct {
			Value int `json:"value"`
			Max   int `json:"max"`
			Min   int `json:"min"`
		} `json:"sat"`
		Ct struct {
			Value int `json:"value"`
			Max   int `json:"max"`
			Min   int `json:"min"`
		} `json:"ct"`
		ColorMode string `json:"colorMode"`
	} `json:"state"`
	Effects struct {
		Select string   `json:"select"`
		List   []string `json:"list"`
	} `json:"effects"`
	Panels      []*Panel `json:"panels"`
	PanelLayout struct {
		Layout struct {
			LayoutData string `json:"layoutData"`
		} `json:"layout"`
		GlobalOrientation struct {
			Value int `json:"value"`
			Max   int `json:"max"`
			Min   int `json:"min"`
		} `json:"globalOrientation"`
	} `json:"panelLayout"`
}

type Panel struct {
	ID         int `json:"id"`
	X          int `json:"x"`
	Y          int `json:"y"`
	Rotation   int `json:"rotation"`
	SideLength int `json:"length"`
}

func (c *Client) GetInfo() (*HardwareInfo, error) {
	req, err := http.NewRequest("GET", c.url(infoURL), nil)
	if err != nil {
		return nil, err
	}
	dat := &HardwareInfo{}
	if err = c.makeReq(req, dat); err != nil {
		return nil, err
	}
	parts := strings.Split(dat.PanelLayout.Layout.LayoutData, " ")
	if len(parts) <= 2 {
		return dat, nil
	}
	n, _ := strconv.Atoi(parts[0])
	side, _ := strconv.Atoi(parts[1])
	if len(parts) < (n*4)+2 {
		return nil, fmt.Errorf("Invalid panel layout data")
	}
	for i := 0; i < n; i++ {
		p := &Panel{
			SideLength: side,
		}
		p.ID, _ = strconv.Atoi(parts[2+i*4])
		p.X, _ = strconv.Atoi(parts[3+i*4])
		p.Y, _ = strconv.Atoi(parts[4+i*4])
		p.Rotation, _ = strconv.Atoi(parts[5+i*4])
		dat.Panels = append(dat.Panels, p)
	}
	return dat, nil
}

func (c *Client) makeReq(req *http.Request, dat interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad status code from nanoleaf device: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = dec.Decode(dat); err != nil {
		return err
	}
	return nil
}

func (c *Client) url(path string) string {
	u := fmt.Sprintf("%s%s%s", c.addr, baseURL, path)
	if strings.Contains(u, "%s") {
		u = fmt.Sprintf(u, c.token)
	}
	return u
}
