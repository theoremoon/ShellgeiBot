package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type BotConfigJson struct {
	DockerImage string   `json:"dockerimage"`
	Workdir     string   `json:"workdir"`
	Memory      string   `json:"memory"`
	MediaSize      int64   `json:"mediasize"`
	Timeout     string   `json:"timeout"`
	Tags        []string `json:"tags"`
}

type BotConfig struct {
	DockerImage string
	Workdir     string
	Memory      string
	MediaSize      int64
	Timeout     time.Duration
	Tags        []string
}

func ParseBotConfig(file string) (BotConfig, error) {
	var c BotConfigJson
	var config BotConfig

	// read json
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return config, err
	}

	// convert json to config type
	config.DockerImage = c.DockerImage
	config.Workdir, err = filepath.Abs(c.Workdir)
	if err != nil {
		return config, err
	}
	config.Memory = c.Memory // TODO: check memory size string
	config.MediaSize = c.MediaSize
	config.Timeout, err = time.ParseDuration(c.Timeout)
	if err != nil {
		return config, err
	}
	config.Tags = c.Tags
	return config, nil
}

func RandStr(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randstr := make([]byte, 0, length)
	keys := make([]byte, length)
	_, err := rand.Read(keys)
	if err != nil {
		return "", err
	}
	for _, v := range keys {
		k := int(v) % len(chars)
		randstr = append(randstr, chars[k])
	}

	return string(randstr), nil
}

// https://golangcode.com/download-a-file-from-a-url/
// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

type StdError struct {
	Msg string
}

func (e *StdError) Error() string {
	return e.Msg
}

func RunCmd(cmdstr string, media_urls []string, botConfig BotConfig) (string, []string, error) {
	// create shellgei script file and write shellgei content
	name, err := RandStr(16)
	if err != nil {
		return "", []string{}, err
	}
	path := filepath.Join(botConfig.Workdir, name)
	file, err := os.Create(path)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %s, directory permission denied?", err)
	}
	defer func() { _ = os.RemoveAll(path) }()
	_, err = file.WriteString(cmdstr)
	if err != nil {
		return "", []string{}, fmt.Errorf("errors: %s, failed to write", err)
	}
	file.Close()

	// create images directory
	imgdir_path := filepath.Join(botConfig.Workdir, name+"__images")
	err = os.MkdirAll(imgdir_path, 0777)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %s, could not create directory", err)
	}
	defer func() { _ = os.RemoveAll(imgdir_path) }()

	// create media directory
	mediadir_path := filepath.Join(botConfig.Workdir, name+"__media")
	err = os.MkdirAll(mediadir_path, 0777)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %s, could not create directory", err)
	}
	defer func() { _ = os.RemoveAll(mediadir_path) }()

	// download medias
	for i, url := range media_urls {
		err = DownloadFile(filepath.Join(mediadir_path, strconv.Itoa(i)), url)
		if err != nil {
			return "", nil, fmt.Errorf("error: %s, failed to download a media", err)
		}
	}

	// execute shellgei in the docker
	cmd := exec.Command("docker", "run", "--rm",
		"--net=none",
		"-m", botConfig.Memory,
		"--oom-kill-disable",
		"--pids-limit", "1024",
		"--name", name,
		"-v", path+":/"+name,
		"-v", imgdir_path+":/images",
		"-v", mediadir_path+":/media",
		botConfig.DockerImage,
		"bash", "-c", fmt.Sprintf("chmod +x /%s && sync &&  ./%s | stdbuf -o0 head -c 100K", name, name))

	// get result
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, botConfig.Timeout)
	defer cancel()

	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		errChan <- cmd.Run()
	}(ctx)

	select {
	case <-ctx.Done():
		// kill send SIGKILL immediately
		// though stop send SIGKILL after sending SIGTERM
		_ = exec.Command("docker", "kill", name).Run()
	case err = <-errChan:
		// do nothing
	}

	// search image data
	files, err := ioutil.ReadDir(imgdir_path)

	// without image
	if err != nil || len(files) == 0 {
		return out.String(), []string{}, nil
	}

	// with image
	b64imgs := make([]string, 0, 4)
	readcount := 0

	for i := 0; readcount < 4; i++ {
		if len(files) <= i {
			break
		}
		path := filepath.Join(imgdir_path, files[i].Name())

		// do not follow the symlink
		lfinfo, err := os.Lstat(path)
		if err != nil || lfinfo.Mode() & os.ModeSymlink != 0 {
			continue
		}

		// if file size is zero or bigger than MediaSize[MB]
		finfo, err := os.Stat(path)
		if err != nil || finfo.Size() == 0 || finfo.Size() >= 1024 * 1024 * botConfig.MediaSize {
			continue
		}

		// read image file into memory
		img, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			continue
		}

		// encode to base64
		b64img := base64.StdEncoding.EncodeToString(img)
		b64imgs = append(b64imgs, b64img)
		readcount += 1
	}
	return out.String(), b64imgs, nil
}
