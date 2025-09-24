/**
 * Utility class to load components
 */
const VueComponentUtil = {

    /**
     * Load component-object from given URL.
     * Load resources (subcomponents, template)
     */
    loadComponent: function (url, initArgs) {

        const self = this;

        return loader.getObject(url).then((component) => {
            return self.loadParent(component).then((component) => {

                component = Object.assign({}, component); //Copy component, it will be reused.
                if (initArgs && component.data) {
                    const data = component.data();
                    for (const [k, v] of Object.entries(initArgs))
                        data[k] = v;
                    component.data = function () {
                        return data
                    }
                }

                return self.loadSubComponents(component).then((component) => {
                    return self.loadTemplate(component).then((component) => {
                        if (component.componentLoaded)
                            component.componentLoaded();
                        return component;
                    });
                });
            });
        });
    },

    /**
     * Load resources (subcomponents, template)
     */
    loadComponentResources(component) {

        const self = this;

        return self.loadSubComponents(component).then((component) => {
            return self.loadTemplate(component);
        });
    },

    /**
     * Check if given component contains member 'components' with URLs.
     * Load objects from URLs.
     */
    loadSubComponents: function (component) {

        if (component.components) {

            const subComponentLoaders = [];
            for (const name in component.components) {
                if (typeof component.components[name] === 'string') { //check if it is a string, otherwise it's a component already
                    const url = component.components[name];
                    subComponentLoaders.push(
                        VueComponentUtil.loadComponent(url).then((subComponent) => {
                            console.log("Loaded sub-component '" + name + "' from " + url);
                            return {name, component: subComponent};
                        })
                    )
                } else {
                    subComponentLoaders.push({name, component: component.components[name]});
                }
            }
            return Promise.all(subComponentLoaders).then((subComponents) => {
                component.components = {};
                for (subComponent of subComponents)
                    component.components[subComponent.name] = subComponent.component;
                return component;
            });
        } else {
            return Promise.resolve(component);
        }
    },

    loadParent: function (component) {
        if (component.parentUrl) {
            const url = component.parentUrl;

            return loader.getObject(url).then((parent) => {
                let extended = Object.assign({}, parent);
                extended = Object.assign(extended, component);

                if (parent.methods) {
                    const methods = Object.assign({}, parent.methods);
                    if (extended.methods)
                        Object.assign(methods, extended.methods);
                    extended.methods = methods;
                }

                if (parent.data) {
                    const merged = Object.assign(extended.data ? extended.data() : {}, parent.data());
                    extended.data = () => {
                        return merged
                    };
                }

                return extended;
            });
        } else {
            return Promise.resolve(component);
        }
    },

    /**
     * Check if given component contains member 'templateUrl' with URLs.
     * Load template from URL an assign it to member 'template'.
     */
    loadTemplate: function (component) {

        if (component.templateUrl) {
            const url = component.templateUrl;
            return loader.getTemplate(url).then((template) => {
                console.log("Loaded template from " + url);
                component.template = template;
                return component;
            });
        } else {
            return Promise.resolve(component);
        }
    }
}