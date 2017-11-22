## [v0.2.2](https://github.com/koudaiii/sltd/releases/tag/v0.2.1)

Fix process #27

## [v0.2.1](https://github.com/koudaiii/sltd/releases/tag/v0.2.1)

Remove debug prints #25

## [v0.2.0](https://github.com/koudaiii/sltd/releases/tag/v0.2.0)

Supported internal elb #23

## [v0.1.0](https://github.com/koudaiii/sltd/releases/tag/v0.1.0)

Update labels Datadog format.

- service_name

```
{
  key: "kube_service",
  value: kubenetes.io/service-name,
}
```

- kubernetescluster

```
{
  key: "kubernetescluster",
  value: cluster name,
}
```

- labels

```
{
  key: "kube_" + key,
  value: value,
}
```

## [v0.0.1](https://github.com/koudaiii/sltd/releases/tag/v0.0.1)

Initial release.
