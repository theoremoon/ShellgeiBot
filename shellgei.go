package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BotConfigJson struct {
	DockerImage string   `json:"dockerimage"`
	Workdir     string   `json:"workdir"`
	Timeout     string   `json:"timeout"`
	Tags        []string `json:"tags"`
}

type BotConfig struct {
	DockerImage string
	Workdir     string
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

type StdError struct {
	Msg string
}

func (e *StdError) Error() string {
	return e.Msg
}

func RunCmd(cmdstr string, botConfig BotConfig) (string, []string, error) {
	// create shellgei script file
	name, err := RandStr(16)
	if err != nil {
		return "", []string{}, err
	}
	dirname := name + "__images"

	path := filepath.Join(botConfig.Workdir, name)
	file, err := os.Create(path)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %s, directory permission denied?", err)
	}
	defer func() { _ = file.Close() }()
	defer func() { _ = os.RemoveAll(path) }()

	imgdir_path := filepath.Join(botConfig.Workdir, dirname)
	err = os.MkdirAll(imgdir_path, 0777)
	if err != nil {
		return "", []string{}, fmt.Errorf("error: %s, could not create directory", err)
	}
	defer func() { _ = os.RemoveAll(imgdir_path) }()

	_, err = file.WriteString(cmdstr)
	if err != nil {
		return "", []string{}, fmt.Errorf("errors: %s, write failed", err)
	}
	file.Close()

	// execute shellgei in the docker
	cmd := exec.Command("docker", "run", "--rm",
		"--net=none",
		"-m", "10M", "--oom-kill-disable",
		"--pids-limit", "1024",
		"--cap-add", "sys_ptrace",
		"--name", name,
		"-v", path+":/"+name, "-v", imgdir_path+":/images", botConfig.DockerImage,
		"bash", "-c", fmt.Sprintf("chmod +x /%s && sync && ./%s | head -c 1K", name, name))

	defer func() {
		cmd := exec.Command("docker", "stop", name)
		_ = cmd.Run()
		cmd = exec.Command("docker", "rm", name)
		_ = cmd.Run()
	}()

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
		// do nothing
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
	b64imgs := make([]string, 4, 4)

	for i := 0; i < 4; i++ {
		if len(files) <= i {
			break
		}

		img, err := ioutil.ReadFile(filepath.Join(imgdir_path, files[i].Name()))
		if err != nil {
			log.Println(err)
			return out.String(), []string{}, nil
		}
		b64img := base64.StdEncoding.EncodeToString(img)
		b64imgs = append(b64imgs, b64img)
	}

	return out.String(), b64imgs, nil
}
