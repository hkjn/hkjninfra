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
	systemdDropin struct{
		Name string `json:"name"`
		Contents string `json:"contents"`
	}
	systemdUnit struct{
		Enable bool `json:"enable"`
		Name string `json:"name"`
		Contents string `json:"contents"`
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
	// binary to use on a node
	binary struct{
		// url to fetch binary from, e.g. "https://github.com/hkjn/hkjninfra/releases/download/1.1.7/tserver_x86_64"
		url string
		// checksum of the file, e.g. "sha512-123cec939d7c03c239ee6040185ccb8b74d5d875764479444448ca2ea31d25f364a891363a5850fba2564ce238c7548b3677d713ce69ed7caf421950cd3cd5c6"
		checksum string
		// path on the remote node for the binary, e.g. "/opt/bin/tserver"
		path string
	}
	// node is a single instance
	node struct{
		name string
		// binaries are the files to install on the node
		binaries []binary
		systemdUnits []systemdUnit
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
		// TODO: could version the systemd units as well.
		units := []systemdUnit{}
		for project, _ := range versions {
			newunits, err := getUnits(project)
			if err != nil {
				return nil, err
			}
			units = append(units, newunits...)
		}
		result[n] = node{
			name: n,
			binaries: bins,
			systemdUnits: units,
		}
	}
	return result, nil
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

func (n node) write(bc config) error {
	f, err := os.Create(fmt.Sprintf("bootstrap/%s.json", n.name))
	if err != nil {
		return err
	}
	defer f.Close()
	bc.Storage.Files = append(bc.Storage.Files, n.getFiles()...)
	bc.Systemd.Units = append(bc.Systemd.Units, n.getSystemdUnits()...)
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

// getBinaries returns the binaries for specified version of project.
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

	type nodeFile struct {
		name, path, url string
	}
	mustLoadFiles := func(nodeFiles ...nodeFile) ([]binary, error) {
		binaries := []binary{}
		for _, file := range nodeFiles {
			checksum, exists := checksums[file.name]
			if !exists {
				return nil, fmt.Errorf("missing checksum for %s", file.name)
			}
			binaries = append(binaries, binary{
				url: file.url,
				checksum: checksum,
				path: file.path,
			})
		}
		return binaries, nil
	}

	if project == "hkjninfra" {
		return mustLoadFiles(
			nodeFile{
				name: "gather_facts",
				path: "/opt/bin/gather_facts",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", project, version, "gather_facts"),
			},
			nodeFile{
				name: fmt.Sprintf("tclient_%s", arch),
				path: "/opt/bin/tclient",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", project, version, "tclient", arch),
			},
			nodeFile{
				name: "mon_ca.pem",
				path: "/etc/ssl/mon_ca.pem",
				url: fmt.Sprintf("https://admin1.hkjn.me/%s/files/certs/%s", sshash, "mon_ca.pem"),
			},
		)
		// TODO: versioning for secretservice URLs
	} else if project == "bitcoin" {
		return nil, nil
	} else if project == "decenter.world" {
		// TODO: client.pem, client-key.pem
		return mustLoadFiles(
			nodeFile{
				name: fmt.Sprintf("decenter_world_%s", arch),
				path: "/opt/bin/decenter_world",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", project, version, "decenter_world", arch),
			},
			nodeFile{
				name: fmt.Sprintf("decenter_redirector_%s", arch),
				path: "/opt/bin/decenter_redirector",
				url: fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", project, version, "decenter_redirector", arch),
			},
		)
	}
	return nil, fmt.Errorf("bug: unknown release %q", project)
}

// getUnits returns the systemd units for specified project.
func getUnits(project string) ([]systemdUnit, error) {
	mustLoadUnits := func (unitFiles ...string) ([]systemdUnit, error) {
		units := []systemdUnit{}
		for _, unitFile := range unitFiles {
			unit, err := newSystemdUnit(unitFile)
			if err != nil {
				return nil, err
			}
			units = append(units, *unit)
		}
		return units, nil
	}
	alsoMustLoadDropin := func(
		units []systemdUnit,
		err error, unitFile,
		dropinFile string) ([]systemdUnit, error) {
		if err != nil {
			return nil, err
		}
		dropin, err := newSystemdDropin(unitFile, dropinFile)
		if err != nil {
			return nil, err
		}
		return append(units, *dropin), nil
	}
	if project == "hkjninfra" {
		units, err := mustLoadUnits("tclient.service", "tclient.timer")
		return alsoMustLoadDropin(
			units,
			err,
			"docker.service",
			"10_override_storage.conf",
		)
	} else if project == "bitcoin" {
		return mustLoadUnits(
			"bitcoin.service",
			"containers.mount",
		)
	} else if project == "decenter.world" {
		return mustLoadUnits(
			"decenter.service",
			"decenter_redirector.service",
			"etc-secrets.mount",
		)
	} else {
		return nil, fmt.Errorf("unknown project: %q", project)
	}
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
