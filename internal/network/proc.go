package network

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"laxy-ports/internal/models"
)

const tcpListenState = 0x0A

func ScanPorts() []models.PortEntry {
	inodePIDMap := buildInodePIDMap()
	dockerPortMap := fetchDockerPortMap()

	var entries []models.PortEntry
	entries = append(entries, readNetFile("/proc/net/tcp", "TCP", true, inodePIDMap)...)
	entries = append(entries, readNetFile("/proc/net/tcp6", "TCP", true, inodePIDMap)...)
	entries = append(entries, readNetFile("/proc/net/udp", "UDP", false, inodePIDMap)...)
	entries = append(entries, readNetFile("/proc/net/udp6", "UDP", false, inodePIDMap)...)

	unique := deduplicateEntries(entries)
	enrichWithDocker(unique, dockerPortMap)

	sort.Slice(unique, func(i, j int) bool {
		if unique[i].Port != unique[j].Port {
			return unique[i].Port < unique[j].Port
		}
		return unique[i].Protocol < unique[j].Protocol
	})

	return unique
}

func enrichWithDocker(entries []models.PortEntry, dockerPortMap map[string]string) {
	for i := range entries {
		key := strconv.Itoa(int(entries[i].Port)) + "|" + entries[i].Protocol
		if containerName, ok := dockerPortMap[key]; ok {
			entries[i].Process = containerName
		}
	}
}

func buildInodePIDMap() map[uint64]int {
	inodes := make(map[uint64]int)

	procDir, err := os.Open("/proc")
	if err != nil {
		return inodes
	}
	pidNames, _ := procDir.Readdirnames(-1)
	procDir.Close()

	for _, pidStr := range pidNames {
		if _, err := strconv.Atoi(pidStr); err != nil {
			continue
		}

		fdDirPath := "/proc/" + pidStr + "/fd"
		fdDir, err := os.Open(fdDirPath)
		if err != nil {
			continue
		}
		fdNames, _ := fdDir.Readdirnames(-1)
		fdDir.Close()

		for _, fdName := range fdNames {
			linkTarget, err := os.Readlink(fdDirPath + "/" + fdName)
			if err != nil {
				continue
			}

			inode, ok := parseSocketInode(linkTarget)
			if !ok {
				continue
			}

			if _, alreadyMapped := inodes[inode]; !alreadyMapped {
				pid, _ := strconv.Atoi(pidStr)
				inodes[inode] = pid
			}
		}
	}

	return inodes
}

func parseSocketInode(link string) (uint64, bool) {
	if !strings.HasPrefix(link, "socket:[") || len(link) < 10 {
		return 0, false
	}
	inode, err := strconv.ParseUint(link[8:len(link)-1], 10, 64)
	if err != nil {
		return 0, false
	}
	return inode, true
}

func readNetFile(path, protocol string, tcpListenOnly bool, inodePIDMap map[uint64]int) []models.PortEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var entries []models.PortEntry

	for lineIdx, line := range lines {
		if lineIdx == 0 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		if tcpListenOnly {
			stateVal, err := strconv.ParseUint(fields[3], 16, 8)
			if err != nil || stateVal != tcpListenState {
				continue
			}
		}

		localAddr := fields[1]
		colonPos := strings.LastIndex(localAddr, ":")
		if colonPos < 0 {
			continue
		}

		portVal, err := strconv.ParseUint(localAddr[colonPos+1:], 16, 16)
		if err != nil || portVal == 0 {
			continue
		}

		inodeVal, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			continue
		}

		pid := inodePIDMap[inodeVal]
		entries = append(entries, models.PortEntry{
			Port:     uint16(portVal),
			Protocol: protocol,
			PID:      pid,
			Process:  resolveProcessName(pid),
		})
	}

	return entries
}

func resolveProcessName(pid int) string {
	if pid == 0 {
		return "kernel"
	}
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/comm")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func deduplicateEntries(entries []models.PortEntry) []models.PortEntry {
	seen := make(map[string]bool)
	result := make([]models.PortEntry, 0, len(entries))

	for _, e := range entries {
		key := strconv.Itoa(int(e.Port)) + "|" + e.Protocol + "|" + strconv.Itoa(e.PID)
		if !seen[key] {
			seen[key] = true
			result = append(result, e)
		}
	}

	return result
}
