package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"hkjn.me/hkjninfra/ignite"
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

// download downloads the checksum files.
func download(conf ignite.NodeConfigs) error {
	fetched := map[string]bool{}
	for node, nc := range conf {
		log.Printf("Fetching checksums for node %q..\n", node)
		for _, p := range nc.Projects {
			url := p.GetChecksumURL()
			if !fetched[url] {
				log.Printf("Fetching %q..\n", url)
				if err := fetch(url, p.Name, p.Version); err != nil {
					return err
				}
				fetched[url] = true
			}
		}
	}
	return nil
}

func main() {
	conf, err := ignite.ReadNodeConfigs()
	if err != nil {
		log.Fatalf("Failed to read node config: %v\n", err)
	}
	log.Printf("Read %d node configs..\n", len(conf))
	if err := download(conf); err != nil {
		log.Fatalf("Failed to download checksums: %v\n", err)
	}
}
