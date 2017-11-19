// ignite.go generates Ignite JSON configs.
//
// TODO: could version the systemd units as well.
package main

import (
	"fmt"
	"log"

	"hkjn.me/hkjninfra/ignite"
	"hkjn.me/hkjninfra/secretservice"
)

func main() {
	sshash, err := secretservice.GetHash()
	if err != nil {
		log.Fatalf("Unable to fetch secret service hash: %v\n", err)
	}
	log.Printf("Read %d character secret service hash.\n", len(sshash))
	arch := "x86_64" // TODO
	pc, err := ignite.GetProjectConfigs(map[ignite.ProjectName]ignite.ProjectFiles{
		"hkjninfra": {
			UnitNames: []string{
				"tclient.service",
				"tclient.timer",
			},
			Files: []ignite.NodeFile{
				{
					Name: "gather_facts",
					Path: "/opt/bin/gather_facts",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s", "hkjninfra", v, "gather_facts")
					},
				}, {
					Name: fmt.Sprintf("tclient_%s", arch),
					Path: "/opt/bin/tclient",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", "hkjninfra", v, "tclient", arch)
					},
				}, {
					Name: "mon_ca.pem",
					Path: "/etc/ssl/mon_ca.pem",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://admin1.hkjn.me/%s/files/%s/%s/certs/%s", sshash, "hkjninfra", v, "mon_ca.pem")
					},
				},
			},
		},
		"bitcoin": {
			UnitNames: []string{
				"bitcoin.service",
				"containers.mount", // TODO: better name
			},
			DropinNames: []ignite.DropinName{
				{
					Unit:   "docker.service",
					Dropin: "10_override_storage.conf",
				},
			},
		},
		"decenter.world": {
			UnitNames: []string{
				"decenter.service",
				"decenter_redirector.service",
				"etc-secrets.mount",
			},
			Files: []ignite.NodeFile{
				{
					Name: fmt.Sprintf("decenter_world_%s", arch),
					Path: "/opt/bin/decenter_world",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", "decenter.world", v, "decenter_world", arch)
					},
				}, {
					Name: fmt.Sprintf("decenter_redirector_%s", arch),
					Path: "/opt/bin/decenter_redirector",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://github.com/hkjn/%s/releases/download/%s/%s_%s", "decenter.world", v, "decenter_redirector", arch)
					},
				}, {
					Name:        "client.pem",
					ChecksumKey: "decenter.world.pem",
					Path:        "/etc/ssl/client.pem",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://admin1.hkjn.me/%s/files/%s/%s/certs/%s", sshash, "decenter.world", v, "decenter.world.pem")
					},
				}, {
					Name:        "client-key.pem",
					ChecksumKey: "decenter.world-key.pem",
					Path:        "/etc/ssl/client-key.pem",
					GetUrl: func(v ignite.Version) string {
						return fmt.Sprintf("https://admin1.hkjn.me/%s/files/%s/%s/certs/%s", sshash, "decenter.world", v, "decenter.world-key.pem")
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create project configs: %v\n", err)
	}

	nc, err := ignite.ReadNodeConfigs()
	if err != nil {
		log.Fatalf("Failed to read node configs: %v\n", err)
	}
	if err := nc.Load(*pc); err != nil {
		log.Fatalf("Failed to create node configs: %v\n", err)
	}

	for _, n := range nc.CreateNodes() {
		log.Printf("Writing Ignition config for %v..\n", n)
		err := n.Write()
		if err != nil {
			log.Fatalf("Failed to write node config: %v\n", err)
		}
	}
}
