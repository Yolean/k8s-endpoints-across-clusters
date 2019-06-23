NOTE: This is WIP, not even PoC yet.

# Synchronize Kubernetes Services Between Clusters

Amazingly we couldn't find any existing solution to non-federated cross-cluster DNS name resolution to pod IP, so here we go ...

Like [k8s-endpoints-sync-controller](https://github.com/vmware/k8s-endpoints-sync-controller)
this controller:
 - Assumes that pods can reach other clusters' pods using the `podIP`s.
 - Watches `Service`s with the purpose to replicate their [Endpoints](https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors) in other clusters.
 - Synchronizes these endpoints so that cross-cluster access is transparent to containers.

However it stays even further away from [Federation](https://github.com/kubernetes-sigs/federation-v2) by avoding cross-cluster Kubernetes API access. The reason being to reduce ops risk, which sort of is the reason we're going multi-cluster at all.

Scope:
 - Test with GKE clusters created using [--enable-ip-alias](https://cloud.google.com/kubernetes-engine/docs/how-to/alias-ips),
   but design for any topology that meets the Pod IP requirement.
 - Handle failures where one cluster is gone and comes back with all `podIP`s changed.
 - Use GCP PubSub as a reference out-of-cluster messaging means.
