k3d_cluster_name    = ["k3d-jenkins"]
k3d_cluster_port    = 6550
k3d_cluster_ip      = "127.0.0.1"
k3d_cluster_lb_port = 80
k3d_host_lb_port    = 8085
server_count        = 1     # Increase to >=3 for HA cluster
agent_count         = 2
k3s_version         = "v1.20.2-k3s1"
