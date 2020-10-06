package rules

import (
	"fmt"
	"testing"

	"github.com/grupozap/aegir/internal/pkg/utils"
	"gotest.tools/assert"
)

var test_pod = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": "2019-06-18T16:06:58Z",
        "generateName": "authnetes-7b4c684b96-",
        "labels": {
            "app": "authnetes",
            "pod-template-hash": "3607240652",
            "process": "web",
            "product": "platform"
        },
        "name": "authnetes-7b4c684b96-ntnsl",
        "namespace": "platform",
        "ownerReferences": [
            {
                "apiVersion": "apps/v1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "authnetes-7b4c684b96",
                "uid": "bcb8eb05-71b0-11e9-aa37-1242bd31f10a"
            }
        ],
        "resourceVersion": "252081177",
        "selfLink": "/api/v1/namespaces/platform/pods/authnetes-7b4c684b96-ntnsl",
        "uid": "19da054e-91e3-11e9-99c1-12f8d8095f20"
    },
    "spec": {
        "containers": [
            {
                "args": [
                    "serve",
                    "--google-client-id",
                    "$(CLIENT_ID)",
                    "--google-client-secret",
                    "$(CLIENT_SECRET)",
                    "--github-client-id",
                    "$(GITHUB_CLIENT_ID)",
                    "--github-client-secret",
                    "$(GITHUB_CLIENT_SECRET)",
                    "--github-redirect-url",
                    "$(GITHUB_REDIRECT_URL)",
                    "--github-organization",
                    "$(GITHUB_ORGANIZATION_NAME)",
                    "--domain",
                    "$(DOMAIN)",
                    "--google-redirect-url",
                    "$(REDIRECT_URL)",
                    "--secret-key",
                    "$(SECRET_KEY)",
                    "--port",
                    "8000",
                    "--cluster-endpoint",
                    "$(CLUSTER_ENDPOINT)",
                    "--kubernetes-version",
                    "$(KUBERNETES_VERSION)"
                ],
                "env": [
                    {
                        "name": "GITHUB_REDIRECT_URL",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_REDIRECT_URL",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "GITHUB_CLIENT_ID",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_CLIENT_ID",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "GITHUB_CLIENT_SECRET",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_CLIENT_SECRET",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "GITHUB_ORGANIZATION_NAME",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_ORGANIZATION_NAME",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "CLUSTER_ENDPOINT",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "CLUSTER_ENDPOINT",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "DOMAIN",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "DOMAIN",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "KUBERNETES_VERSION",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "KUBERNETES_VERSION",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "CLIENT_ID",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "CLIENT_ID",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "CLIENT_SECRET",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "CLIENT_SECRET",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "REDIRECT_URL",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "REDIRECT_URL",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "SECRET_KEY",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "SECRET_KEY",
                                "name": "authnetes"
                            }
                        }
                    }
                ],
                "image": "vivareal/authnetes:master",
                "imagePullPolicy": "Always",
                "livenessProbe": {
                    "failureThreshold": 3,
                    "httpGet": {
                        "path": "/",
                        "port": 8000,
                        "scheme": "HTTP"
                    },
                    "initialDelaySeconds": 30,
                    "periodSeconds": 10,
                    "successThreshold": 1,
                    "timeoutSeconds": 30
                },
                "name": "authnetes",
                "ports": [
                    {
                        "containerPort": 8000,
                        "protocol": "TCP"
                    },
                    {
                        "containerPort": 9090,
                        "protocol": "TCP"
                    }
                ],
                "readinessProbe": {
                    "failureThreshold": 3,
                    "httpGet": {
                        "path": "/",
                        "port": 8000,
                        "scheme": "HTTP"
                    },
                    "initialDelaySeconds": 5,
                    "periodSeconds": 10,
                    "successThreshold": 1,
                    "timeoutSeconds": 10
                },
                "resources": {
                    "limits": {
                        "cpu": "100m",
                        "memory": "50Mi"
                    },
                    "requests": {
                        "cpu": "100m",
                        "memory": "50Mi"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-605vx",
                        "readOnly": true
                    }
                ]
            },
            {
                "args": [
                    "serve",
                    "--google-client-id",
                    "$(CLIENT_ID)",
                    "--google-client-secret",
                    "$(CLIENT_SECRET)",
                    "--github-client-id",
                    "$(GITHUB_CLIENT_ID)",
                    "--github-client-secret",
                    "$(GITHUB_CLIENT_SECRET)",
                    "--github-redirect-url",
                    "$(GITHUB_REDIRECT_URL)",
                    "--github-organization",
                    "$(GITHUB_ORGANIZATION_NAME)",
                    "--domain",
                    "$(DOMAIN)",
                    "--google-redirect-url",
                    "$(REDIRECT_URL)",
                    "--secret-key",
                    "$(SECRET_KEY)",
                    "--port",
                    "8000",
                    "--cluster-endpoint",
                    "$(CLUSTER_ENDPOINT)",
                    "--kubernetes-version",
                    "$(KUBERNETES_VERSION)"
                ],
                "env": [
                    {
                        "name": "GITHUB_REDIRECT_URL",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_REDIRECT_URL",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "GITHUB_CLIENT_ID",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_CLIENT_ID",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "GITHUB_CLIENT_SECRET",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_CLIENT_SECRET",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "GITHUB_ORGANIZATION_NAME",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "GITHUB_ORGANIZATION_NAME",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "CLUSTER_ENDPOINT",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "CLUSTER_ENDPOINT",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "DOMAIN",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "DOMAIN",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "KUBERNETES_VERSION",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "KUBERNETES_VERSION",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "CLIENT_ID",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "CLIENT_ID",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "CLIENT_SECRET",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "CLIENT_SECRET",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "REDIRECT_URL",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "REDIRECT_URL",
                                "name": "authnetes"
                            }
                        }
                    },
                    {
                        "name": "SECRET_KEY",
                        "valueFrom": {
                            "secretKeyRef": {
                                "key": "SECRET_KEY",
                                "name": "authnetes"
                            }
                        }
                    }
                ],
                "image": "vivareal/authnetes:master",
                "imagePullPolicy": "Always",
                "livenessProbe": {
                    "failureThreshold": 3,
                    "httpGet": {
                        "path": "/",
                        "port": 8000,
                        "scheme": "HTTP"
                    },
                    "initialDelaySeconds": 30,
                    "periodSeconds": 10,
                    "successThreshold": 1,
                    "timeoutSeconds": 30
                },
                "name": "authnetes",
                "ports": [
                    {
                        "containerPort": 8000,
                        "protocol": "TCP"
                    },
                    {
                        "containerPort": 9090,
                        "protocol": "TCP"
                    }
                ],
                "readinessProbe": {
                    "failureThreshold": 3,
                    "httpGet": {
                        "path": "/",
                        "port": 8000,
                        "scheme": "HTTP"
                    },
                    "initialDelaySeconds": 5,
                    "periodSeconds": 10,
                    "successThreshold": 1,
                    "timeoutSeconds": 10
                },
                "resources": {
                    "limits": {
                        "cpu": "100m",
                        "memory": "50Mi"
                    },
                    "requests": {
                        "cpu": "100m",
                        "memory": "50Mi"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-605vx",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "ip-10-160-105-143.ec2.internal",
        "priority": 0,
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
            {
                "effect": "NoSchedule",
                "key": "node-role.kubernetes.io/master"
            },
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/not-ready",
                "operator": "Exists",
                "tolerationSeconds": 300
            },
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/unreachable",
                "operator": "Exists",
                "tolerationSeconds": 300
            }
        ],
        "volumes": [
            {
                "name": "default-token-605vx",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-605vx"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2019-06-18T16:06:58Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2019-06-18T16:07:34Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": null,
                "status": "True",
                "type": "ContainersReady"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2019-06-18T16:06:58Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://6d3e8a6c7cb09ee7f34b3ba215ab8c1c6dd33713b0fc9b227519645fd5717776",
                "image": "vivareal/authnetes:master",
                "imageID": "docker-pullable://vivareal/authnetes@sha256:4b763a566cc3a8d8287ca4f7208923f0913c2cb8ff035dfc2846e8c11ad1dfc0",
                "lastState": {},
                "name": "authnetes",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2019-06-18T16:07:20Z"
                    }
                }
            }
        ],
        "hostIP": "10.160.105.143",
        "phase": "Running",
        "podIP": "172.20.175.198",
        "qosClass": "Guaranteed",
        "startTime": "2019-06-18T16:06:58Z"
    }
}
`

func TestGetJSONObjectByPath(t *testing.T) {
	jsonObject := GetJSONObjectByPath(test_pod, "metadata.labels")
	expected := `{
            "app": "authnetes",
            "pod-template-hash": "3607240652",
            "process": "web",
            "product": "platform"
        }`
	if jsonObject[0].String() != expected {
		t.Errorf("expected '%s' but got '%s'", expected, jsonObject)
	}
}

func TestGetJSONObjectByPathTwoArrays(t *testing.T) {
	jsonObject := GetJSONObjectByPath(test_pod, "spec.containers.#.ports.#.protocol")
	expected := `TCP`
	for _, value := range jsonObject {
		if value.Value() != expected {
			t.Errorf("expected '%s' but got '%s'", expected, jsonObject)
		}
	}
}

func TestGetJSONObjectByPathOneArray(t *testing.T) {
	jsonObject := GetJSONObjectByPath(test_pod, "spec.containers.#.name")
	expected := `authnetes`
	if jsonObject[0].Value() != expected {
		t.Errorf("expected '%s' but got '%s'", expected, jsonObject)
	}
}

func TestGetViolationsFoundViolations(t *testing.T) {
	rulesloaded := RulesLoader("testing_rules.yaml")
	violation := &utils.Violation{
		Description: "Container name must be skull",
		JSONPath:    "spec.containers.#.name",
		Object:      map[string]interface{}{"name": string("authnetes")},
		Message:     fmt.Sprint("validation error"),
	}
	for _, rule := range rulesloaded.Rules {
		v := rule.RulesDefinitions[0].GetViolations(test_pod)
		assert.DeepEqual(t, v[0], violation)
	}
}

func TestGetViolationsFieldIsOptional(t *testing.T) {
	rulesloaded := RulesLoader("testing_rules.yaml")
	for _, rule := range rulesloaded.Rules {
		v := rule.RulesDefinitions[2].GetViolations(test_pod)
		if len(v) != 0 {
			t.Error("Expecting 0 violations")
		}
	}
}

func TestGetViolationsFieldIsRequired(t *testing.T) {
	rulesloaded := RulesLoader("testing_rules.yaml")
	violation := &utils.Violation{
		Description: "This field is NOT optional and should fail if it does not exists",
		JSONPath:    "spec.containers.#.THIS_PATH_DOES_NOT_EXIST_BUT_MUST",
		Message:     fmt.Sprint("Field: spec.containers.#.THIS_PATH_DOES_NOT_EXIST_BUT_MUST is required"),
	}
	for _, rule := range rulesloaded.Rules {
		v := rule.RulesDefinitions[3].GetViolations(test_pod)
		if len(v) < 1 {
			t.Error("Expecting 1 violation")
		}
		assert.DeepEqual(t, v[0], violation)
	}
}
