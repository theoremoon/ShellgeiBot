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

type botConfigJSON struct {
	DockerImage string   `json:"dockerimage"`
	Workdir     string   `json:"workdir"`
	Memory      string   `json:"memory"`
	MediaSize   int64    `json:"mediasize"`
	Timeout     string   `json:"timeout"`
	Tags        []string `json:"tags"`
}

type botConfig struct {
	DockerImage string
	Workdir     string
	Memory      string
	MediaSize   int64
	Timeout     time.Duration
	Tags        []string
}

func parseBotConfig(file string) (botConfig, error) {
	var c botConfigJSON
	var config botConfig

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

func randStr(length int) (string, error) {
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

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// https://golangcode.com/download-a-file-from-a-url/
func downloadFile(filepath string, url string) error {

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

type stdError struct {
	Msg string
}

func (e *stdError) Error() string {
	return e.Msg
}

func runCmd(cmdstr string, mediaUrls []string, config botConfig) (string, []string, error) {
	// create shellgei script file and write shellgei content
	name, err := randStr(16)
	if err != nil {
		return "", []string{}, err
	}
	path := filepath.Join(config.Workdir, name)
	file, err := os.Create(path)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %v, directory permission denied?", err)
	}
	defer func() { err := os.RemoveAll(path); log.Println(err) }()
	_, err = file.WriteString(cmdstr)
	if err != nil {
		return "", []string{}, fmt.Errorf("errors: %v, failed to write", err)
	}
	file.Close()

	// use images volume intead of directory
	// c.f. https://github.com/theoldmoon0602/ShellgeiBot/issues/41
	imagesVolume := name + "__volume"
	defer func() {
		_ = exec.Command("docker", "volume", "rm", imagesVolume).Run()
	}()

	// create media directory
	mediadirPath := filepath.Join(config.Workdir, name+"__media")
	err = os.MkdirAll(mediadirPath, 0777)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %v, could not create directory", err)
	}
	defer func() { _ = os.RemoveAll(mediadirPath) }()

	// download medias
	for i, url := range mediaUrls {
		err = downloadFile(filepath.Join(mediadirPath, strconv.Itoa(i)), url)
		if err != nil {
			return "", nil, fmt.Errorf("error: %v, failed to download a media", err)
		}
	}

	// execute shellgei in the docker
	cmd := exec.Command("docker", "run", "--rm",
		"--net=none",
		"-m", config.Memory,
		"--oom-kill-disable",
		"--pids-limit", "1024",
		"--name", name,
		"-v", path+":/"+name,
		"-v", imagesVolume+":/images",
		"-v", mediadirPath+":/media:ro",
		config.DockerImage,

		"bash", "-c", fmt.Sprintf("chmod +x /%s && sync &&  ./%s | stdbuf -o0 head -c 100K | stdbuf -o0 head -n 15", name, name))

	// get result
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
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
	case <-errChan:
		// do nothing
	}

	// create images directory
	imgdirPath := filepath.Join(config.Workdir, name+"__images")
	err = os.MkdirAll(imgdirPath, 0777)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %v, could not create directory", err)
	}
	defer func() { err := os.RemoveAll(imgdirPath); log.Println(err) }()

	// get images from docker volume
	if err := getImagesFromDockerVolume(imgdirPath, imagesVolume, config.MediaSize); err != nil {
		log.Println(err)
	}

	// search image data
	b64img, err := encodeImages(imgdirPath, config.MediaSize)
	return out.String(), b64img, err
}

func getImagesFromDockerVolume(dstPath, vol string, size int64) error {
	// do not use 'cp'. special device files hurts the system
	sizeStr := strconv.FormatInt(size*1024*1024, 10)
	return exec.Command("docker", "run", "--rm", "-v", dstPath+":/dst", "-v", vol+":/src", "bash", "-c", "ls -A -1d /src/* | while read -r f; do [[ -f \"$f\" ]] && head -c "+sizeStr+" \"$f\" > \"${f/#\\/src/\\/dst}\"; done").Run()
}

func encodeImages(imgdirPath string, size int64) ([]string, error) {
	files, err := ioutil.ReadDir(imgdirPath)
	if err != nil || len(files) == 0 {
		return []string{}, nil
	}
	// with image
	b64imgs := make([]string, 0, 4)
	readcount := 0

	for i := 0; readcount < 4; i++ {
		if len(files) <= i {
			break
		}
		path := filepath.Join(imgdirPath, files[i].Name())

		// do not follow the symlink
		lfinfo, err := os.Lstat(path)
		if err != nil || lfinfo.Mode()&os.ModeSymlink != 0 {
			continue
		}

		// if file size is zero or bigger than MediaSize[MB]
		finfo, err := os.Stat(path)
		if err != nil || finfo.Size() == 0 || finfo.Size() >= 1024*1024*size {
			continue
		}

		// unnecessary because [[ -f "$f" ]] checks this
		// // check file is regular to avoid read special files
		// // e.g. /dev/zero, named pipe, etc.
		// if !finfo.Mode().IsRegular() {
		// 	continue
		// }

		// read image file into memory
		img, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			continue
		}

		// encode to base64
		b64img := base64.StdEncoding.EncodeToString(img)
		b64imgs = append(b64imgs, b64img)
		readcount++
	}
	return b64imgs, nil
}
