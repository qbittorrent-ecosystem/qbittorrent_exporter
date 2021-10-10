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
	DHTNodes          int    `json:"dht_nodes"`
	Downloaded        int    `json:"dl_info_data"`
	DownloadSpeed     int    `json:"dl_info_speed"`
	DownloadRateLimit int    `json:"dl_rate_limit"`
	Uploaded          int    `json:"up_info_data"`
	UploadSpeed       int    `json:"up_info_speed"`
	UploadRateLimit   int    `json:"up_rate_limit"`
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
	AddedOn                int     `json:"added_on"`
	AmountLeft             int     `json:"amount_left"`
	AutoTMM                bool    `json:"auto_tmm"`
	Category               string  `json:"category"`
	Completed              int     `json:"completed"`
	CompletionOn           int     `json:"completion_on"`
	DownloadLimit          int     `json:"dl_limit"`
	DownloadSpeed          int     `json:"dlspeed"`
	Downloaded             int     `json:"downloaded"`
	DownloadedSession      int     `json:"downloaded_session"`
	ETA                    int     `json:"eta"`
	FirstLastPiecePriority bool    `json:"f_l_piece_prio"`
	ForceStart             bool    `json:"force_start"`
	Hash                   string  `json:"hash"`
	LastActivity           int     `json:"last_activity"`
	MagnetURI              string  `json:"magnet_uri"`
	MaxRatio               float64 `json:"max_ratio"`
	MaxSeedingTime         int     `json:"max_seeding_time"`
	Name                   string  `json:"name"`
	NumComplete            int     `json:"num_complete"`
	NumIncomplete          int     `json:"num_incomplete"`
	NumLeechs              int     `json:"num_leechs"`
	NumSeeds               int     `json:"num_seeds"`
	Priority               int     `json:"priority"`
	Progress               int     `json:"progress"`
	Ratio                  float64 `json:"ratio"`
	RatioLimit             int     `json:"ratio_limit"`
	SavePath               string  `json:"save_path"`
	SeedingTimeLimit       int     `json:"seeding_time_limit"`
	SeenComplete           int     `json:"seen_complete"`
	SeqDownload            bool    `json:"seq_dl"`
	Size                   int     `json:"size"`
	State                  string  `json:"state"`
	SuperSeeding           bool    `json:"super_seeding"`
	Tags                   string  `json:"tags"`
	TimeActive             int     `json:"time_active"`
	TotalSize              int     `json:"total_size"`
	Tracker                string  `json:"tracker"`
	UploadLimit            int     `json:"up_limit"`
	Uploaded               int     `json:"uploaded"`
	UploadedSession        int     `json:"uploaded_session"`
	UploadSpeed            int     `json:"upspeed"`
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
