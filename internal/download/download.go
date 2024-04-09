package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"
	"reflect"
	"strings"

	"net/url"
	"os"

	"github.com/kgf1980/go-luxpower/internal/config"
	"golang.org/x/net/publicsuffix"
)

type Download struct {
	StationNumber string
	Username      string
	Password      string
	BaseURL       *url.URL
	Jar           *cookiejar.Jar
	Client        *http.Client
}

type LiveData struct {
	PhotoVoltaic1Watts     int `json:"ppv1"`
	PhotoVoltaic2Watts     int `json:"ppv2"`
	PhotoVoltaic3Watts     int `json:"ppv3"`
	PhotoVoltaicTotalWatts int `json:"ppv"`
	InverterToBattery      int `json:"pCharge"`
	BatteryToInverter      int `json:"pDisCharge"`
	BatteryChargePercent   int `json:"soc"`
	InverterToLoad         int `json:"pinv"`
	GridToLoad             int `json:"pToUser"`
	InverterToGrid         int `json:"pToGrid"`
}

type LiveDataDisplay LiveData

func (ldd LiveDataDisplay) MarshalJSON() ([]byte, error) {
	lddVal := reflect.ValueOf(ldd)
	kvpairs := []string{}

	for i := 0; i < lddVal.NumField(); i++ {
		k := lddVal.Type().Field(i).Name
		v := lddVal.Field(i).Interface() //TODO: proper JSON encoding of things
		kvpairs = append(kvpairs, fmt.Sprintf("\"%s\":%#v", k, v))
	}

	return []byte(fmt.Sprintf("{%s}", strings.Join(kvpairs, ","))), nil
}

type TodayData struct {
	SolarYield         float64 `json:"todayYielding"`
	SolarYieldTotal    float64 `json:"totalYielding"`
	BatteryCharge      float64 `json:"todayCharging"`
	BatteryChargeTotal float64 `json:"totalCharging"`
	Export             float64 `json:"todayExport"`
	ExportTotal        float64 `json:"totalExport"`
	Import             float64 `json:"todayImport"`
	ImportTotal        float64 `json:"totalImport"`
	Usage              float64 `json:"todayUsage"`
	UsageTotal         float64 `json:"totalUsage"`
}

type TodayDataDisplay TodayData

func (tdd TodayDataDisplay) MarshalJSON() ([]byte, error) {
	tddVal := reflect.ValueOf(tdd)
	kvpairs := []string{}

	for i := 0; i < tddVal.NumField(); i++ {
		k := tddVal.Type().Field(i).Name
		v := tddVal.Field(i).Interface() //TODO: proper JSON encoding of things
		kvpairs = append(kvpairs, fmt.Sprintf("\"%s\":%#v", k, v))
	}

	return []byte(fmt.Sprintf("{%s}", strings.Join(kvpairs, ","))), nil
}

func NewDownloader(config *config.Config) (*Download, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	url, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}
	return &Download{
		StationNumber: config.StationNumber,
		Username:      config.AccountName,
		Password:      config.Password,
		BaseURL:       url,
		Jar:           jar,
		Client:        client,
	}, nil
}

func (d *Download) authenticate() error {
	v := url.Values{
		"account":  {d.Username},
		"password": {d.Password},
	}
	postUrl, err := url.Parse(fmt.Sprintf("%s/web/login", d.BaseURL))
	if err != nil {
		return err
	}
	_, err = d.Client.PostForm(postUrl.String(), v)
	if err != nil {
		return err
	}

	return nil
}

func (d *Download) GetTodayData() (*TodayData, error) {
	if len(d.Jar.Cookies(d.BaseURL)) == 0 {
		d.authenticate()
	}
	liveUrl, err := url.Parse(fmt.Sprintf("%s/api/inverter/getInverterEnergyInfo", d.BaseURL.String()))
	if err != nil {
		return &TodayData{}, err
	}

	v := url.Values{
		"serialNum": {d.StationNumber},
	}

	r, err := d.Client.PostForm(liveUrl.String(), v)
	if err != nil {
		return &TodayData{}, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &TodayData{}, err
	}
	var data TodayData
	if err := json.Unmarshal(body, &data); err != nil {
		return &TodayData{}, err
	}
	data.SolarYield = data.SolarYield / 10
	data.SolarYieldTotal = data.SolarYieldTotal / 10
	data.BatteryCharge = data.BatteryCharge / 10
	data.BatteryChargeTotal = data.BatteryChargeTotal / 10
	data.Export = data.Export / 10
	data.ExportTotal = data.ExportTotal / 10
	data.Import = data.Import / 10
	data.ImportTotal = data.ImportTotal / 10
	data.Usage = data.Usage / 10
	data.UsageTotal = data.UsageTotal / 10
	return &data, nil
}

func (d *Download) GetLiveData() (*LiveData, error) {
	if len(d.Jar.Cookies(d.BaseURL)) == 0 {
		d.authenticate()
	}
	liveUrl, err := url.Parse(fmt.Sprintf("%s/api/inverter/getInverterRuntime", d.BaseURL.String()))
	if err != nil {
		return &LiveData{}, err
	}

	v := url.Values{
		"serialNum": {d.StationNumber},
	}

	r, err := d.Client.PostForm(liveUrl.String(), v)
	if err != nil {
		return &LiveData{}, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &LiveData{}, err
	}
	var data LiveData
	if err := json.Unmarshal(body, &data); err != nil {
		return &LiveData{}, err
	}
	return &data, nil
}

func (d *Download) DownloadFile(ExportDate string, Folder string) error {
	filename := fmt.Sprintf("%s - %s.xls", d.StationNumber, ExportDate)
	filename = filepath.Join(Folder, filename)
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(filename), 0750)
	info, err := os.Stat(filename)
	if !errors.Is(err, fs.ErrNotExist) && info.Size() > 0 {
		fmt.Println(filename, " already exists - skipping")
		return nil
	}
	if len(d.Jar.Cookies(d.BaseURL)) == 0 {
		d.authenticate()
	}
	exportUrl, err := url.Parse(fmt.Sprintf("%s/web/analyze/data/export/%s/%s", d.BaseURL.String(), d.StationNumber, ExportDate))
	if err != nil {
		return err
	}

	r, err := d.Client.Get(exportUrl.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()

	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Disposition"))
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %s\n", params["filename"])
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r.Body)
	if err != nil {
		return err
	}
	return nil
}
