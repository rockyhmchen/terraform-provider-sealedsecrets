provider "sealedsecrets" {
    // optional
    kubectl_bin = "/usr/local/bin/kubectl"

    // optional
    kubeseal_bin = "/usr/local/bin/kubeseal"
}

resource "sealedsecrets_secret" "docker_pull" {
  name      = "docker_pull"
  namespace = kubernetes_namespace.example_ns.metadata.0.name
  type = "kubernetes.io/dockerconfigjson"

  secret_source = "${join("/", [var.keys_directory, "docker", "pull", ".dockerconfigjson"])}"
  sealed_secret_source = "${join("/", [var.secret_directory, "docker-pull.yaml"])}"
  certificate = "${join("/", [var.tmp_directory, "sealed-secrets", "tls.crt"])}"

  depends_on = [kubernetes_namespace.example_ns, var.sealed_secrets_controller_id]
}