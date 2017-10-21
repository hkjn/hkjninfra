// ignite.go generates Ignite JSON configs.
package main

import (
	"crypto/sha512"
	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"os"
)

const (
	// saltFile is the path to the secretservice salt file.
	saltFile = "/etc/secrets/secretservice/salt"
	// seedFile is the path to the secretservice seed file.
	seedFile = "/etc/secrets/secretservice/seed"
)

type (
	fileVerification struct{
		Hash string `json:"hash"`
	}
	fileContents struct{
		Source string `json:"source"`
		Verification fileVerification `json:"verification"`
	}
	file struct{
		Filesystem string `json:"filesystem"`
		Path string `json:"path"`
		Contents fileContents `json:"contents"`
		Mode int `json:"mode"`
		User map[string]string `json:"user"`
		Group map[string]string `json:"version"`
	}
	storage struct{
		Filesystem []string `json:"filesystem"`
		Files []file `json:"files"`
	}
	systemdUnit struct{
		Enable bool `json:"enable"`
		Name string `json:"name"`
		Contents string `json:"contents"`
	}
	systemd struct{
		Units []systemdUnit `json:"units"`
		Passwd map[string]string `json:"passwd"`
		Networkd map[string]string `json:"networkd"`
	}
	ignition struct{
		Version string `json:"version"`
		Config map[string]string `json:"config"`
	}
	config struct{
		Ignition ignition `json:"ignition"`
		Storage storage `json:"storage"`
		Systemd systemd `json:"systemd"`
	}
	// binary to use on a node
	binary struct{
		// url to fetch binary from, e.g. "https://github.com/hkjn/hkjninfra/releases/download/1.1.7/tserver_x86_64"
		url string
		// checksum of the file, e.g. "sha512-123cec939d7c03c239ee6040185ccb8b74d5d875764479444448ca2ea31d25f364a891363a5850fba2564ce238c7548b3677d713ce69ed7caf421950cd3cd5c6"
		checksum string
		// path on the remote node for the binary, e.g. "/opt/bin/tserver"
		path string
	}
	node struct{
		name string
		binaries []binary
	}
	nodes map[string]node

	nodeConfig map[string]map[string]string
)

func (b binary) toFile() file {
	return file{
		Filesystem: "root",
		Path: b.path,
		Contents: fileContents{
			Source: b.url,
			Verification: fileVerification{
				Hash: b.checksum,
			},
		},
		Mode: 493,
		User: map[string]string{},
		Group: map[string]string{},
	}
}

func (nc nodeConfig) getNodes(sshash string) (nodes, error) {
	result := nodes{}
	for n, versions := range nc {
		arch := "x86_64" // TODO: support other archs.
		bins := []binary{}
		for project, version := range versions {
			newbins, err := getBinaries(project, version, arch, sshash)
			if err != nil {
				return nil, err
			}
			bins = append(bins, newbins...)
		}
		result[n] = node{
			name: n,
			binaries: bins,
		}
	}
	return result, nil
}

func (n node) getFiles() []file {
	result := []file{}
	for _, bin := range n.binaries {
		fmt.Printf("FIXMEH: node %q: %+v\n", n.name, bin)
		result = append(result, bin.toFile())
	}
	return result
}

func (n node) write(bc config) error {
	f, err := os.Create(fmt.Sprintf("bootstrap/%s.json", n.name))
	if err != nil {
		return err
	}
	defer f.Close()
	bc.Storage.Files = append(bc.Storage.Files, n.getFiles()...)
	// TODO: also append bc.Systemd.Units
	log.Printf("Serializing json to JSON..\n")
	bc.serialize(f)
	return nil
}

func newConfig() config {
	return config{
		Ignition: ignition{
			Version: "2.0.0",
			Config: map[string]string{},
		},
		Storage: storage{
			Filesystem: []string{},
			Files: []file{
				file{
					Filesystem: "root",
					Path: "/etc/coreos.update.conf",
					Contents: fileContents{
						Source: "data:,GROUP%3Dbeta%0AREBOOT_STRATEGY%3D%22etcd-lock%22",
						Verification: fileVerification{},
					},
					Mode: 420,
					User: map[string]string{},
					Group: map[string]string{},
				},
			},
		},
		Systemd: systemd{
			Units: []systemdUnit{},
			Passwd: map[string]string{},
			Networkd: map[string]string{},
		},
	}
}

func (c config) serialize(w io.Writer) error {
	return json.NewEncoder(w).Encode(&c)
}

func getBinaries(project, version, arch, sshash string) ([]binary, error) {
	checksum_data, err := ioutil.ReadFile(fmt.Sprintf("checksums/%s_%s.sha512", project, version))
	if err != nil {
		return nil, fmt.Errorf("unable to read checksums for %q version %q: %v", project, version, err)
	}
	checksums := map[string]string{}
	for _, line := range strings.Split(string(checksum_data), "\n") {
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line in checksum file %s.sha512: %q", version, line)
		}
		checksums[parts[1]] = parts[0]
	}

	bins := []binary{}
	if project == "hkjninfra" {
		filename := "gather_facts"
		checksum, exists := checksums[filename]
		if !exists {
			return nil, fmt.Errorf("missing checksum for %s", filename)
		}

		url := fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", project, version, filename)
		bins = append(bins, binary {
			url: url,
			checksum: checksum,
			path: fmt.Sprintf("/opt/bin/%s", filename),
		})

		filename = fmt.Sprintf("tclient_%s", arch)
		checksum, exists = checksums[filename]
		if !exists {
			log.Fatalf("missing checksum for %q", filename)
		}
		url = fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", project, version, filename)
		bins = append(bins, binary {
			url: url,
			checksum: checksum,
			path: fmt.Sprintf("/opt/bin/%s", filename),
		})

		filename = "mon_ca.pem"
		checksum, exists = checksums[filename]
		if !exists {
			log.Fatalf("missing checksum for %q for version %v", filename, version)
		}
		// TODO: filenames for secretservice URLs
		url = fmt.Sprintf("https://admin1.hkjn.me/%s/files/certs/%s", sshash, filename)
		bins = append(bins, binary {
			url: url,
			checksum: checksum,
			path: fmt.Sprintf("/etc/ssl/%s", filename),
		})
	} else if project == "bitcoin" {
	} else if project == "decenter.world" {
		filename := fmt.Sprintf("decenter_world_%s", arch)
		checksum, exists := checksums[filename]
		if !exists {
			log.Fatalf("missing checksum for %q", filename)
		}
		url := fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", project, version, filename)
		bins = append(bins, binary {
			url: url,
			checksum: checksum,
			path: fmt.Sprintf("/opt/bin/%s", filename),
		})

		filename = fmt.Sprintf("decenter_redirector_%s", arch)
		checksum, exists = checksums[filename]
		if !exists {
			log.Fatalf("missing checksum for %q", filename)
		}
		url = fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", project, version, filename)
		bins = append(bins, binary {
			url: url,
			checksum: checksum,
			path: fmt.Sprintf("/opt/bin/%s", filename),
		})
	} else {
		log.Fatalf("bug: unknown release: %q\n", project)
	}
	return bins, nil
}

func getSecretServiceHash() (string, error) {
	salt, err := ioutil.ReadFile(saltFile)
	if err != nil {
		return "", err
	}
	seed, err := ioutil.ReadFile(seedFile)
	if err != nil {
		return "", err
	}
	val := fmt.Sprintf("%s|%s\n", seed, salt)
	digest := sha512.Sum512([]byte(val))
	return fmt.Sprintf("%x", digest), nil
}

func main() {
	nc := nodeConfig{
		"core": map[string]string{
			"hkjninfra": "1.5.0",
			"bitcoin": "0.0.15",
		},
		"decenter_world": map[string]string{
			"hkjninfra": "1.5.0",
			"decenter.world": "1.1.7",
		},
	}
	sshash, err := getSecretServiceHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d character secret service hash\n", len(sshash))

	ns, err := nc.getNodes(sshash)
	if err != nil {
		log.Fatalf("Unable to get node versions: %v\n", err)
	}
	log.Printf("Parsed configs for %d nodes..\n", len(ns))

	for _, n := range ns {
		err := n.write(newConfig())
		if err != nil {
			log.Fatalf("Failed to write node config: %v\n", err)
		}
	}
}
