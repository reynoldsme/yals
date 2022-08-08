# Yet Another Link Shortener (YALS) Design considerations and assumptions

* It is assumed that the service may eventually be extended to scale horizontally.
* It is assumed that the service may eventually be extended to provide some form of analytics to the users generating shortened links.
* It is assumed that users of the service may not want the existence of a given shortened link to be discoverable.

* Requesting a link be shortened that has already been shortened by the service will result in a *new* shortened link pointing to the same URL.
  * This enables multiple users to generate links to the same resource and track link consumption independently from one another.
  * This enables a single user to generate multiple links to the same resource and track link consumption by multiple users or groups.
  * This prevents the service from leaking information.
    * If user A generates a shortened link for url X and any subsequent request to shorten url X produces different behavior (for example, a 409 "Conflict" instead of a 200 "Ok") this allows user B to determine that someone else already shortened the link, which may not be desirable.
  * Not having to check if a URL has already been shortened greatly reduces a core bottleneck when attempting to scale the service horizontally and improves performance.

* Shortened URLs are case sensitive.
  * URLs are intended to be machine readable, but also act as a form of [user interface](https://www.w3.org/Provider/Style/URI) interacted with directly by humans.
    * URL case sensitivity decreases the ergonomics of the URL for humans, but increases the number of encodable URLs significantly (36^10 vs 62^10).
      * A shortened URL is already opaque and unergonomic, thus it seems like an acceptable tradeoff given the benefit.

* Are shortened URLs actually "secret"?
  * "yes"
  * Identifiers are randomly generated and are are always a full 10 characters.
  * 10 alphanumeric characters cases sensitive is 62^10 permutations
    * This is a sufficiently large number that brute force attacks against the service are not viable.

* A randomly generated identifier will *not* be checked against the list of existing identifiers before attempting to issue it to the API caller.
  * There is a 1 out of 62^10 percent approximate chance (yes, the probability increases as more shortened links are created, but even millions of short links do not meaningfully change things at this order of magnitude) of any given identifier conflicting with an existing identifier.
  * If an identifier does conflict with an existing identifier, it will do so for <1% of requests, meeting our acceptable error rate.
  * Not having to check if an identifier already exists greatly reduces a core bottleneck when attempting to scale the service horizontally and improves performance.

* Why use a `302` instead of a `301`
  * We may want to extend the service to support changing the target URL for a given identifier.
    * A `301` *usually* causes a client browser to cache the content of the response `Location` header.
      * The duration of this caching varies and may cause some users to unexpectedly get directed to the old target URL.
      * The lack of client caching may increase server load but the server responses can be cached easily using standard tools and mechanisms if necessary.

* If The shortener supported multiple serializations, by what mechanism should the client specify the format to the server?
  * The most correct method is probably the `Accept` request header since that is what it exists for, but I'm leaning towards a query parameter.
    * We already provide the link to shorten via the API URL, so why not do both there?
    * In practice, I think it's not unlikely to encounter proxies and load balancers that mutate the Accept header before the request hits the backend service.
      * We don't always directly control the full path between our users and our backend services.
* How should configs be handled?
  * Ideally they should be handled by a robust library like [viper](https://github.com/spf13/viper) to easily allow configs to be provided via environment variables (preferable in kubernetes land) or a number of interchangeable machine friendly text formats like JSON or YAML.
* What about logging?
  * I think it's reasonable to assume that a service such as this will live behind a reverse proxy which will provide adequate logging of web requests.
    * Logging to files or even standard error / standard out is really only acceptable if you take the extra step to configure log levels, which I'm considering beyond the scope for this exercise. Failing to add configurable log levels can have drastic consequences in the real world.
      * While not exactly logging, distributed tracing via something like OpenTelemetry is fantastic for getting deeper insight into application performance.
  