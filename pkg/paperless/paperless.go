package paperless

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type Instance struct {
	server    string
	ignoreTLS bool
	username  string
	password  string
}

func New(server, username, password string, ignoretls bool) (*Instance, error) {
	url, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	url.Path = strings.TrimRight(url.Path, "/")
	if len(url.Scheme) == 0 {
		url.Scheme = "http"
	}
	i := &Instance{
		server:    url.String(),
		username:  username,
		password:  password,
		ignoreTLS: ignoretls,
	}
	if len(i.server) == 0 ||
		len(i.username) == 0 ||
		len(i.password) == 0 {
		return nil, errors.New("server, username and password are required")
	}
	return i, nil
}

func (i *Instance) MergePDF(ids []int) error {
	if len(ids) < 2 {
		return errors.New("need to specify atleast two ids")
	}
	tmp, err := os.MkdirTemp("", "plm-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	rsc := make([]string, len(ids))
	for n, v := range ids {
		rsc[n], err = i.DownloadPDF(v, tmp)
		if err != nil {
			return err
		}
	}
	if err := api.MergeAppendFile(rsc[1:], rsc[0], nil); err != nil {
		return err
	}
	return i.UploadPDF(rsc[0])
}

func (i *Instance) DownloadPDF(id int, dstdir string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: i.ignoreTLS}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/documents/%d/download/?original=true", i.server, id), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", i.authHeader())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}
	defer resp.Body.Close()
	fpath := filepath.Join(dstdir, fmt.Sprintf("%d.pdf", id))
	f, err := os.Create(fpath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	return fpath, nil
}

func (i *Instance) UploadPDF(path string) error {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	fw, err := w.CreateFormFile("document", "document.pdf")
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: i.ignoreTLS}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/documents/post_document/", i.server), b)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", i.authHeader())
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed: %s", resp.Status)
	}
	return nil
}

func (i *Instance) DeleteDocuments(ids []int) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: i.ignoreTLS}
	for _, id := range ids {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/documents/%d/", i.server, id), nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", i.authHeader())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK &&
			resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("delete for %d failed: %s", id, resp.Status)
		}
		defer resp.Body.Close()
	}
	return nil
}

func (i *Instance) authHeader() string {
	authstr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", i.username, i.password)))
	return fmt.Sprintf("Basic %s", authstr)
}
