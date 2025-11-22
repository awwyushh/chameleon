package honeypot

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	Image      string
	TTLMinutes int
}

type SessionInfo struct {
	ContainerID string
	HostPort    int
}

func NewManager(image string, ttl int) *Manager {
	return &Manager{Image: image, TTLMinutes: ttl}
}

func (m *Manager) Spawn(ctx context.Context, payload, srcIP, hpType string) (*SessionInfo, error) {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("chp-%06d", rand.Intn(999999))

	// Spawn: dynamic port on 8080
	args := []string{
		"run", "-d",
		"--name", name,
		"--cap-drop=ALL",
		"--read-only",
		"-p", "0:8080",
		"-e", "CHM_PAYLOAD=" + payload,
		"-e", "CHM_SRC_IP=" + srcIP,
		"-e", "CHM_TYPE=" + hpType,
		m.Image,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker run error: %w", err)
	}
	containerID := strings.TrimSpace(string(out))

	// Inspect host port mapping
	inspectCmd := exec.CommandContext(ctx,
		"docker", "inspect", containerID,
		"--format", "{{json .NetworkSettings.Ports}}",
	)
	inspectOut, err := inspectCmd.Output()
	if err != nil {
		return &SessionInfo{ContainerID: containerID, HostPort: 0}, nil
	}

	// JSON is like:
	// {"8080/tcp":[{"HostIp":"0.0.0.0","HostPort":"32768"}]}
	var parsed map[string][]map[string]string
	if err := json.Unmarshal(inspectOut, &parsed); err != nil {
		return &SessionInfo{ContainerID: containerID, HostPort: 0}, nil
	}

	entries, ok := parsed["8080/tcp"]
	if !ok || len(entries) == 0 {
		return &SessionInfo{ContainerID: containerID, HostPort: 0}, nil
	}

	hostPortStr := entries[0]["HostPort"]
	hostPort, _ := strconv.Atoi(hostPortStr)

	// TTL-based auto removal
	go func(id string, ttl int) {
		time.Sleep(time.Duration(ttl) * time.Minute)
		exec.Command("docker", "rm", "-f", id).Run()
	}(containerID, m.TTLMinutes)

	return &SessionInfo{
		ContainerID: containerID,
		HostPort:    hostPort,
	}, nil
}
