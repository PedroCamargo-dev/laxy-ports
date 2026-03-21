package network

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type dockerContainer struct {
	Names []string     `json:"Names"`
	Ports []dockerPort `json:"Ports"`
}

type dockerPort struct {
	PublicPort uint16 `json:"PublicPort"`
	Type       string `json:"Type"`
}

func fetchDockerPortMap() map[string]string {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}

	resp, err := client.Get("http://localhost/containers/json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var containers []dockerContainer
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil
	}

	portMap := make(map[string]string)
	for _, c := range containers {
		if len(c.Names) == 0 {
			continue
		}
		name := strings.TrimPrefix(c.Names[0], "/")
		for _, p := range c.Ports {
			if p.PublicPort == 0 {
				continue
			}
			key := strconv.Itoa(int(p.PublicPort)) + "|" + strings.ToUpper(p.Type)
			portMap[key] = "[docker] " + name
		}
	}
	return portMap
}
