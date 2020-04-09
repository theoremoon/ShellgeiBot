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
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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

var dkclient, _ = client.NewEnvClient()

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

	mem, _ := strconv.ParseInt(config.Memory, 10, 64)
	to := int(config.Timeout.Seconds())
	f := false

	// get result
	var out bytes.Buffer
	var stderr bytes.Buffer

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	resp, err := dkclient.ContainerCreate(
		ctx,
		&container.Config{
			Image:           config.DockerImage,
			NetworkDisabled: true,
			Cmd: []string{
				"bash", "-c",
				oneLiner(
					"chmod", "+x", "/"+name, "&& sync &&", "./"+name, "|",
					"stdbuf -o0 head -c 100K", "|",
					"stdbuf -o0 head -n 15",
				),
			},
			AttachStdout: true,
			AttachStderr: true,
			StopTimeout:  &to,
		},
		&container.HostConfig{
			AutoRemove:   true, // AutoRemove を true にすることで --rm と同じになる
			NetworkMode:  "none",
			VolumeDriver: "local",
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					Source:   path,
					Target:   "/" + name,
					ReadOnly: false,
				},
				{
					Type:     mount.TypeVolume,
					Source:   imagesVolume,
					Target:   "/images",
					ReadOnly: false,
				},
				{
					Type:     mount.TypeBind,
					Source:   mediadirPath,
					Target:   "/media",
					ReadOnly: true,
				},
			},
			Resources: container.Resources{
				Memory:         mem,
				OomKillDisable: &f,
				PidsLimit:      1024,
			},
		},
		&network.NetworkingConfig{},
		name,
	)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %v, could not container create correctly", err)
	}

	if err := dkclient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", []string{}, fmt.Errorf("error: %v ContainerStartError", err)
	}

	r, err := dkclient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return "", []string{}, fmt.Errorf("error containerlogs : %v", err)
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

	_, err = stdcopy.StdCopy(&out, &stderr, r)
	if err == context.DeadlineExceeded {
		c := context.Background()
		// timeoutで落ちたときには終了しないため、Container stopでコンテナを終了させる
		stoperr := dkclient.ContainerStop(c, resp.ID, nil)
		if stoperr != nil {
			return "", []string{}, fmt.Errorf("error: %v container timeout and could not stop container", stoperr)
		}
		return "", []string{}, err
	} else if err != nil {
		c := context.Background()
		stoperr := dkclient.ContainerStop(c, resp.ID, nil)
		if stoperr != nil {
			return "", []string{}, fmt.Errorf("error: %v could not run and stop container", stoperr)
		}
		return "", []string{}, fmt.Errorf("error: %v, could not run correctly", err)
	}

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

func oneLiner(args ...string) string {
	var oneline string
	for _, arg := range args {
		oneline = oneline + arg + " "
	}

	return strings.TrimSpace(oneline)
}
