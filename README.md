# Golang Multi-Service Example

Simple Golang source code, based upon original Golang example for OpenShift (https://github.com/sclorg/golang-ex.git).

> The code is for playground purposes only. No intent for 'quality' Golang code.


## Install

Note: you may need a pull secret to obtain the 'go-toolset' base image.

```shell
git clone https://github.com/avoigtma/golang-example.git 

oc import-image -n openshift rhel8/go-toolset:1.16 --from=registry.redhat.io/rhel8/go-toolset:1.16 --confirm

oc new-project a-goex
oc apply -n a-goex -R -f OpenShift

```

## Tracing

The Deployments used the `sidecar.istio.io/inject: "true"` annotation to inject the Envoy sidecar on OpenShift with OpenShift ServiceMesh installed.

Jaeger can be used to gather tracing information.

As the default configuration does not apply 'clock skew' adjustments to the traces, thus resulting in unadjusted display of the trace spans, add to the Jaeger configuration (`oc edit jaeger jaeger -n istio-system`) the query configuration parameter `max-clock-skew-adjustment`, e.g. `Jaeger.spec.allInOne.options.query.max-clock-skew-adjustment: 20s`.
