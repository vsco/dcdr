
Username = "dcdr_admin"
Namespace = "dcdr"
Storage = "consul"

Consul {
	 Address = "consul:8500"
 }

Watcher {
   OutputPath = "/etc/dcdr/decider.json"
}

Server {
   JsonRoot = "dcdr"
   Endpoint = "/dcdr.json"
}

Git {
   RepoURL = "file:///etc/dcdr/git-backup"
   RepoPath = "/etc/dcdr/audit"
}
