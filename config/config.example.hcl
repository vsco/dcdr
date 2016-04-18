
Username = "dcdr admin"
Namespace = "dcdr"
Storage = "consul"

//Etcd {
//  Endpoints = ["http://127.0.0.1:2379"]
//}

Watcher {
  OutputPath = "/etc/dcdr/decider.json"
}

Server {
  JsonRoot = "dcdr"
  Endpoint = "/dcdr.json"
  Host = "0.0.0.0:9000"
}

Git {
  //RepoPath = "/etc/dcdr/audit"
}

Stats {
  Namespace = "dcdr"
  Host = "127.0.0.1"
  Port = 8125
}

