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
	DHTNodes          int64  `json:"dht_nodes"`
	Downloaded        int64  `json:"dl_info_data"`
	DownloadSpeed     int64  `json:"dl_info_speed"`
	DownloadRateLimit int64  `json:"dl_rate_limit"`
	Uploaded          int64  `json:"up_info_data"`
	UploadSpeed       int64  `json:"up_info_speed"`
	UploadRateLimit   int64  `json:"up_rate_limit"`
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
	AddedOn                int64   `json:"added_on"`
	AmountLeft             int64   `json:"amount_left"`
	AutoTMM                bool    `json:"auto_tmm"`
	Category               string  `json:"category"`
	Completed              int64   `json:"completed"`
	CompletionOn           int64   `json:"completion_on"`
	DownloadLimit          int64   `json:"dl_limit"`
	DownloadSpeed          int64   `json:"dlspeed"`
	Downloaded             int64   `json:"downloaded"`
	DownloadedSession      int64   `json:"downloaded_session"`
	ETA                    int64   `json:"eta"`
	FirstLastPiecePriority bool    `json:"f_l_piece_prio"`
	ForceStart             bool    `json:"force_start"`
	Hash                   string  `json:"hash"`
	LastActivity           int64   `json:"last_activity"`
	MagnetURI              string  `json:"magnet_uri"`
	MaxRatio               float64 `json:"max_ratio"`
	MaxSeedingTime         int64   `json:"max_seeding_time"`
	Name                   string  `json:"name"`
	NumComplete            int64   `json:"num_complete"`
	NumIncomplete          int64   `json:"num_incomplete"`
	NumLeechs              int64   `json:"num_leechs"`
	NumSeeds               int64   `json:"num_seeds"`
	Priority               int64   `json:"priority"`
	Progress               float64 `json:"progress"`
	Ratio                  float64 `json:"ratio"`
	RatioLimit             int64   `json:"ratio_limit"`
	SavePath               string  `json:"save_path"`
	SeedingTimeLimit       int64   `json:"seeding_time_limit"`
	SeenComplete           int64   `json:"seen_complete"`
	SeqDownload            bool    `json:"seq_dl"`
	Size                   int64   `json:"size"`
	State                  string  `json:"state"`
	SuperSeeding           bool    `json:"super_seeding"`
	Tags                   string  `json:"tags"`
	TimeActive             int64   `json:"time_active"`
	TotalSize              int64   `json:"total_size"`
	Tracker                string  `json:"tracker"`
	UploadLimit            int64   `json:"up_limit"`
	Uploaded               int64   `json:"uploaded"`
	UploadedSession        int64   `json:"uploaded_session"`
	UploadSpeed            int64   `json:"upspeed"`
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
