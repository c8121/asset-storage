const loader = {

    DEFAULT_POST_REQUEST_CONTENT_TYPE: "application/json",

    loadedScripts: [],
    loadedTemplates: {},

    loadedObjects: {},

    /**
     * Request data via HTTP GET.
     * Returns a Promise
     */
    get: function (url) {

        return new Promise(function (resolve, reject) {

            const request = new XMLHttpRequest();
            request.onreadystatechange = function () {
                if (this.readyState === 4) {
                    if (this.status === 200) {
                        resolve(this);
                    } else if (this.status !== 0)
                        reject(this.status);
                }
            }

            request.ontimeout = function (e) {
                console.error("timeout: " + url);
                console.error(e);
                reject("Timeout");
            }

            request.onerror = function (e) {
                console.error("error: " + url);
                console.error(e);
                reject("Network Error");
            }

            request.open("GET", url, true);
            request.send();
        });
    },

    /**
     * Request data via HTTP POST.
     * Returns a Promise
     */
    post: function (url, postData, contentType) {

        const self = this;

        return new Promise(function (resolve, reject) {

            const request = new XMLHttpRequest();
            request.onreadystatechange = function () {
                if (this.readyState === 4) {
                    if (this.status === 200) {
                        resolve(this);
                    } else if (this.status !== 0)
                        reject(this);
                }
            }

            request.ontimeout = function (e) {
                console.error("timeout: " + url);
                console.error(e);
                reject("Timeout");
            }

            request.onerror = function (e) {
                console.error("error: " + url);
                console.error(e);
                reject("Network Error");
            }

            request.open("POST", url, true);
            if (contentType)
                request.setRequestHeader("Content-Type", contentType);
            else
                request.setRequestHeader("Content-Type", self.DEFAULT_POST_REQUEST_CONTENT_TYPE);
            if (postData)
                request.send(postData);
            else
                request.send();
        });
    },

    /**
     * Loads a template if it not was loaded already.
     *
     */
    getTemplate: function (url) {

        if (loader.loadedTemplates.hasOwnProperty(url))
            return Promise.resolve(loader.loadedTemplates[url]);

        return loader.get(url).then((response) => {
            loader.loadedTemplates[url] = response.responseText;
            return response.responseText;
        });
    },

    /**
     * Loads a JavaScript-Object from given URL.
     */
    getObject: function (url) {

        if (loader.loadedObjects.hasOwnProperty(url))
            return Promise.resolve(loader.loadedObjects[url]);

        return loader.get(url).then((response) => {

            if (loader.head == null)
                loader.head = document.getElementsByTagName('head');

            const script = document.createElement('script');
            script.text = "loader.loadedObjects['" + url + "'] = " + response.responseText;
            try {
                loader.head[0].appendChild(script);
                return loader.loadedObjects[url];
            } catch (e) {
                console.log(e);
            }
        });
    }
};