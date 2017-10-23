// ignite.go generates Ignite JSON configs.
//
// TODO: Update fetch to generate checksums/ correctly, including for secrets.
// TODO: could version the systemd units as well.
package main

import (
	"crypto/sha512"
	"fmt"
	"encoding/json"
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
	ignitionConfig struct{
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
	version string
	binaries map[version][]binary
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
		version version
		// units are the systemd units for the project
		units []systemdUnit
		// binaries are the binaries needed for the project
		binaries []binary
	}
	projectConfig map[projectName][]systemdUnit
	// projects is a list of projects that a node should run
	projects []project
	// nodeConfig is the configuration of a single node
	nodeConfig struct{
		// name is the name of the node
		name nodeName
		// projects is all the projects the node should run
		projects projects
		// arch is the CPU architecture the node runs, e.g. "x86_64"
		arch string
	}
	dropinName struct{
		unit, dropin string
	}
	// projectFiles represents the systemd files to include for a project.
	projectFiles struct{
		// units are the names of the systemd units for the project
		unitNames []string
		// dropins are the names of the systemd units and overrides for the project
		dropinNames []dropinName
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

func (dn dropinName) load() (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("units/%s", dn.dropin))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Name: dn.unit,
		Dropins: []systemdDropin{
			{
				Name: dn.dropin,
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

// write writes the Ignition config to disk.
func (n node) write() error {
	f, err := os.Create(fmt.Sprintf("bootstrap/%s.json", n.name))
	if err != nil {
		return err
	}
	defer f.Close()

	conf := newIgnitionConfig()
	conf.Storage.Files = append(conf.Storage.Files, n.getFiles()...)
	conf.Systemd.Units = append(conf.Systemd.Units, n.getSystemdUnits()...)
	return json.NewEncoder(f).Encode(&conf)
}

func newIgnitionConfig() ignitionConfig {
	return ignitionConfig{
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

// getBinaries returns the binaries for project.
func getBinaries(name projectName, v version, arch, sshash string) ([]binary, error) {
	checksumFile := fmt.Sprintf("checksums/%s_%s.sha512", name, v)
	checksum_data, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read checksums for %q version %q: %v", name, v, err)
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

	if name == "hkjninfra" {
		return mustLoadFiles(
			nodeFile{
				name: "gather_facts",
				path: "/opt/bin/gather_facts",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", name, v, "gather_facts"),
			},
			nodeFile{
				name: fmt.Sprintf("tclient_%s", arch),
				path: "/opt/bin/tclient",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", name, v, "tclient", arch),
			},
			nodeFile{
				name: "mon_ca.pem",
				path: "/etc/ssl/mon_ca.pem",
				url: fmt.Sprintf("https://admin1.hkjn.me/%s/files/certs/%s", sshash, "mon_ca.pem"),
			},
		)
		// TODO: versioning for secretservice URLs
	} else if name == "bitcoin" {
		return nil, nil
	} else if name == "decenter.world" {
		return mustLoadFiles(
			nodeFile{
				name: fmt.Sprintf("decenter_world_%s", arch),
				path: "/opt/bin/decenter_world",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", name, v, "decenter_world", arch),
			},
			nodeFile{
				name: fmt.Sprintf("decenter_redirector_%s", arch),
				path: "/opt/bin/decenter_redirector",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", name, v, "decenter_redirector", arch),
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
	return nil, fmt.Errorf("bug: unknown project %q", name)
}

// loadUnits returns the systemd units for the project.
func (pf projectFiles) loadUnits() ([]systemdUnit, error) {
	units := []systemdUnit{}
	for _, unitFile := range pf.unitNames {
		unit, err := newSystemdUnit(unitFile)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}
	for _, d := range pf.dropinNames {
		dropin, err := d.load()
		if err != nil {
			return nil, err
		}
		units = append(units, *dropin)
	}
	return units, nil
}

// getSecretServiceHash returns the secret service hash read from files.
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

func (ps projects) getBinaries(sshash, arch string) []binary {
	bins := []binary{}
	for _, p := range ps {
		b, err := getBinaries(p.name, p.version, arch, sshash)
		if err != nil {
			log.Fatalf("Failed to load binaries for %q: %v\n", p.name, err)
		}
		bins = append(bins, b...)
	}
	return bins
}

func (ps projects) getUnits() []systemdUnit {
	units := []systemdUnit{}
	for _, p := range ps {
		units = append(units, p.units...)
	}
	return units
}

// createNodes returns nodes created from the configs.
func (nc nodeConfigs) createNodes(sshash string) nodes {
	result := nodes{}
	for name, conf := range nc {
		log.Printf("Generating config for node %q..\n", name)
		result[name] = node{
			name: name,
			binaries: conf.projects.getBinaries(sshash, conf.arch),
			systemdUnits: conf.projects.getUnits(),
		}
	}
	return result
}

// getProjectConfigs returns the project configs, given files to load.
func getProjectConfigs(pf map[projectName]projectFiles) projectConfig {
	conf := map[projectName][]systemdUnit{}
	for name, files := range pf {
		units, err := files.loadUnits()
		if err != nil {
			log.Fatalf("Failed to load systemd units: %v\n", err)
		}
		conf[name] = units
	}
	return conf
}


func main() {
	pc := getProjectConfigs(map[projectName]projectFiles{
		"hkjninfra": {
			unitNames: []string{
				"tclient.service",
				"tclient.timer",
			},
		},
		"bitcoin": {
			unitNames: []string{
				"bitcoin.service",
				"containers.mount", // TODO: better name
			},
			dropinNames: []dropinName{
				{
					unit: "docker.service",
					dropin: "10_override_storage.conf",
				},
			},
		},
		"decenter.world": {
			unitNames: []string{
				"decenter.service",
				"decenter_redirector.service",
				"etc-secrets.mount",
			},
		},
	})
	nc := nodeConfigs{
		"core": nodeConfig{
			name: "core",
			arch: "x86_64",
			projects: []project{
				{
					name: "hkjninfra",
					version: "1.5.0",
					units: pc["hkjninfra"],
				}, {
					name: "bitcoin",
					version: "0.0.15",
					units: pc["bitcoin"],
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
					units: pc["hkjninfra"],
				}, {
					name: "decenter.world",
					version: "1.1.7",
					units: pc["decenter.world"],
				},
			},
		},
	}

	sshash, err := getSecretServiceHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d character secret service hash.\n", len(sshash))
	for _, n := range nc.createNodes(sshash) {
		log.Printf("Writing Ignition config for %v..\n", n)
		err := n.write()
		if err != nil {
			log.Fatalf("Failed to write node config: %v\n", err)
		}
	}
}
