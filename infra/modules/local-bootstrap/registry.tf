resource "kubectl_manifest" "registry" {
  server_side_apply = true
  yaml_body = yamlencode({
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name       = "registry"
      namespace  = helm_release.argocd.namespace
      finalizers = ["resources-finalizer.argocd.argoproj.io"]
    }
    spec = {
      project = "default"
      destination = {
        name      = "in-cluster"
        namespace = "registry"
      }
      syncPolicy = local.sync_policy
      source = {
        repoURL        = "http://zotregistry.dev/helm-charts"
        chart          = "zot"
        targetRevision = "0.1.67"
        helm = {
          valuesObject = {
            image = {
              repository = "ghcr.io/project-zot/zot"
            }
            podLabels = {
              "istio.io/dataplane-mode" = "ambient"
            }
            service = {
              type = "NodePort"
              port = 80
              # HACK Use node port for k3s registry mirror.
              # See also ../../cluster/roles/k3s/templates/registries.yaml.j2
              # The range of valid ports is 30000-32767
              nodePort = 30000
            }
            persistence = true
            pvc = {
              create  = true
              storage = "10Gi"
            }
          }
        }
      }
    }
  })
}
