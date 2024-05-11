bundle: {
    apiVersion: "v1alpha1"
    name:       "gitops"
    instances: {
        "flux": {
            module: url: "oci://ghcr.io/stefanprodan/modules/flux-aio"
            namespace: "flux-system"
            values: {
                controllers: {
                    notification: enabled: false
                }
            }
        }
    }
}
