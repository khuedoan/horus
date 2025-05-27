resource "helm_release" "argocd" {
  name             = "argocd"
  namespace        = "argocd"
  create_namespace = true
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  version          = "8.0.6"
  timeout          = 60 * 10

  values = [
    yamlencode({
      global = {
        domain = "argocd.127-0-0-1.nip.io"
      }
      configs = {
        params = {
          "server.insecure"             = true
          "controller.diff.server.side" = true
        }
        cm = {
          "timeout.reconciliation.jitter"         = "60s"
          "resource.ignoreResourceUpdatesEnabled" = true
          "resource.customizations.ignoreResourceUpdates.all" = yamlencode({
            jsonPointers = [
              "/status"
            ]
          })
          "users.anonymous.enabled" = true
        }
        rbac = {
          "policy.default" = "role:admin"
        }
      }
      server = {
        ingress = {
          enabled          = true
          ingressClassName = "nginx"
          tls              = false
        }
      }
      repoServer = {
        hostNetwork = true
      }
      dex = {
        enabled = false
      }
    })
  ]
}
