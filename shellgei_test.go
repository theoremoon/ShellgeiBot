package main

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	dockerImage     = "alpine-bash"
	dockerImagePath = "./testdata/alpine-bash.tar"
)

var dockerclient, _ = client.NewEnvClient()

func setupDocker(t *testing.T) func(t *testing.T) {
	ctx := context.Background()
	imgdata, err := os.Open(dockerImagePath)
	if err != nil {
		t.Fatalf("docker load unexpected error, %v", err)
	}
	loadResponce, err := dockerclient.ImageLoad(ctx, imgdata, false)
	if err != nil {
		t.Fatalf("docker load unexpected error, %v", err)
	}
	defer loadResponce.Body.Close()
	body, _ := ioutil.ReadAll(loadResponce.Body)
	t.Logf(string(body))
	return func(t *testing.T) {
		ctx := context.Background()
		_, err := dockerclient.ImageRemove(ctx, dockerImage, types.ImageRemoveOptions{})
		if err != nil {
			t.Errorf("docker rmi unexpected error, %v", err)
		}
	}
}

func TestRunCmd(t *testing.T) {
	teardown := setupDocker(t)
	defer teardown(t)

	t.Run("simple-command", func(t *testing.T) {
		out, imgs, err := runCmd("echo hello", nil, botConfig{
			DockerImage: dockerImage,
			Workdir:     "/tmp",
			Memory:      "100M",
			MediaSize:   250,
			Timeout:     time.Duration(20 * time.Second),
			Tags:        []string{},
		})
		if err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		if expected := "hello\n"; out != expected {
			t.Errorf("got %q, expected %v", out, expected)
		}
		if len(imgs) != 0 {
			t.Errorf("got %+v, expected empty", imgs)
		}
	})

	t.Run("using /images dir", func(t *testing.T) {
		out, imgs, err := runCmd("echo hello; echo goodbye > /images/a.png", nil, botConfig{
			DockerImage: dockerImage,
			Workdir:     "/tmp",
			Memory:      "100M",
			MediaSize:   250,
			Timeout:     time.Duration(20 * time.Second),
			Tags:        []string{},
		})
		if err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		if expected := "hello\n"; out != expected {
			t.Errorf("got %q, expected %v", out, expected)
		}
		if len(imgs) != 1 {
			t.Errorf("got %+v, expected length 1", imgs)
		} else if got, expected := imgs[0], base64.StdEncoding.EncodeToString([]byte("goodbye\n")); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	})

	t.Run("with media urls", func(t *testing.T) {
		urls := []string{"a.png", "b.png"}
		mux := http.NewServeMux()
		for _, url := range urls {
			url := url
			mux.HandleFunc(
				"/"+url,
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("hello-" + url))
				},
			)
		}
		ts := httptest.NewServer(mux)
		defer ts.Close()
		for i := range urls {
			urls[i] = ts.URL + "/" + urls[i]
		}

		out, imgs, err := runCmd("cat /media/0 /media/1", urls, botConfig{
			DockerImage: dockerImage,
			Workdir:     "/tmp",
			Memory:      "100M",
			MediaSize:   250,
			Timeout:     time.Duration(20 * time.Second),
			Tags:        []string{},
		})
		if err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		if expected := "hello-a.pnghello-b.png"; out != expected {
			t.Errorf("got %q, expected %v", out, expected)
		}
		if len(imgs) != 0 {
			t.Errorf("got %+v, expected empty", imgs)
		}
	})

	t.Run("check timeout", func(t *testing.T) {
		out, imgs, err := runCmd("sleep 2", nil, botConfig{
			DockerImage: dockerImage,
			Workdir:     "/tmp",
			Memory:      "100M",
			MediaSize:   250,
			Timeout:     time.Duration(1 * time.Second),
			Tags:        []string{},
		})
		if err != context.DeadlineExceeded {
			t.Fatalf("unexpected error, %v", err)
		}
		if len(imgs) != 0 && out == "" {
			t.Errorf("got %+v, expected empty", imgs)
		}
	})
}
