// Package ignite ...
package ignite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type (
	fileVerification struct {
		Hash string `json:"hash,omitempty"`
	}
	fileContents struct {
		Source       string           `json:"source"`
		Verification fileVerification `json:"verification"`
	}
	file struct {
		Filesystem string            `json:"filesystem"`
		Path       string            `json:"path"`
		Contents   fileContents      `json:"contents"`
		Mode       int               `json:"mode"`
		User       map[string]string `json:"user"`
		Group      map[string]string `json:"group"`
	}
	storage struct {
		Filesystem []string `json:"filesystem"`
		Files      []file   `json:"files"`
	}
	systemdDropin struct {
		Name     string `json:"name"`
		Contents string `json:"contents"`
	}
	systemdUnit struct {
		Enable   bool            `json:"enable"`
		Name     string          `json:"name"`
		Contents string          `json:"contents,omitempty"`
		Dropins  []systemdDropin `json:"dropins,omitempty"`
	}
	systemd struct {
		Units    []systemdUnit     `json:"units"`
		Passwd   map[string]string `json:"passwd"`
		Networkd map[string]string `json:"networkd"`
	}
	ignition struct {
		Version string            `json:"version"`
		Config  map[string]string `json:"config"`
	}
	ignitionConfig struct {
		Ignition ignition `json:"ignition"`
		Storage  storage  `json:"storage"`
		Systemd  systemd  `json:"systemd"`
	}
	// binary to fetch on a node
	binary struct {
		// url to fetch binary from, e.g. "https://github.com/hkjn/hkjninfra/releases/download/1.1.7/tserver_x86_64"
		url string
		// checksum of the file, e.g. "sha512-123cec939d7c03c239ee6040185ccb8b74d5d875764479444448ca2ea31d25f364a891363a5850fba2564ce238c7548b3677d713ce69ed7caf421950cd3cd5c6"
		checksum string
		// path on the remote node for the binary, e.g. "/opt/bin/tserver"
		path string
	}
	Version  string
	binaries map[Version][]binary
	// nodeName is the name of a node, e.g. "core"
	nodeName string
	// node is a single instance
	node struct {
		// name is the name of the node
		name nodeName
		// binaries are the files to install on the node
		binaries []binary
		// systemdUnits are the systemd units to use for the node
		systemdUnits []systemdUnit
	}

	nodes       map[nodeName]node
	ProjectName string
	// Project is something that a node should run
	Project struct {
		// name is the name of a project the node should run node, e.g. "hkjninfra"
		Name ProjectName `json:"name"`
		// version is the version of the project that should run on the node, e.g. "1.0.1"
		Version Version `json:"version"`
		// units are the systemd units for the project
		units []systemdUnit
		// binaries are the binaries needed for the project
		binaries []binary
	}
	ProjectConfig struct {
		units []systemdUnit
		files []NodeFile
	}
	ProjectConfigs map[ProjectName]ProjectConfig
	// Projects is a list of projects that a node should run
	Projects []Project
	// NodeConfig is the configuration of a single node
	NodeConfig struct {
		// sshash is the secretservice hash to use
		sshash string
		// projects is all the projects the node should run
		Projects Projects `json:"projects"`
		// arch is the CPU architecture the node runs, e.g. "x86_64"
		Arch string `json:"arch"`
	}
	// NodeConfigs is the configuration of all nodes
	NodeConfigs map[nodeName]NodeConfig
	DropinName  struct {
		Unit, Dropin string
	}
	NodeFile struct {
		Name, ChecksumKey, Path string
		GetUrl                  func(Version) string
	}
	// projectFiles represents the files to include for a project.
	ProjectFiles struct {
		// units are the names of the systemd units for the project
		UnitNames []string
		// dropins are the names of the systemd units and overrides for the project
		DropinNames []DropinName

		// files are the non-systemd files for the project
		Files []NodeFile
	}
)

func (b binary) toFile() file {
	return file{
		Filesystem: "root",
		Path:       b.path,
		Contents: fileContents{
			Source: b.url,
			Verification: fileVerification{
				Hash: fmt.Sprintf("sha512-%s", b.checksum),
			},
		},
		Mode:  493,
		User:  map[string]string{},
		Group: map[string]string{},
	}
}

func newSystemdUnit(unitFile string) (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("units/%s", unitFile))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Enable:   true,
		Name:     unitFile,
		Contents: string(b),
	}, nil
}

func (dn DropinName) Load() (*systemdUnit, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("units/%s", dn.Dropin))
	if err != nil {
		return nil, err
	}
	return &systemdUnit{
		Name: dn.Unit,
		Dropins: []systemdDropin{
			{
				Name:     dn.Dropin,
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
func (n node) Write() error {
	bp := "bootstrap"
	_, err := os.Stat(bp)
	if os.IsNotExist(err) {
		if mkerr := os.Mkdir(bp, 755); mkerr != nil {
			return fmt.Errorf("failed to create dir %q: %v", bp, mkerr)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat %q: %v", bp, err)
	}
	f, err := os.Create(fmt.Sprintf("%s/%s.json", bp, n.name))
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
			Config:  map[string]string{},
		},
		Storage: storage{
			Filesystem: []string{},
			Files: []file{
				file{
					Filesystem: "root",
					Path:       "/etc/coreos/update.conf",
					Contents: fileContents{
						Source:       "data:,GROUP%3Dbeta%0AREBOOT_STRATEGY%3D%22etcd-lock%22",
						Verification: fileVerification{},
					},
					Mode:  420,
					User:  map[string]string{},
					Group: map[string]string{},
				},
			},
		},
		Systemd: systemd{
			Units:    []systemdUnit{},
			Passwd:   map[string]string{},
			Networkd: map[string]string{},
		},
	}
}

// loadFiles returns the non-systemd files for the project.
func (p *Project) loadFiles(arch, sshash string, files []NodeFile) error {
	checksumFile := fmt.Sprintf("checksums/%s_%s.sha512", p.Name, p.Version)
	checksumData, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		return fmt.Errorf("unable to read checksums for %q version %q: %v", p.Name, p.Version, err)
	}
	checksums := map[string]string{}
	for _, line := range strings.Split(string(checksumData), "\n") {
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in checksum file %s: %q", checksumFile, line)
		}
		checksums[parts[1]] = parts[0]
	}
	binaries := make([]binary, len(files), len(files))
	for i, file := range files {
		key := file.ChecksumKey
		if key == "" {
			key = file.Name
		}
		checksum, exists := checksums[key]
		if !exists {
			return fmt.Errorf("missing checksum %q in %s", key, checksumFile)
		}
		binaries[i] = binary{
			url:      file.GetUrl(p.Version),
			checksum: checksum,
			path:     file.Path,
		}
	}
	p.binaries = binaries
	return nil
}

// loadUnits returns the systemd units for the project.
func (pf ProjectFiles) loadUnits() ([]systemdUnit, error) {
	units := []systemdUnit{}
	for _, unitFile := range pf.UnitNames {
		unit, err := newSystemdUnit(unitFile)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}
	for _, d := range pf.DropinNames {
		dropin, err := d.Load()
		if err != nil {
			return nil, err
		}
		units = append(units, *dropin)
	}
	return units, nil
}

// GetURL returns the URL to fetch the checksums for the project.
func (p Project) GetChecksumURL() string {
	return fmt.Sprintf(
		"https://github.com/hkjn/%s/releases/download/%s/SHA512SUMS",
		p.Name,
		p.Version,
	)
}

func (ps Projects) getBinaries() []binary {
	bins := []binary{}
	for _, p := range ps {
		bins = append(bins, p.binaries...)
	}
	return bins
}

func (ps Projects) getUnits() []systemdUnit {
	units := []systemdUnit{}
	for _, p := range ps {
		units = append(units, p.units...)
	}
	return units
}

// createNodes returns nodes created from the configs.
func (nc NodeConfigs) CreateNodes() nodes {
	result := nodes{}
	for name, conf := range nc {
		log.Printf("Generating config for node %q..\n", name)
		result[name] = node{
			name:         name,
			binaries:     conf.Projects.getBinaries(),
			systemdUnits: conf.Projects.getUnits(),
		}
	}
	return result
}

// GetProjectConfigs returns the project configs, given files to load.
func GetProjectConfigs(pfs map[ProjectName]ProjectFiles) (*ProjectConfigs, error) {
	conf := ProjectConfigs{}
	for name, pfs := range pfs {
		units, err := pfs.loadUnits()
		if err != nil {
			return nil, err
		}
		conf[name] = ProjectConfig{
			units: units,
			files: pfs.Files,
		}
	}
	return &conf, nil
}

// load loads the project's systemd units and binaries.
func (p *Project) Load(sshash, arch string, conf ProjectConfig) error {
	p.units = conf.units
	return p.loadFiles(arch, sshash, conf.files)
}

// loadProjects loads the systemd units and binaries for the node config.
func (nc *NodeConfig) loadProjects(projectConf ProjectConfigs) error {
	projects := make([]Project, len(nc.Projects), len(nc.Projects))
	for i, p := range nc.Projects {
		p := p
		pc, exists := projectConf[p.Name]
		if !exists {
			return fmt.Errorf("bug: missing projectConfig for project %q", p.Name)
		}
		if err := p.Load(nc.sshash, nc.Arch, pc); err != nil {
			return err
		}
		projects[i] = p
	}
	nc.Projects = projects
	return nil
}

// load loads the systemd units and binaries for each project in the node configs.
func (nc NodeConfigs) Load(pc ProjectConfigs) error {
	for name, conf := range nc {
		conf := conf
		log.Printf("Loading projects for node %q..\n", name)
		if err := conf.loadProjects(pc); err != nil {
			return err
		}
		nc[name] = conf
	}
	return nil
}

// ReadNodeConfigs returns the node configs, read from disk.
func ReadNodeConfigs() (NodeConfigs, error) {
	conf := NodeConfigs{}
	f, err := os.Open("nodes.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return conf, json.NewDecoder(f).Decode(&conf)
}
