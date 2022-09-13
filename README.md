# Kubernetes Auth Example

This repo contains some very simple services, an `api` and a `data-store` that demonstrate how to authenticate calls from one service to another using Kubernetes Service-Account tokens.

## The Exercise

To get this working you must perform a number of steps.  The services are written in Go; if you understand Go then feel free to examine the source code, though this is not necessary to get the solution working.
The `Optional` tasks can be skipped if it is taking too long to figure out.

1. These services will be deployed separately so you must create a Dockerfile for each.
2. [Optional] Think of optimizing both the size a security of the resulting container images; i.e. use best practices.
3. Create Helm charts for each of the services.  Required configurations for each service are detailed below (this may also affect how you configure your Dockerfile).
4. [Optional] Use Helm best practices and be sure to consider security here too.
5. The build: use a combination of docs, and/or scripts, and/or CI system integration to automate the Docker builds and possibly even installing the Helm Apps into a Kubernetes cluster.
6. [Optional] How would you expose the `api` Service outside of the Kubernetes cluster.  (This may require further changes or additions to the associated Helm chart.)
7. [**Bonus**/Optional] Add terraform, or use some other configuration tool to setup a test Kubernetes cluster.  This can be a managed cluster in a common cloud provider (we use GCP and DigitalOcean).

## Expected Results

 - The Docker image builds should be easily reproducible. The Dockerfiles and resulting Docker images will be inspected to ensure they are working as expected.
 - The helm charts should be easily configured and used once the Docker images have been built and pushed.  Note that due to configuration requirements the `data-store` Service should be deployed first.
 - Installing the Services with the Helm charts should work *even if* the 2 services are installed to different Namespaces in the same cluster.
 - Use Git.  We use GitHub, but your fork of the repo can be hosted elsewhere though it should be sharable with us.
 - Feel free to make notes about issues that were encountered/resolved.

## Configurations
 
 The following are critical configurations needed to be passed to each Service.  This should not be the complete set of configurations passed to the Helm charts -- other incidental settings may be needed to fully configure the Kubernetes manifests that get rendered from the charts.
 
 > Note that both of the Services are using HTTP as the communication protocol.
 
### api Service
 
#### Environment Variables

- `LISTEN_ADDR` -- The port this Service is exposed on in Kubernetes.  This should be a string of the form `:<number>`.
- `DATA_STORE_CONNECTION` -- The connection string use by `api` to connect to the the `data-store`.
 
#### Other Settings

- `AUDIENCE_NAME` -- The `audience` field to be added to the token used to authenticate with the `data-store` Service.
 
 ---

### data-store Service

#### Environment Variables

- `LISTEN_ADDR` -- The port this Service is exposed on in Kubernetes.  This should be a string of the form `:<number>`.
- `AUDIENCE_NAME` -- The expected `audience` field on the incoming auth token (passed from `api`).

---
 
 > Note (Hint #1): The `api` Service must use a projected volume to set both the number of seconds before the token expires as well as set an audience on the token.  Also note (*important*) the `audience` setting must match between both services.
 
 > Note (Hint #2): The `data-store` Service needs to be able to validate ("authenticate") any tokens passed to it using the Kubernetes API.

 > Note (Hint #3): Each of the Services implements global variables for `version` (app version), `commit` (commit hash being built), and `date` (built timestamp) which are displayed on start-up.  These should be set at build time (this requires knowledge of `go build` options).

---

## Questions?

If you get stuck don't hesitate to contact us and ask questions.  This is an advanced problem so you may not get it all done, but that is ok.  Contact info will be provided in the same email when this exercise is given to you.
