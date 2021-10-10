package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type QBittorrentClient struct {
	client   *http.Client
	Address  string
	Username string
	Password string
	baseURL  string
	sid      string
}

// NewQBittorrentClient returns authenticated QBittorrentClient.
func NewQBittorrentClient(Address, Username, Password string) (*QBittorrentClient, error) {
	c := &QBittorrentClient{
		client:   http.DefaultClient,
		Address:  Address,
		Username: Username,
		Password: Password,
		baseURL:  fmt.Sprintf("%s/api/v2", Address),
	}
	if err := c.login(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *QBittorrentClient) login() error {
	loginInfo := url.Values{}
	loginInfo.Set("username", c.Username)
	loginInfo.Set("password", c.Password)

	resp, err := http.PostForm(fmt.Sprintf("%s/auth/login", c.baseURL), loginInfo)
	if err != nil {
		return err
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			c.sid = cookie.Value
		}
	}

	return nil
}

type Status struct {
	Connection        string `json:"connection_status"`
	DHTNodes          uint   `json:"dht_nodes"`
	Downloaded        uint64 `json:"dl_info_data"`
	DownloadSpeed     uint   `json:"dl_info_speed"`
	DownloadRateLimit uint   `json:"dl_rate_limit"`
	Uploaded          uint64 `json:"up_info_data"`
	UploadSpeed       uint   `json:"up_info_speed"`
	UploadRateLimit   uint   `json:"up_rate_limit"`
}

func (c *QBittorrentClient) GetStatus() (Status, error) {
	var status Status

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/transfer/info", c.baseURL), nil)
	if err != nil {
		return status, err
	}
	req.AddCookie(&http.Cookie{Name: "SID", Value: c.sid})

	resp, err := c.client.Do(req)
	if err != nil {
		return status, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return status, err
	}

	return status, nil
}

func (c *QBittorrentClient) GetCategories() ([]string, error) {
	type category struct {
		Name     string `json:"name"`
		SavePath string `json:"savePath"`
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/torrents/categories", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "SID", Value: c.sid})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var categories = make(map[string]category)
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, err
	}

	var result = make([]string, 0, len(categories))
	for _, v := range categories {
		result = append(result, v.Name)
	}

	return result, nil
}

type Torrent struct {
	AddedOn                uint    `json:"added_on"`
	AmountLeft             uint    `json:"amount_left"`
	AutoTMM                bool    `json:"auto_tmm"`
	Category               string  `json:"category"`
	Completed              uint64  `json:"completed"`
	CompletionOn           uint    `json:"completion_on"`
	DownloadLimit          int     `json:"dl_limit"`
	DownloadSpeed          uint    `json:"dlspeed"`
	Downloaded             uint64  `json:"downloaded"`
	DownloadedSession      uint64  `json:"downloaded_session"`
	ETA                    uint    `json:"eta"`
	FirstLastPiecePriority bool    `json:"f_l_piece_prio"`
	ForceStart             bool    `json:"force_start"`
	Hash                   string  `json:"hash"`
	LastActivity           uint    `json:"last_activity"`
	MagnetURI              string  `json:"magnet_uri"`
	MaxRatio               float64 `json:"max_ratio"`
	MaxSeedingTime         int     `json:"max_seeding_time"`
	Name                   string  `json:"name"`
	NumComplete            uint    `json:"num_complete"`
	NumIncomplete          uint    `json:"num_incomplete"`
	NumLeechs              uint    `json:"num_leechs"`
	NumSeeds               uint    `json:"num_seeds"`
	Priority               uint    `json:"priority"`
	Progress               uint    `json:"progress"`
	Ratio                  float64 `json:"ratio"`
	RatioLimit             int     `json:"ratio_limit"`
	SavePath               string  `json:"save_path"`
	SeedingTimeLimit       int     `json:"seeding_time_limit"`
	SeenComplete           uint64  `json:"seen_complete"`
	SeqDownload            bool    `json:"seq_dl"`
	Size                   uint64  `json:"size"`
	State                  string  `json:"state"`
	SuperSeeding           bool    `json:"super_seeding"`
	Tags                   string  `json:"tags"`
	TimeActive             uint    `json:"time_active"`
	TotalSize              uint64  `json:"total_size"`
	Tracker                string  `json:"tracker"`
	UploadLimit            int     `json:"up_limit"`
	Uploaded               uint64  `json:"uploaded"`
	UploadedSession        uint64  `json:"uploaded_session"`
	UploadSpeed            uint    `json:"upspeed"`
}

func (c *QBittorrentClient) GetTorrents() ([]Torrent, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/torrents/info", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "SID", Value: c.sid})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var torrents []Torrent
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return nil, err
	}

	return torrents, nil
}
