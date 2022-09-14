# ServiceAccount Authentication between Kubernetes Services

This repo contains a simple demonstration of how to use Kubernetes ServiceAccount tokens to authenticate between services.

There are a number of key factors that have to be considered:

1. The source service must pass a ServiceAccount token to the target service.  Within an HTTP request, this is done by passing it encoded as a string within the http header `X-Client-Id`.
2. For security reasons we want our token to be short lived and uniquely targeted for communication between those 2 service.  To accomplish this we use a `Projected Volume` to add the `expirationSeconds` and `audience` fields to the generated token,  Finally, our code must regularly fetch and refresh the token into a global variable -- this is done with a simple GO-Routine in this example.  Note that the Kubernetes API will look after expiring and rotating the token.
3. The `audience` field added to the token is verified by the receiving service.  To do this it must call the `tokenReview` API in Kubernetes.  In doing this it will pass in a `tokenReviewSpec` that contains a list of valid audiences (e.g. `Audiences: []string{audienceName}`).
4. In order to call the `tokenReview` API a service must have the appropriate permissions granted to it.  To do this we need to add a `ClusterRoleBinding` that grants the `system:auth-delegator` permission to the Pod (or more specifically to the ServiceAccount it is running under).  Note that we don't use a `RoleBinding` here because that would confine the solution to only work in a specified target namespace, but we want to support having the source and target services running in separate namespaces.

> Note that the token provided by the `Projected Volume` in the source service *is* associated with the ServiceAccount that Pod is running under.  What is actually happening here is the ServiceAccount can (and does in this case) have multiple tokens associated with it.  The one statically associated with the ServiceAccount within a specified Secret never expires nor can it have an `audience` field added to it.

This solution is both secure and efficient.  Every Pod in the cluster should already be strongly identified (using ServiceAccount and Namespace) by the Kubernetes control-plane.  Note that this is one of many reasons why we should never use the `default` ServiceAccount for a running Pod.
