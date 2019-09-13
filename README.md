# Ægir

![Aegir](https://upload.wikimedia.org/wikipedia/commons/9/98/%C3%86gir%2C_ruler_of_the_ocean.jpg)


### Aegir is a simple and generic webhook admission controller for Kubernetes.

It allows you to write custom rules for your cluster resources. If your rule is violated, Aegir will not allow the resource to be created and will display a message on the terminal, optionally it can send a notification in a Slack channel.

Aegir uses [LIVR](http://livr-spec.org) to validate the rules. Any LIVR rule is supported.

Example of rules:

```yaml
rules:
- name: required_labels
# All Namespaces
  namespace: "*"
  resource_type: "Deployment"
  rules_definitions:
  - field: "metadata.labels"
    livr_rule:
      description: "Labels should have an app label"
      rule:
      # THE LAST FIELD SHOULD ALSO BE DECLARED HERE
        labels:
        # NOW THE RULE ITSELF
          nested_object:
            app: required
            version:
            - required
            - positive_integer
  - field: "spec.template.spec.containers.#.port.#.protocol"
    livr_rule:
      description: "Containers protocol should be http or https"
      rule:
        name:
          one_of: ['https', 'http']
  slack_notification_channel: "#some_team_channel"
  # Another rule
- name: container_user_could_not_be_root
  namespace: "*"
  resource_type: "Deployment"
  rules_definitions:
  - field: "spec.template.spec.securityContext.runAsUser"
    livr_rule:
      description: "Only non-root users are allowed"
      rule:
        runAsUser:
          number_between: [1, 1000]
  slack_notification_channel: "#some_team_channel"
  ```

### Usage

```shell
A generic admission controller to validate Kubernetes resources using LIVR rules.

Usage:
  aegir [command]

Available Commands:
  help        Help about any command
  server      Runs Aegir's admission controller.

Flags:
  -h, --help      help for aegir
      --version   version for aegir

Use "aegir [command] --help" for more information about a command.
```

### Running Aegir on your Kubernetes cluster

Create a `Deployment` and a `Service`

```shell
kubectl apply -f examples/aegir-deployment.yaml
kubectl apply -f examples/aegir-service.yaml
```

To make aegir be able validating cluster resources create a `ValidatingWebhookConfiguration` like this:

```yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: aegir-webhook
webhooks:
  - name: aegir.example.svc
    sideEffects: NoneOnDryRun
    clientConfig:
      service:
        name: aegir
        namespace: example
        # This path should be /admission
        path: "/admission"
      caBundle: base64 encoded CA certificate
    rules:
      - apiGroups:
        - apps
        - extensions
        apiVersions:
        - v1
        - v1beta1
        operations:
        - UPDATE
        - CREATE
        resources:
        - deployments
        - services
        - ingresses
  ```

  ### Important note
  `sideEffects` should be set to `NoneOnDryRun` so `Aegir` can validate the rules when you run `--server-dry-run` with `kubectl`. This is useful
  running CI/CD pipelines or trying to validate the configuration of the object before persisting it on ETCD


### Skipping some namespaces

If you have defined a rule with `*` this rule will run against all namespaces. Sometimes is useful to skip some namespaces, like `kube-system`, `istio-system` and etc.
To do this you can set the environment variable `SKIP_NAMESPACES=namespace1,namespace2,namespace3`, and these namespaces will be skipped at rule evaluation.

### TLS certificates

The Kubernetes API needs to trust the certificate to connect to Aegir's webhook.
Use the `genkey.sh` script to generate self-signed certificates, passing some directory and Commom Name as parameters.

```shell
$ ./genkey.sh dir/ foobar.com
```

This will create some certificate files inside directory `dir`.

Encode the file `dir/ca.crt` into base64

```shell
$ base64 dir/ca.crt
```

Use the output to fill the field `caBundle` in the `ValidationWebhookConfiguration`

Use the files `webhook-server-tls.crt` and `webhook-server-tls.key` passing the flags:
```
aegir server \
--rules-file=rules.yaml \
--tls-cert-file=dir/webhook-server-tls.crt \
--tls-key-file=dir/webhook-server-tls.key
```

And that's it!

### Limitations and Warnings
Aegir is pretty new and have some limitations for now:
- Can't validate if a field is part of a Kubernetes Object.
- There is no parsing or validation for the configuration file format.
- Only a few unit tests aiming the main part of the validation rules.

All this problems will be addressed in the future.

Aegir is under development, changes and improvements will come.

In the future Aegir should be converted into CRD's.

Feedbacks and PR's are welcome.
