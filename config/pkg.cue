package config

import (
    corev1 "k8s.io/api/core/v1"
    rbacv1 "k8s.io/api/rbac/v1"
    appsv1 "k8s.io/api/apps/v1"
    extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

#Namespace: corev1.#Namespace & {
    apiVersion: "v1"
    kind: "Namespace"
}



#Name: {
    metadata: name: "kink"
}


DeploymentNamespace: #Namespace & {
    #Name
}

#NamespacedName: {
    #Name
    metadata: namespace: DeploymentNamespace.metadata.name

}

#ServiceAccount: corev1.#ServiceAccount & {
    apiVersion: "v1"
    kind: "ServiceAccount"
}

DeployServiceAccount: #ServiceAccount & {
    #NamespacedName
}

#RBACGroup: "rbac.authorization.k8s.io"

#RBACV1: {
   apiVersion: "\( #RBACGroup )/v1"
}

#ClusterRole: rbacv1.#ClusterRole & {
   #RBACV1
   kind: "ClusterRole"
}

DeployClusterRole: #ClusterRole & {
   #Name
   rules: [
       {
           resources: ["*"]
           apiGroups: ["*"]
           verbs: ["*"]
       },
       {
           nonResourceURLs: ["*"]
           verbs: ["*"]
       }
   ]
}

#ClusterRoleBinding: rbacv1.#ClusterRoleBinding & {
   #RBACV1
   kind: "ClusterRoleBinding"
}


DeployClusterRoleBinding: #ClusterRoleBinding & {
    #Name
    roleRef: {
        apiGroup: "\( #RBACGroup )"
        kind: "ClusterRole"
        name: DeployClusterRole.metadata.name
    }
    subjects: [
        {
            kind: "ServiceAccount"
            name: DeployServiceAccount.metadata.name
            namespace: DeployServiceAccount.metadata.namespace
        }
    ]
}

#Deployment: appsv1.#Deployment & {
    apiVersion: "apps/v1"
    kind: "Deployment"
}

DeployDeployment: #Deployment & {
    #NamespacedName
    spec: {
      selector: {
        matchLabels: {
            "control-plane": "controller-manager"
        }
      }
      replicas: 1
      template: {
        metadata: {
          labels: {
            "control-plane": "controller-manager"
          }
        }
        spec: {
          containers: [
            {
                command: [
                    "/manager"
                ]
                args: [
                    "--enable-leader-election"
                ]
                image: "controller:latest"
                name: "manager"
                resources: {
                  limits: {
                    cpu: "100m"
                    memory: "30Mi"
                  }
                  requests: {
                    cpu: "100m"
                    memory: "20Mi"
                  }
                }
            }
          ]
          terminationGracePeriodSeconds: 10
      }
    }
  }
}

#CustomResourceDefinition: extv1.#CustomResourceDefinition & {
    apiVersion: "apiextensions.k8s.io/v1"
    kind:       "CustomResourceDefinition"
}

DeployCustomResourceDefinition: #CustomResourceDefinition & {
    metadata: {
        annotations: "controller-gen.kubebuilder.io/version": "v0.4.1"
        creationTimestamp: null
        name:              "clusters.kink.x-k8s.io"
    }
    spec: {
        group: "kink.x-k8s.io"
        names: {
            kind:     "Cluster"
            listKind: "ClusterList"
            plural:   "clusters"
            singular: "cluster"
        }
        scope: "Namespaced"
        versions: [{
            name: "v1alpha1"
            schema: openAPIV3Schema: {
                description: "Cluster is the Schema for the clusters API"
                properties: {
                    apiVersion: {
                        description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources"

                        type: "string"
                    }
                    kind: {
                        description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"

                        type: "string"
                    }
                    metadata: type: "object"
                    spec: {
                        description: "ClusterSpec defines the desired state of Cluster."
                        type:        "object"
                    }
                    status: {
                        description: "ClusterStatus defines the observed state of Cluster"
                        properties: serviceRef: {
                            description: "ObjectReference contains enough information to let you inspect or modify the referred object. --- New uses of this type are discouraged because of difficulty describing its usage when embedded in APIs.  1. Ignored fields.  It includes many fields which are not generally honored.  For instance, ResourceVersion and FieldPath are both very rarely valid in actual usage.  2. Invalid usage help.  It is impossible to add specific help for individual usage.  In most embedded usages, there are particular     restrictions like, \"must refer only to types A and B\" or \"UID not honored\" or \"name must be restricted\".     Those cannot be well described when embedded.  3. Inconsistent validation.  Because the usages are different, the validation rules are different by usage, which makes it hard for users to predict what will happen.  4. The fields are both imprecise and overly precise.  Kind is not a precise mapping to a URL. This can produce ambiguity     during interpretation and require a REST mapping.  In most cases, the dependency is on the group,resource tuple     and the version of the actual struct is irrelevant.  5. We cannot easily change it.  Because this type is embedded in many locations, updates to this type     will affect numerous schemas.  Don't make new APIs embed an underspecified API type they do not control. Instead of using this type, create a locally provided and used type that is well-focused on your reference. For example, ServiceReferences for admission registration: https://github.com/kubernetes/api/blob/release-1.17/admissionregistration/v1/types.go#L533 ."

                            properties: {
                                apiVersion: {
                                    description: "API version of the referent."
                                    type:        "string"
                                }
                                fieldPath: {
                                    description: "If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For example, if the object reference is to a container within a pod, this would take on a value like: \"spec.containers{name}\" (where \"name\" refers to the name of the container that triggered the event) or if no container name is specified \"spec.containers[2]\" (container with index 2 in this pod). This syntax is chosen only to have some well-defined way of referencing a part of an object. TODO: this design is not final and this field is subject to change in the future."

                                    type: "string"
                                }
                                kind: {
                                    description: "Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"
                                    type:        "string"
                                }
                                name: {
                                    description: "Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names"
                                    type:        "string"
                                }
                                namespace: {
                                    description: "Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/"
                                    type:        "string"
                                }
                                resourceVersion: {
                                    description: "Specific resourceVersion to which this reference is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency"

                                    type: "string"
                                }
                                uid: {
                                    description: "UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids"
                                    type:        "string"
                                }
                            }
                            type: "object"
                        }
                        required: [
                            "serviceRef",
                        ]
                        type: "object"
                    }
                }
                type: "object"
            }
            served:  true
            storage: true
        }]
    }
}
