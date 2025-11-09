package ftp

import (
	"bytes"
	"context"
	"github.com/jlaffaye/ftp"
	"io"
	"regexp"
)

func ReadFileByUrl(ctx context.Context, url string) ([]byte, error) {

	// parse ftp url
	pu, err := parseUrl(ctx, url)
	if err != nil {
		return nil, err
	}

	// connect to FTP server
	c, err := ftp.Dial(pu.HostAndPort)
	if err != nil {
		return nil, ErrFtpServerConnect(ctx, err, url)
	}

	// close
	defer func() { _ = c.Quit() }()

	// login
	err = c.Login(pu.Login, pu.Password)
	if err != nil {
		return nil, ErrFtpServerLogin(ctx, err, url)
	}

	r, err := c.Retr(pu.Path)
	if err != nil {
		return nil, ErrFtpRead(ctx, err, url)
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, ErrFtpRead(ctx, err, url)
	}

	return buf, nil
}

func LoadFileByUrl(ctx context.Context, url string, content []byte) error {

	// parse ftp url
	pu, err := parseUrl(ctx, url)
	if err != nil {
		return err
	}

	// connect to FTP server
	c, err := ftp.Dial(pu.HostAndPort)
	if err != nil {
		return ErrFtpServerConnect(ctx, err, url)
	}

	// close
	defer func() { _ = c.Quit() }()

	// login
	err = c.Login(pu.Login, pu.Password)
	if err != nil {
		return ErrFtpServerLogin(ctx, err, url)
	}

	err = c.Stor(pu.Path, bytes.NewReader(content))
	if err != nil {
		return ErrFtpLoad(ctx, err, url)
	}

	return nil
}

type parsedUrl struct {
	Schema      string
	Login       string
	Password    string
	HostAndPort string
	Path        string
}

func parseUrl(ctx context.Context, url string) (*parsedUrl, error) {

	// parse ftp url
	re, err := regexp.Compile(`(ftp|ftps)://([a-zA-Z0-9\-]+):([a-zA-Z0-9\-]+)@([a-zA-Z0-9\-.:]+)/(.+)`)
	if err != nil {
		return nil, ErrFtpUrlParse(ctx, err, url)
	}

	// extract parts from the url
	m := re.FindStringSubmatch(url)
	if len(m) != 6 {
		return nil, ErrFtpUrlParse(ctx, err, url)
	}

	return &parsedUrl{
		Schema:      m[1],
		Login:       m[2],
		Password:    m[3],
		HostAndPort: m[4],
		Path:        m[5],
	}, nil

}
