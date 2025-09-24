/**
 * RestWebservice client
 */
const client = {

    initialized: false,
    loadSpinner: null,

    defaultErrorHandler: function (response) {
        if (response && response.error)
            client.showError(response.error);
        else if (response && response.errormessage)
            client.showError(response.errormessage);
        else if (response && response.responseText) {
            const json = JSON.parse(response.responseText);
            if (json.error)
                client.showError(json.error);
            else if (json.errormessage)
                client.showError(json.errormessage);
            else
                client.showError(response.responseText);
        } else
            client.showError("Undefined error (" + (typeof response) + ")");
    },

    progressInfoElement: null,
    messagesElement: null,

    eventListener: [],

    /**
     *
     */
    initialize: function () {
        const self = this;

        if (self.initialized)
            return;
        self.initialized = true;

        self.loadSpinner = document.createElement('div');
        self.loadSpinner.setAttribute('style', 'display: none');
        self.loadSpinner.setAttribute('class', 'loadSpinnerWrapper');

        document.getElementsByTagName('body')[0].appendChild(self.loadSpinner);

        self.progressInfoElement = document.getElementById('clientProgressInfo');
        self.messagesElement = document.getElementById('clientMessages');
    },

    /**
     * Sends a GET-Request to the server and returns the JSON response (Promise-Object).
     * Always expects a JSON-Response.
     * If the response is not JSON, use loader.get(url).
     */
    get: function (url, errorhandler, showProgressInfo) {
        const self = this;

        self.initialize();

        if ((typeof errorhandler == 'undefined') || (errorhandler == null))
            errorhandler = self.defaultErrorHandler;

        if (typeof showProgressInfo == 'undefined')
            showProgressInfo = true;

        if (showProgressInfo)
            self.showProgress("Loading...");

        return loader.get(url).then(
            (response) => {
                const jsonResponse = JSON.parse(response.responseText);

                self.hideProgress();

                if (jsonResponse && jsonResponse.errormessage) {
                    errorhandler(jsonResponse);
                }
                return jsonResponse;
            },
            (errorStatus) => {
                self.hideProgress();
                errorhandler(errorStatus);
            });
    },


    /**
     * Sends a POST-Request to the server and returns the JSON response (Promise-Object).
     * Always expects a JSON-Response.
     * If the response is not JSON, use loader.post(url, data).
     */
    post: function (url, json, errorhandler, showProgressInfo) {
        const self = this;

        self.initialize();

        if ((typeof errorhandler == 'undefined') || (errorhandler == null))
            errorhandler = self.defaultErrorHandler;

        if (typeof showProgressInfo == 'undefined')
            showProgressInfo = true;

        if (showProgressInfo)
            self.showProgress("Loading...");

        return loader.post(url, JSON.stringify(json)).then(
            (response) => {
                const jsonResponse = JSON.parse(response.responseText);

                self.hideProgress();

                if (typeof jsonResponse.errormessage != 'undefined') {
                    errorhandler(jsonResponse);
                }
                return jsonResponse;
            },
            (errorStatus) => {
                self.hideProgress();
                errorhandler(errorStatus);
            });
    },

    /**
     *
     */
    showProgress: function (message) {
        const self = this;

        self.fireProgressEvents();

        if (self.progressInfoElement == null) {
            console.log("PROGRESS: " + message);
            return;
        }

        self.loadSpinner.setAttribute('style', 'display: block');

        while (self.progressInfoElement.firstChild) {
            self.progressInfoElement.firstChild.remove();
        }

        self.progressInfoElement.appendChild(self.createMessageElement(
            message, 'alert alert-info m-2'
        ));
    },

    /**
     *
     */
    hideProgress: function () {
        const self = this;
        if (self.progressInfoElement == null) {
            return;
        }

        while (self.progressInfoElement.firstChild) {
            self.progressInfoElement.firstChild.remove();
        }

        self.loadSpinner.setAttribute('style', 'display: none');
    },


    /**
     *
     */
    showError: function (message) {
        const self = this;
        console.error(message);
        if (self.messagesElement == null) {
            return;
        }

        while (self.messagesElement.firstChild) {
            self.messagesElement.firstChild.remove();
        }

        self.messagesElement.appendChild(self.createMessageElement(
            message, 'alert alert-danger alert-fixed m-2'
        ));
    },

    /**
     *
     */
    hideError: function () {
        const self = this;
        if (self.messagesElement == null) {
            return;
        }

        while (self.messagesElement.firstChild) {
            self.messagesElement.firstChild.remove();
        }
    },

    /**
     *
     */
    createMessageElement(message, cssClasses) {
        const self = this;

        const messageElement = document.createElement('div');
        messageElement.setAttribute('class', cssClasses);
        messageElement.setAttribute('role', 'alert');
        messageElement.innerText = message;

        const buttonElement = document.createElement('button');
        messageElement.appendChild(buttonElement);
        buttonElement.setAttribute('class', 'btn close');
        buttonElement.setAttribute('type', 'button');
        buttonElement.setAttribute('data-dismiss', 'alert');
        buttonElement.setAttribute('aria-label', 'Close');
        buttonElement.innerHTML = '<i class="fas fa-times-circle"></i>';
        buttonElement.onclick = () => self.hideError();

        return messageElement;
    },

    /**
     * Add obj to listen to client events.
     * obj can implement
     *    onClientProgress(...)
     */
    addEventListener: function (obj) {
        //console.trace("add event listener");
        client.eventListener.push(obj);
    },

    /**
     *
     */
    removeEventListener: function (obj) {
        const idx = client.eventListener.indexOf(obj);
        if (idx < 0) {
            console.error("client.removeEventListener: No such listener: " + obj);
            return;
        }

        client.eventListener.splice(idx, 1);
        console.log("client: " + client.eventListener.length + " listeners remaining");
    },

    /**
     *
     */
    fireProgressEvents: function (infoChanges) {
        for (const i in client.eventListener) {
            const listener = client.eventListener[i];
            if (typeof listener.onClientProgress != 'undefined') {
                listener.onClientProgress();
            }
        }
    },
};