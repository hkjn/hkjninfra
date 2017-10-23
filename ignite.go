// ignite.go generates Ignite JSON configs.
//
// TODO: Update fetch to generate checksums/ correctly, including for secrets.
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
		Hash string `json:"hash,omitempty"`
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
		Group map[string]string `json:"group"`
	}
	storage struct{
		Filesystem []string `json:"filesystem"`
		Files []file `json:"files"`
	}
	systemdDropin struct{
		Name string `json:"name"`
		Contents string `json:"contents"`
	}
	systemdUnit struct{
		Enable bool `json:"enable"`
		Name string `json:"name"`
		Contents string `json:"contents,omitempty"`
		Dropins []systemdDropin `json:"dropins,omitempty"`
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
	// binary to fetch on a node
	binary struct{
		// url to fetch binary from, e.g. "https://github.com/hkjn/hkjninfra/releases/download/1.1.7/tserver_x86_64"
		url string
		// checksum of the file, e.g. "sha512-123cec939d7c03c239ee6040185ccb8b74d5d875764479444448ca2ea31d25f364a891363a5850fba2564ce238c7548b3677d713ce69ed7caf421950cd3cd5c6"
		checksum string
		// path on the remote node for the binary, e.g. "/opt/bin/tserver"
		path string
	}
	// nodeName is the name of a node, e.g. "core"
	nodeName string
	// node is a single instance
	node struct{
		// name is the name of the node
		name nodeName
		// binaries are the files to install on the node
		binaries []binary
		// systemdUnits are the systemd units to use for the node
		systemdUnits []systemdUnit
	}

	nodes map[nodeName]node
	projectName string
	// project is something that a node should run
	project struct {
		// name is the name of a project the node should run node, e.g. "hkjninfra"
		name projectName
		// version is the version of the project that should run on the node, e.g. "1.0.1"
		version string
	}
	// nodeConfig is the configuration of a single node
	nodeConfig struct{
		// name is the name of the node
		name nodeName
		// projects is all the projects the node should run
		projects []project
		// arch is the CPU architecture the node runs, e.g. "x86_64"
		arch string
	}
	dropinName struct{
		unit, dropin string
	}
	projectConfigs map[projectName]projectConfig
	projectConfig struct{
		units []string
		dropins []dropinName
	}
	// nodeConfigs is the configuration of all nodes
	nodeConfigs map[nodeName]nodeConfig
)

func (b binary) toFile() file {
	return file{
		Filesystem: "root",
		Path: b.path,
		Contents: fileContents{
			Source: b.url,
			Verification: fileVerification{
				Hash: fmt.Sprintf("sha512-%s", b.checksum),
			},
		},
		Mode: 493,
		User: map[string]string{},
		Group: map[string]string{},
	}
}

func newSystemdUnit(unitFile string) (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("units/%s", unitFile))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Enable: true,
		Name: unitFile,
		Contents: string(b),
	}, nil
}

func newSystemdDropin(unitFile, dropinFile string) (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("units/%s", dropinFile))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Name: unitFile,
		Dropins: []systemdDropin{
			{
				Name: dropinFile,
				Contents: string(b),
			},
		},
	}, nil
}

func (n node) getFiles() []file {
	result := []file{}
	for _, bin := range n.binaries {
		result = append(result, bin.toFile())
	}
	return result
}

func (n node) getSystemdUnits() []systemdUnit {
	result := []systemdUnit{}
	for _, unit := range n.systemdUnits {
		result = append(result, unit)
	}
	return result
}

func (n node) String() string {
	return fmt.Sprintf("%q (%d binaries, %d systemd units)", n.name, len(n.binaries), len(n.systemdUnits))
}

func (n node) write(bc config) error {
	f, err := os.Create(fmt.Sprintf("bootstrap/%s.json", n.name))
	if err != nil {
		return err
	}
	defer f.Close()
	bc.Storage.Files = append(bc.Storage.Files, n.getFiles()...)
	bc.Systemd.Units = append(bc.Systemd.Units, n.getSystemdUnits()...)
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
					Path: "/etc/coreos/update.conf",
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

// getBinaries returns the binaries for project.
func (p project) getBinaries(arch, sshash string) ([]binary, error) {
	checksumFile := fmt.Sprintf("checksums/%s_%s.sha512", p.name, p.version)
	checksum_data, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read checksums for %q version %q: %v", p.name, p.version, err)
	}
	checksums := map[string]string{}
	for _, line := range strings.Split(string(checksum_data), "\n") {
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line in checksum file %s: %q", checksumFile, line)
		}
		checksums[parts[1]] = parts[0]
	}

	type nodeFile struct {
		name, checksumKey, path, url string
	}
	mustLoadFiles := func(nodeFiles ...nodeFile) ([]binary, error) {
		binaries := []binary{}
		for _, file := range nodeFiles {
			key := file.checksumKey
			if key == "" {
				key = file.name
			}
			checksum, exists := checksums[key]
			if !exists {
				return nil, fmt.Errorf("missing checksum %q in %s", key, checksumFile)
			}
			binaries = append(binaries, binary{
				url: file.url,
				checksum: checksum,
				path: file.path,
			})
		}
		return binaries, nil
	}

	if p.name == "hkjninfra" {
		return mustLoadFiles(
			nodeFile{
				name: "gather_facts",
				path: "/opt/bin/gather_facts",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", p.name, p.version, "gather_facts"),
			},
			nodeFile{
				name: fmt.Sprintf("tclient_%s", arch),
				path: "/opt/bin/tclient",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", p.name, p.version, "tclient", arch),
			},
			nodeFile{
				name: "mon_ca.pem",
				path: "/etc/ssl/mon_ca.pem",
				url: fmt.Sprintf("https://admin1.hkjn.me/%s/files/certs/%s", sshash, "mon_ca.pem"),
			},
		)
		// TODO: versioning for secretservice URLs
	} else if p.name == "bitcoin" {
		return nil, nil
	} else if p.name == "decenter.world" {
		return mustLoadFiles(
			nodeFile{
				name: fmt.Sprintf("decenter_world_%s", arch),
				path: "/opt/bin/decenter_world",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", p.name, p.version, "decenter_world", arch),
			},
			nodeFile{
				name: fmt.Sprintf("decenter_redirector_%s", arch),
				path: "/opt/bin/decenter_redirector",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", p.name, p.version, "decenter_redirector", arch),
			},
			nodeFile{
				name: "client.pem",
				checksumKey: "decenter.world.pem",
				path: "/etc/ssl/client.pem",
				url: fmt.Sprintf("https://admin1.hkjn.me/%s/files/certs/%s", sshash, "decenter.world.pem"),
			},
			nodeFile{
				name: "client-key.pem",
				checksumKey: "decenter.world-key.pem",
				path: "/etc/ssl/client-key.pem",
				url: fmt.Sprintf("https://admin1.hkjn.me/%s/files/certs/%s", sshash, "decenter.world-key.pem"),
			},
		)
	}
	return nil, fmt.Errorf("bug: unknown project %q", p.name)
}

// getUnits returns the systemd units for the project.
func (p project) getUnits(conf projectConfig) ([]systemdUnit, error) {
	units := []systemdUnit{}
	for _, unitFile := range conf.units {
		unit, err := newSystemdUnit(unitFile)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}
	for _, d := range conf.dropins {
		dropin, err := newSystemdDropin(d.unit, d.dropin)
		if err != nil {
			return nil, err
		}
		units = append(units, *dropin)
	}
	return units, nil
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
	seed = []byte(strings.TrimSpace(string(seed)))
	salt = []byte(strings.TrimSpace(string(salt)))
	val := fmt.Sprintf("%s|%s\n", seed, salt)
	digest := sha512.Sum512([]byte(val))
	return fmt.Sprintf("%x", digest), nil
}

// newNode returns a new node created from the config.
func (nc nodeConfig) newNode(sshash string, pc projectConfigs) (*node, error) {
	bins := []binary{}
	for _, p := range nc.projects {
		newbins, err := p.getBinaries(nc.arch, sshash)
		if err != nil {
			return nil, err
		}
		bins = append(bins, newbins...)
	}
	// TODO: could version the systemd units as well.
	units := []systemdUnit{}
	for _, p := range nc.projects {
		c := pc[p.name]
		newunits, err := p.getUnits(c)
		if err != nil {
			return nil, err
		}
		units = append(units, newunits...)
	}
	return &node{
		name: nc.name,
		binaries: bins,
		systemdUnits: units,
	}, nil
}

// createNodes returns nodes created from the configs.
func (nc nodeConfigs) createNodes(sshash string, pc projectConfigs) (nodes, error) {
	result := nodes{}
	for name, conf := range nc {
		log.Printf("Generating config for node %q..\n", name)
		n, err := conf.newNode(sshash, pc)
		if err != nil {
			return nil, err
		}
		result[name] = *n
	}
	return result, nil
}

func main() {
	pc := projectConfigs{
		"hkjninfra": {
			units: []string{
				"tclient.service",
				"tclient.timer",
			},
		},
		"bitcoin": {
			units: []string{
				"bitcoin.service",
				"containers.mount", // TODO: better name
			},
			dropins: []dropinName{
				{
					unit: "docker.service",
					dropin: "10_override_storage.conf",
				},
			},
		},
		"decenter.world": {
			units: []string{
				"decenter.service",
				"decenter_redirector.service",
				"etc-secrets.mount",
			},
		},
	}
	nc := nodeConfigs{
		"core": nodeConfig{
			name: "core",
			arch: "x86_64",
			projects: []project{
				{
					name: "hkjninfra",
					version: "1.5.0",
				}, {
					name: "bitcoin",
					version: "0.0.15",
				},
			},
		},
		"decenter_world": nodeConfig{
			name: "decenter_world",
			arch: "x86_64",
			projects: []project{
				{
					name: "hkjninfra",
					version: "1.5.0",
				}, {
					name: "decenter.world",
					version: "1.1.7",
				},
			},
		},
	}

	sshash, err := getSecretServiceHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d character secret service hash.\n", len(sshash))

	ns, err := nc.createNodes(sshash, pc)
	if err != nil {
		log.Fatalf("Unable to get node versions: %v\n", err)
	}
	log.Printf("Parsed configs for %d nodes.\n", len(ns))

	for _, n := range ns {
		log.Printf("Writing Ignition config for %v..\n", n)
		err := n.write(newConfig())
		if err != nil {
			log.Fatalf("Failed to write node config: %v\n", err)
		}
	}
}
