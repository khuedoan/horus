locals {
  sync_policy = {
    automated = {
      prune    = true
      selfHeal = true
    }
    syncOptions = [
      "CreateNamespace=true",
      "ApplyOutOfSyncOnly=true",
      "ServerSideApply=true"
    ]
  }
}
