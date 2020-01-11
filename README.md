# Rimegate

Rimegate is a Grafana dashboard rendering and caching proxy, with a custom-built front-end for displaying Grafana dashboards on devices with severely outdated browsers.

![Rimegate Interface](https://images.ebornet.com/uploads/big/c5658c997f85075faa4f40e0a5339299.png)

## Why?

I bought some cheap, early-generation iPads off eBay, but unfortunately these old versions of iOS have neither support for most recent versions of mobile Chrome / Firefox, nor Safari versions supporting up-to-date JavaScript standards.

These still work (with some unpatched browser vulnerabilities) for most of the internet, but sadly not the recent version of Grafana I'm using. Instead, Grafana refuses to load any graphs, or sometimes any UI at all; and instead throws many angry AngularJS errors.

While Grafana already has a built-in [Render API](https://grafana.com/docs/grafana/latest/administration/image_rendering/), remote rendering is computationally expensive, so some form of caching is required to display the same renders of dashboards on multiple devices efficiently.

As a result, I wrote Rimegate as a Go microservice, serving a custom API interface repackaging and proxying Grafana API requests; along with a basic HTML/JavaScript frontend that is **deliberatly using no modern browser JavaScript support** to ensure broad browser compatibility even on the most outdated browsers (e.g. Mobile Safari running on iOS 5, which is the cap for first-generation iPads).

## How to use it

Rimegate is only semi-ready in terms of providing for general use, and may be particularly difficult to deploy outside Docker/Kubernetes. If you are interested in trying it, however:

* Make sure your Grafana is set up with [remote rendering](https://grafana.com/docs/grafana/latest/administration/image_rendering/#remote-rendering-service), which is essentially an officially-packaged headless Chromium serving requests from the main Grafana instance
* Clone this repository
* In [backend.yaml](http://github.com/icydoge/rimegate/blob/master/backend.yaml), update `GRAFANA_HOST` if your Grafana runs on a different hostname and port (in my case it runs in Kubernetes and has a cluster DNS name)
* Deploy [backend.yaml](http://github.com/icydoge/rimegate/blob/master/backend.yaml) to spawn the backend Rimegate service in Kubernetes
* Expose the backend service somehow through a load balancer on the edge of your cluster (Rimegate proxies Grafana basic auth credentials for authentication, it is therefore unprivileged when exposed)
* Update the `API_BASE` in [web/frontend.yaml](https://github.com/icydoge/rimegate/blob/master/web/frontend.yaml) to be your Rimegate backend API host.
* Deploy [web/frontend.yaml](https://github.com/icydoge/rimegate/blob/master/web/frontend.yaml) to spawn the front Rimegate web service in Kubernetes
* Expose the frontend service somehow through a load balancer on the edge of your cluster (the frontend is unprivileged)

Now visit the web interface, and log in with a pair of read-only Grafana credentials. You can then select a dashboard to be rendered. Once in the rendered view, the page will automatically request re-renders periodically, with intervals defined in [web/rimegate/main.js](https://github.com/icydoge/rimegate/blob/master/web/rimegate/main.js).
