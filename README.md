# Rimegate

Rimegate is a Grafana dashboard rendering and caching proxy, with a custom-built front-end for displaying Grafana dashboards on devices with severely outdated browsers.

![Rimegate Interface](https://images.ebornet.com/uploads/big/c5658c997f85075faa4f40e0a5339299.png)

## Why?

I bought some cheap, early-generation iPads off eBay, but unfortunately these old versions of iOS have neither support for most recent versions of mobile Chrome / Firefox, nor Safari versions supporting up-to-date JavaScript standards.

These still work (insecurely) for most of the internet, but sadly not the recent version Grafana I'm using. Instead, Grafana refuses to load any graphs, or any UI at all; and instead throws many angry AngularJS errors.

While Grafana already has a built-in [Render API](https://grafana.com/docs/grafana/latest/administration/image_rendering/), remote rendering is computationally expensive, so some form of caching is required to display the same renders of dashboards on multiple devices efficiently. Therefore I wrote Rimegate as a Go-based microservice with a custom API interface proxying Grafana API requests, and an HTML/JavaScript frontend that is **deliberating using no modern browser features to provide automatic refreshes and ensure broad browser compatibility.**

## How to use it

Rimegate is half-finished in terms of providing for general use, as there is still a significant amount of configurations specific to my Kubernetes cluster in code. But in case you are interested in trying it:

* Make sure your Grafana is set up with [remote rendering](https://grafana.com/docs/grafana/latest/administration/image_rendering/#remote-rendering-service), which is essentially an officially-packaged headless Chromium serving requests from the main Grafana instance
* Fork this repository
* Adapt [backend.yaml](http://github.com/icydoge/rimegate/blob/master/backend.yaml), update `GRAFANA_HOST` if your Grafana runs on a different hostname and port (in my case it runs in Kubernetes and has a cluster DNS name)
* Deploy [backend.yaml](http://github.com/icydoge/rimegate/blob/master/backend.yaml) to spawn the backend Rimegate service in Kubernetes
* Expose the backend service somehow through a load balancer on the edge of your cluster (Rimegate proxies Grafana basic auth credentials for authentication, it is otherwise unprivileged when exposed)
* Adapt [web/rimegate/main.js](https://github.com/icydoge/rimegate/blob/master/web/rimegate/main.js) and update `apiBase` to be your Rimegate API base URL exposed above
* Adapt [Makefile](http://github.com/icydoge/rimegate/blob/master/Makefile) and update `WEB_REPOSITORY` to be a Docker repository you can push to and pull from
* Run `make web` to build and push the frontend image
* With the image built and pushed, update the image in [web/frontend.yaml](https://github.com/icydoge/rimegate/blob/master/web/frontend.yaml) and change the service port as necessary
* Deploy [web/frontend.yaml](https://github.com/icydoge/rimegate/blob/master/web/frontend.yaml) to spawn the front Rimegate web service in Kubernetes
* Expose the frontend service somehow through a load balancer on the edge of your cluster (image is unprivileged)

Now visit the web interface, and login with a pair of read-only Grafana credentials. You can then select a dashboard to be rendered. Once in the rendered view, the page will automatically request re-renders periodically, with intervals defined in [web/rimegate/main.js](https://github.com/icydoge/rimegate/blob/master/web/rimegate/main.js).