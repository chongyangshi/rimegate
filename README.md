# Rimegate

Rimegate is a Grafana dashboard rendering and caching proxy, with a custom-built front-end for displaying Grafana dashboards on devices with severely outdated browsers.

![Rimegate Interface](https://i.doge.at/uploads/big/c5658c997f85075faa4f40e0a5339299.png)

## Why?

I bought some cheap, early-generation iPads off eBay, but unfortunately these old versions of iOS have neither support for most recent versions of mobile Chrome / Firefox, nor Safari versions supporting up-to-date JavaScript standards.

These still work (with some unpatched browser vulnerabilities) for most of the internet, but sadly not the recent version of Grafana I'm using. Instead, Grafana refuses to load any graphs, or sometimes any UI at all; and instead throws many angry AngularJS errors.

While Grafana already has a built-in [Render API](https://grafana.com/docs/grafana/latest/administration/image_rendering/), remote rendering is computationally expensive, so some form of caching is required to display the same renders of dashboards on multiple devices efficiently.

As a result, I wrote Rimegate as a Go microservice, serving a custom API interface repackaging and proxying Grafana API requests; along with a basic HTML/JavaScript frontend that is **deliberatly using no modern browser JavaScript support** to ensure broad browser compatibility even on the most outdated browsers (e.g. Mobile Safari running on iOS 5, which is the cap for first-generation iPads).

## How to use it

Rimegate is only semi-ready in terms of providing for general use, and may be particularly difficult to deploy outside Docker/Kubernetes. If you are interested in trying it, please follow steps below.

### Decide how to authenticate clients

There are two options you can authenticate clients, which require different credential setups:
* **Grafana username/password**: Rimegate will pass on Grafana credentials supplied by the client to Grafana for authentication, these will need to be supplied by the user in a login page served by the front-end. 
    * Set up a `Viewer` role user in Grafana admin;
    * **If using Kubernetes**, in [backend.yaml](http://github.com/chongyangshi/rimegate/blob/master/backend.yaml), remove the environment variable `GRAFANA_API_TOKEN` in the container `env`.
* **Grafana static API token + mTLS**: Rimegate will not attempt to authenticate clients, and instead rely on a static Grafana API token to authenticate itself with Grafana.
    * Set up a `Viewer` role api key in Grafana admin, with a sufficiently long duration of validity;
    * Set the `GRAFANA_API_TOKEN` environment variable in Rimegate's execution environment to be this static API key, if you are using Kubernetes, this can instead be supplied in a Kubernetes Secret called `rimegate-grafana-api-token` with key name also `rimegate-grafana-api-token`. You can change the secret/key name in [backend.yaml](http://github.com/chongyangshi/rimegate/blob/master/backend.yaml);
    * Now set up mTLS. The remaining steps are optional, but without them, anyone will be able to use the Rimegate backend as a proxy to query your Grafana, which is especially dangerous if your api key is not read-only;
    * Run [`./p12-ca.sh`](http://github.com/chongyangshi/rimegate/blob/master/p12-ca.sh) to generate a self-signed root CA for Rimegate; upload the root CA certificate (but not the key) to your revese proxy web server;
    * Run [`./p12-client.sh <client-name>`](http://github.com/chongyangshi/rimegate/blob/master/p12-ca.sh) to generate a client certificate for the given `<client-name>`; transfer the resulting `.p12` file to your Rimegate display device by email or a web server, and install it;
    * In your reverse proxy (such as NGINX) serving Rimegate front-end, set up TLS client authentication with the root CA above.

### Deploy Rimegate

Once the appropriate authentication method above has been set up, you are now ready to deploy Rimegate:

* Make sure your Grafana is set up with [remote rendering](https://grafana.com/docs/grafana/latest/administration/image_rendering/#remote-rendering-service), which is essentially an officially-packaged headless Chromium serving requests from the main Grafana instance
* Clone this repository
* In [backend.yaml](http://github.com/chongyangshi/rimegate/blob/master/backend.yaml), update `GRAFANA_HOST` if your Grafana runs on a different hostname and port (in my case it runs in Kubernetes and has a cluster DNS name)
* Deploy [backend.yaml](http://github.com/chongyangshi/rimegate/blob/master/backend.yaml) to spawn the backend Rimegate service in Kubernetes
* Expose the backend service somehow through a load balancer on the edge of your cluster (Rimegate proxies Grafana basic auth credentials for authentication, it is therefore unprivileged when exposed)
* Update the `API_BASE` in [web/frontend.yaml](https://github.com/chongyangshi/rimegate/blob/master/web/frontend.yaml) to be your Rimegate backend API host.
* Deploy [web/frontend.yaml](https://github.com/chongyangshi/rimegate/blob/master/web/frontend.yaml) to spawn the front Rimegate web service in Kubernetes
* Expose the frontend service somehow through a load balancer on the edge of your cluster (the frontend is unprivileged)

Now visit the web interface, and log in with a pair of read-only Grafana credentials. You can then select a dashboard to be rendered. Once in the rendered view, the page will automatically request re-renders periodically, with intervals defined in [web/rimegate/main.js](https://github.com/chongyangshi/rimegate/blob/master/web/rimegate/main.js).
