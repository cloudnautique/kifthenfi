# kifthenfi (Kubernetes if then) (ALPHA)

## Why?

Sometimes things need to be applied to each namespace when they are created, this provides an initial attempt at a solution.

## What

This project watches namespaces based on selectors and deploys Kubernetes manifests stored in a secret. For instance, if a namespace is created with a specific label, then you could deploy a default set of RBAC or network policies.

# How it works

Installing this Controller creates a new Custom Resource type called NamespaceWatcher:

```yaml
apiVersion: kifthenfi.cloudnautique.com/v1
kind: NamespaceWatcher
metadata:
  name: default-config
  namespace: default
spec:
  namespaceLabelSelector:
    matchLabels:
      myapp: app
    matchExpressions:
      - key: "acorn.io/managed"
        operator: In
        values:
        - "true"
  manifestSecretName: "testing-v1"
```

In the example above, the NamespaceWatcher is configured to watch for namespaces that have the label `myapp=app` and `acorn.io/managed=true` applied to them. If the namespace does have the labels, then **kifthenfi** will apply the manifests found in the secret `testing-v1` in the same namespace.

An example secret with manifests:

```yaml
apiVersion: v1
data:
  configmap: YXBpVmVyc2lvbjogdjEKZGF0YToKICBhbm90aGVyS2V5OiBzb21ldGhpbmcgZWxzZQogIHNvbWV0aGluZzogQSBOZXcgSW50ZXJlc3RpbmcgdGhpbmcKa2luZDogQ29uZmlnTWFwCm1ldGFkYXRhOgogIG5hbWU6IGFjb3JuLXRlc3QtbWFwCg==
  # apiVersion: v1
  # data:
  #   anotherKey: something else
  #   something: A New Interesting thing
  # kind: ConfigMap
  # metadata:
  #   name: acorn-test-map
  network-policy: YXBpVmVyc2lvbjogbmV0d29ya2luZy5rOHMuaW8vdjEKa2luZDogTmV0d29ya1BvbGljeQptZXRhZGF0YToKICBuYW1lOiBkZW55LWVjMi1tZXRhZGF0YS1hY2Nlc3MKc3BlYzoKICBwb2RTZWxlY3Rvcjoge30KICBwb2xpY3lUeXBlczoKICAtIEVncmVzcwogIGVncmVzczoKICAtIHRvOgogICAgLSBpcEJsb2NrOgogICAgICAgIGNpZHI6IDAuMC4wLjAvMAogICAgICAgIGV4Y2VwdDoKICAgICAgICAtIDE2OS4yNTQuMTY5LjI1NC8zMgo=
  # apiVersion: networking.k8s.io/v1
  # kind: NetworkPolicy
  # metadata:
  #   name: deny-ec2-metadata-access
  # spec:
  #   podSelector: {}
  #   policyTypes:
  #   - Egress
  #   egress:
  #   - to:
  #     - ipBlock:
  #       cidr: 0.0.0.0/0
  #       except:
  #       - 169.254.169.254/32
kind: Secret
metadata:
  name: testing-v2
  namespace: default
type: Opaque
```

If a namespace is specified in the manifests, it will be ignored, and set for the namespaces needing that manifest applied.

## Issues / future enhancements

- If two NamespaceWatchers watch the same namespaces with the same secret, they can fight.
- Remove the secret before the NamespaceWatcher resource, things could be orphaned.
- Watch and respond to additional resources.
