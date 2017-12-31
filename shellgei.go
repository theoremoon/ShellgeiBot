package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func RunCmd(cmdstr string, botConfig BotConfig) (string, error) {
	// create shellgei script file
	name, err := RandStr(16)
	if err != nil {
		return "", err
	}

	path := filepath.Join(botConfig.Workdir, name)
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("error: %s, directory permission denied?", err)
	}
	defer func() { _ = file.Close() }()
	defer func() { _ = os.RemoveAll(path) }()

	_, err = file.WriteString(cmdstr)
	if err != nil {
		return "", fmt.Errorf("errors: %s, write failed", err)
	}

	// execute shellgei in the docker
	cmd := exec.Command("docker", "run", "--net=none", "--rm", "--name", name, "-v", path+":/"+name, botConfig.DockerImage, "bash", "/"+name)
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
		if err != nil {
			return "", &StdError{fmt.Sprintf("err: %s -- execution error? %s ", err.Error(), stderr.String())}
		}
	}
	return out.String(), nil
}
