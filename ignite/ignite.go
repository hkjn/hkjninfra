// Package ignite deals with Ignite JSON configs.
package ignite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
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
		// checksum of the file, e.g. "sha512-123[...]"
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
	SecretFiles []SecretFile
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
		// secretFilenames are the secrets needed for the project
		secretFiles SecretFiles
		// secretServiceDomain is the base domain for the secret service
		secretServiceDomain string
	}
	ProjectVersion struct {
		// name is the name of a project the node should run node, e.g. "hkjninfra"
		Name ProjectName `json:"name"`
		// version is the version of the project that should run on the node, e.g. "1.0.1"
		Version Version `json:"version"`
	}
	ProjectConfig struct {
		units       []systemdUnit
		files       []nodeFile
		secretFiles SecretFiles
	}
	ProjectConfigs struct {
		secretServiceDomain string
		configs             map[ProjectName]ProjectConfig
	}
	// Projects is a list of projects that a node should run
	Projects []Project
	// NodeConfig is the configuration of a single node
	NodeConfig struct {
		// sshash is the secretservice hash to use
		sshash string
		// projectVersions is the names of all the projects the node should run
		ProjectVersions []ProjectVersion `json:"projects"`
		// arch is the CPU architecture the node runs, e.g. "x86_64"
		Arch string `json:"arch"`
	}
	// NodeConfigs is the configuration of all nodes
	NodeConfigs map[nodeName]NodeConfig

	nodeFile struct {
		Path        string `json:"path"`
		Name        string `json:"name"`
		ChecksumKey string `json:"checksum_key"`
	}
	projectJSON struct {
		Units   []string     `json:"units"`
		Dropins []DropinName `json:"dropins"`
		Files   []nodeFile   `json:"files"`
		Secrets SecretFiles  `json:"secrets"`
	}
	projectsJSON map[ProjectName]projectJSON
	ConfigJSON   struct { // TODO: better name
		Projects projectsJSON `json:"projects"`
		Nodes    NodeConfigs  `json:"nodes"`
	}
	DropinName struct {
		Unit, Dropin string
	}
	// FIXMEH: rename to NodeFile
	SecretFile struct {
		Name        string `json:"name"`
		ChecksumKey string `json:"checksum_key"`
		Path        string `json:"path"`
	}
	// projectFiles represents the files to include for a project.
	ProjectFiles struct {
		// Units are the names of the systemd units for the project
		UnitNames []string
		// Dropins are the names of the systemd units and overrides for the project
		DropinNames []DropinName
		// Files are the non-systemd files for the project
		Files []nodeFile
		// SecretFiles are the secret files for the project
		SecretFiles []SecretFile
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

// Load returns the systemd units.
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
			u, _ := user.Current()
			return fmt.Errorf("failed to create dir %q as %s:%s: %v", bp, u.Uid, u.Gid, mkerr)
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

// newProjectConfig returns the project config.
func newProjectConfig(conf projectJSON) (*ProjectConfig, error) {
	units := []systemdUnit{}
	for _, unitFile := range conf.Units {
		unit, err := newSystemdUnit(unitFile)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}
	for _, d := range conf.Dropins {
		dropin, err := d.Load()
		if err != nil {
			return nil, err
		}
		units = append(units, *dropin)
	}
	return &ProjectConfig{
		units:       units,
		files:       conf.Files,
		secretFiles: conf.Secrets,
	}, nil
}

// loadUnits returns the systemd units for the project.
// FIXMEH: remove
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

func getChecksums(name ProjectName, version Version) (map[string]string, error) {
	checksumFile := fmt.Sprintf("checksums/%s_%s.sha512", name, version)
	checksumData, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read checksums for %q version %q: %v", name, version, err)
	}
	checksums := map[string]string{}
	for _, line := range strings.Split(string(checksumData), "\n") {
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line in checksum file %s: %q", checksumFile, line)
		}
		checksums[parts[1]] = parts[0]
	}
	return checksums, nil
}

// loadFiles returns the non-systemd files for the project.
// FIXMEH: remove
func (p *Project) loadFiles(arch, sshash string, files []nodeFile) error {
	checksums, err := getChecksums(p.Name, p.Version)
	if err != nil {
		return err
	}

	binaries := make([]binary, len(files), len(files))
	for i, file := range files {
		key := file.ChecksumKey
		if key == "" {
			key = file.Name
		}
		checksum, exists := checksums[key]
		if !exists {
			return fmt.Errorf("missing checksum %q in %s", key, checksums)
		}
		binaries[i] = binary{
			url:      fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", p.Name, p.Version, file.Name),
			checksum: checksum,
			path:     file.Path,
		}
	}
	p.binaries = binaries
	return nil
}

// load loads the project's systemd units and binaries.
func (p *Project) load(sshash, arch string, conf ProjectConfigs) error {
	pc, exists := conf.configs[p.Name]
	if !exists {
		return fmt.Errorf("bug: missing project configs for project %q", p.Name)
	}
	p.secretServiceDomain = conf.secretServiceDomain
	p.secretFiles = pc.secretFiles
	fmt.Printf("FIXMEH: in Project.Load(), secretFiles=%v\n", p.secretFiles)
	p.units = pc.units
	return p.loadFiles(arch, sshash, pc.files)
}

// GetChecksumURL returns the URL to fetch the checksums for the project.
func (p Project) GetChecksumURL() string {
	return fmt.Sprintf(
		"https://github.com/hkjn/%s/releases/download/%s/SHA512SUMS",
		p.Name,
		p.Version,
	)
}

// GetSecretURLs returns the URLs for any secrets in the project.
func (p Project) GetSecretURLs(sshash, secretServiceDomain string) []string {
	fmt.Printf("FIXMEH: secretFiles=%v\n", p.secretFiles)
	results := make([]string, len(p.secretFiles), len(p.secretFiles))
	for i, secret := range p.secretFiles {
		results[i] = fmt.Sprintf("https://%s/%s/files/%s/%s/certs/%s", secretServiceDomain, sshash, p.Name, p.Version, secret.Name)
	}
	return results
}

func (ps Projects) getUnits() []systemdUnit {
	units := []systemdUnit{}
	for _, p := range ps {
		units = append(units, p.units...)
	}
	return units
}

func (ps Projects) String() string {
	names := make([]string, len(ps), len(ps))
	for i, p := range ps {
		names[i] = string(p.Name)
	}
	return strings.Join(names, ", ")
}

func (sf SecretFiles) String() string {
	if len(sf) == 0 {
		return "[empty SecretFiles]"
	}
	files := make([]string, len(sf), len(sf))
	for i, f := range sf {
		files[i] = fmt.Sprintf("SecretFile{Name: %s, ChecksumKey: %s, Path: %s}}", f.Name, f.ChecksumKey, f.Path)
	}
	return strings.Join(files, ", ")
}

func (nc NodeConfig) String() string {
	return fmt.Sprintf(fmt.Sprintf("NodeConfig{Arch: %s}", nc.Arch))
	// return fmt.Sprintf(fmt.Sprintf("NodeConfig{Arch: %s, Projects: %s}", nc.Arch, nc.Projects.String()))
}

func (p projectJSON) String() string {
	return fmt.Sprintf("projectJSON{Units: %s, Secrets: %s,..}",
		strings.Join(p.Units, ", "),
		p.Secrets.String(),
	)
}

// getBinaries returns the binaries for the specific project.
// FIXMEH
func (conf projectsJSON) getBinaries(pversions []ProjectVersion) ([]binary, error) {
	result := []binary{}
	for _, p := range pversions {
		fmt.Printf("FIXMEH: should load binaries for project %q version %v\n", p.Name, p.Version)
		pc, exists := conf[p.Name]
		if !exists {
			return nil, fmt.Errorf("bug: no such project %q", p.Name)
		}

		bins, err := pc.getBinaries(p.Name, p.Version)
		if err != nil {
			return nil, err
		}
		result = append(result, bins...)
	}
	return result, nil
}

// getBinaries returns the binaries.
func (conf projectJSON) getBinaries(projectName ProjectName, projectVersion Version) ([]binary, error) {
	// TODO: Find better place to load checksums to avoid loading same ones over
	// and over.
	checksums, err := getChecksums(projectName, projectVersion)
	if err != nil {
		return nil, err
	}

	result := []binary{}
	for _, file := range conf.Files {
		key := file.ChecksumKey
		if key == "" {
			key = file.Name
		}
		checksum, exists := checksums[key]
		if !exists {
			return nil, fmt.Errorf("missing checksum for key %q; all checksums %v", key, checksums)
		}
		result = append(result, binary{
			url: fmt.Sprintf(
				"https://github.com/hkjn/%s/releases/download/%s/%s",
				projectName,
				projectVersion,
				file.Name,
			),
			checksum: checksum,
			path:     file.Path,
		})
	}
	return result, nil
}

// getUnits returns the systemd units for the specific projects.
func (conf projectsJSON) getUnits(pversions []ProjectVersion) ([]systemdUnit, error) {
	result := []systemdUnit{}
	for _, p := range pversions {
		pc, exists := conf[p.Name]
		if !exists {
			return nil, fmt.Errorf("bug: no such project %q", p.Name)
		}
		pconf, err := newProjectConfig(pc)
		if err != nil {
			return nil, err
		}
		// p.Version
		result = append(result, pconf.units...)
	}
	return result, nil
}

func (conf ConfigJSON) CreateNodes() (nodes, error) {
	// createNodes returns nodes created from the configs.
	result := nodes{}
	for name, nc := range conf.Nodes {
		log.Printf("Generating config for node %q..\n", name)
		bins, err := conf.Projects.getBinaries(nc.ProjectVersions)
		if err != nil {
			return nil, err
		}
		units, err := conf.Projects.getUnits(nc.ProjectVersions)
		if err != nil {
			return nil, err
		}
		result[name] = node{
			name:         name,
			binaries:     bins,
			systemdUnits: units,
		}
	}
	return result, nil
}

// GetProjectConfigs returns the project configs, given files to load.
func GetProjectConfigs(secretServiceDomain string, pfconf map[ProjectName]ProjectFiles) (*ProjectConfigs, error) {
	conf := map[ProjectName]ProjectConfig{}
	for name, pfs := range pfconf {
		units, err := pfs.loadUnits()
		if err != nil {
			return nil, err
		}
		conf[name] = ProjectConfig{
			units:       units,
			files:       pfs.Files,
			secretFiles: pfs.SecretFiles,
		}
	}
	return &ProjectConfigs{
		secretServiceDomain: secretServiceDomain,
		configs:             conf,
	}, nil
}

// ReadConfig returns the node/project configs, read from disk.
func ReadConfig() (*ConfigJSON, error) {
	conf := ConfigJSON{}
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		return nil, err
	}
	pconfs := map[ProjectName]ProjectConfig{}
	for name, pconf := range conf.Projects {
		pc, err := newProjectConfig(pconf)
		if err != nil {
			return nil, err
		}
		pconfs[name] = *pc
	}
	return &conf, nil
}
