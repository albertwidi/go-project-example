package imagepath

import (
	"errors"
	"net"
	"net/url"
	"path"
	"strings"
	"time"

	imageentity "github.com/albertwidi/go-project-example/internal/entity/image"
)

// Config of image path library
type Config struct {
	Public  DownloadConfig
	Private DownloadConfig
}

// DownloadConfig struct
type DownloadConfig struct {
	DownloadProto string
	DownloadHost  string
	DownloadPort  string
	DownloadPath  string
}

// ImagePath struct
type ImagePath struct {
	Proto        string
	Host         string
	FilePath     string
	DownloadPath string
	DownloadLink string
	Signed       bool
}

var _config Config

// Init image path library
func Init(c *Config, local bool) error {
	if c == nil {
		return errors.New("imagepath: config cannot be nil")
	}

	_config = *c

	if local {
		dialer := net.Dialer{Timeout: time.Second * 3}
		// TODO: check when internet connection is down
		conn, err := dialer.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return err
		}

		addr := conn.LocalAddr().(*net.UDPAddr)
		// split the port
		s := strings.Split(addr.String(), ":")
		_config.Public.DownloadHost = s[0]
		_config.Private.DownloadHost = s[0]
	}
	return nil
}

// GenerateImagePath return image path struct with all download information needed
func GenerateImagePath(mode imageentity.Mode, imagePath string) (*ImagePath, error) {
	var img ImagePath

	if err := mode.Validate(); err != nil {
		return nil, err
	}

	switch mode {
	case imageentity.ModePublic:
		u, err := url.Parse(_config.Public.DownloadProto + _config.Public.DownloadHost + _config.Public.DownloadPort)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, _config.Public.DownloadPath, imagePath)

		img = ImagePath{
			Proto:        _config.Public.DownloadProto,
			Host:         _config.Public.DownloadHost,
			FilePath:     imagePath,
			DownloadPath: path.Join(_config.Public.DownloadPath, imagePath),
			DownloadLink: u.String(),
		}

	case imageentity.ModePrivate:
		u, err := url.Parse(_config.Private.DownloadProto + _config.Private.DownloadHost + _config.Private.DownloadPort)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, _config.Private.DownloadPath)

		v := u.Query()
		v.Add("image_path", imagePath)
		u.RawQuery = v.Encode()

		img = ImagePath{
			Proto:        _config.Private.DownloadProto,
			Host:         _config.Private.DownloadHost,
			FilePath:     imagePath,
			DownloadPath: path.Join(_config.Private.DownloadPath, imagePath),
			DownloadLink: u.String(),
		}

	case imageentity.ModeSigned:
		img = ImagePath{
			Proto:        "",
			Host:         "",
			FilePath:     "",
			DownloadPath: "",
			DownloadLink: imagePath,
			Signed:       true,
		}
	}

	return &img, nil
}
