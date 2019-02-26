package s3

import (
	"fmt"

	"github.com/golang/glog"
)

// Implements Mounter
type cosMounter struct {
	bucket        *bucket
	url           string
	region        string
	pwFileContent string
	appid         string
}

const (
	cosCmd = "s3fs"
)

func newCosMounter(b *bucket, cfg *Config) (Mounter, error) {
	return &cosMounter{
		bucket:        b,
		url:           cfg.Endpoint,
		region:        cfg.Region,
		pwFileContent: cfg.AccessKeyID + ":" + cfg.SecretAccessKey,
		appid:         cfg.Appid,
	}, nil
}

func (cos *cosMounter) Stage(stageTarget string) error {
	return nil
}

func (cos *cosMounter) Unstage(stageTarget string) error {
	return nil
}

func (cos *cosMounter) Mount(source string, target string) error {
	if err := writes3fsPass(cos.pwFileContent); err != nil {
		return err
	}
	args := []string{
		fmt.Sprintf("%s-%s:/%s", cos.bucket.Name, cos.appid, cos.bucket.FSPath),
		fmt.Sprintf("%s", target),
		"-o", "sigv2",
		"-o", "use_path_request_style",
		"-o", fmt.Sprintf("url=%s", cos.url),
		"-o", "allow_other",
		"-o", "mp_umask=000",
	}

	glog.Info("Mount args: ", args)
	return fuseMount(target, cosCmd, args)
}

func (cos *cosMounter) Unmount(target string) error {
	return fuseUnmount(target, cosCmd)
}
