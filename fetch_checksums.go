// fetch_checksums.go is a tool to read nodes.json
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"hkjn.me/hkjninfra/ignite"
	"hkjn.me/hkjninfra/secretservice"
)

// checkClose closes specified closer and sets err to the result.
func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}

// fetch downloads the checksums from specified url.
func fetch(url string, project ignite.ProjectName, version ignite.Version) (err error) {
	// TODO: Also need to handle secrets, like decenter.world.pem for "decenter.world"..
	// fetch from secret service directly?
	if project == ignite.ProjectName("bitcoin") {
		// TODO: Instead of special-casing "core" (bitcoin) project, which has
		// no checksums since there's no binaries to download, maybe start
		// checksumming / versioning systemd unit (.service, .mount) and
		// dropins (.conf) within the project?
		log.Printf("Skipping bitcoin, no binaries to download..\n")
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer checkClose(resp.Body, &err)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from GET %q, want 200 OK, got %s", url, resp.Status)
	}
	filename := fmt.Sprintf("checksums/%s_%s.sha512", project, version)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer checkClose(f, &err)
	log.Printf("Writing to %q..\n", filename)
	_, err = io.Copy(f, resp.Body)
	return err
}

// downloadChecksums downloads the checksum files.
func downloadChecksums(conf ignite.ConfigJSON, sshash string) error {
	fetched := map[string]bool{}
	for node, nc := range conf.Nodes {
		log.Printf("Fetching checksums for node %q..\n", node)
		for _, pv := range nc.ProjectVersions {
			url := p.GetChecksumURL()
			if !fetched[url] {
				log.Printf("Fetching %q..\n", url)
				if err := fetch(url, p.Name, p.Version); err != nil {
					return err
				}
				fetched[url] = true
			}
			// TODO: using pv.Name and pv.Version, get the SecretURLs from conf.Projects
			secretURLs := p.GetSecretURLs(sshash, secretservice.BaseDomain)
			log.Printf("FIXMEH: secret urls for %q: %v\n", p.Name, secretURLs)
			// FIXMEH: NodeConfig.Load(ProjectConfigs) is what sets p.secretFiles, so without that
			// there's no way to go from nodes.json -> secret service URLs.. maybe refactor out
			// projects.json to hold []SecretFile?
			for _, url := range secretURLs {
				if !fetched[url] {
					log.Printf("Fetching secret %q..\n", url)
					if err := fetch(url, p.Name, p.Version); err != nil {
						return err
					}
					fetched[url] = true
				}
			}
		}
	}
	return nil
}

func main() {
	conf, err := ignite.ReadConfig()
	//conf, err := ignite.ReadNodeConfigs()
	if err != nil {
		log.Fatalf("Failed to read node config: %v\n", err)
	}
	for k, c := range conf.Projects {
		log.Printf("FIXMEH: project %q: %+v\n", k, c)
	}
	for k, c := range conf.Nodes {
		log.Printf("FIXMEH: node %q: %+v\n", k, c)
	}

	sshash, err := secretservice.GetHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d node configs..\n", len(conf.Nodes))
	if err := downloadChecksums(*conf, sshash); err != nil {
		log.Fatalf("Failed to download checksums: %v\n", err)
	}
}
