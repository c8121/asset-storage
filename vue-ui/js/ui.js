/**
 * UI Utils
 */
const ui = {

    vueMixin: {
        methods: {
            confirmDialog(title, message, confirmCallback, rejectCallback) {
                ui.confirm(title, message, confirmCallback, rejectCallback)
            },
            arraySortDialog(ary, itemToStringFunction, title) {
                ui.arraySortDialog(ary, itemToStringFunction, title)
            },
            objecSelectDialog(title, widgetUrl, widgetSelectMethodName, selectCallback) {
                ui.objecSelectDialog(title, widgetUrl, widgetSelectMethodName, selectCallback)
            },
            personSelectDialog(selectCallback) {
                ui.objecSelectDialog('Auswählen', '/vue/components/contact/PersonListWidget.js', 'rowClicked', (row) => {
                    selectCallback(row)
                });
            }
        }
    },


    /**
     * Confirmation Dialog
     */
    confirm(title, message, confirmCallback, rejectCallback) {

        const dlg = dialog.create();
        dlg.setTitle(title);
        dlg.setContent("<p>" + message + "</p>");

        if (rejectCallback)
            dlg.addOnDismissListener(rejectCallback);
        dlg.setDismissText("No");

        if (confirmCallback)
            dlg.addOnConfirmListener(confirmCallback);
        dlg.setConfirmText("Yes");

    },


    /**
     * Array-Sorting Dialog
     */
    arraySortDialog(ary, itemToStringFunction, title) {

        if (!itemToStringFunction) {

            itemToStringFunction = (item) => {
                if (item.Name)
                    return item.Name;
                if (item.Title)
                    return item.Title;
                if (item.Text)
                    return item.Text;
                return item;
            }
        }

        if (!title)
            title = "Sortieren";

        const dlg = dialog.create();
        dlg.setTitle(title)
        dlg.setDismissText("");
        VueComponentUtil.loadComponent('/vue/components/common/ArraySorter.js').then((component) => {
            const vm = dlg.setContentComponent(component, {});
            vm.setItemToStringFunction(itemToStringFunction);
            vm.setArray(ary);
        });
    },

    /**
     * Object-Select Dialog
     */
    objecSelectDialog(title, widgetUrl, widgetSelectMethodName, selectCallback) {

        if (!title)
            title = "Auswählen";

        const dlg = dialog.create();
        dlg.setTitle(title);
        dlg.setConfirmText("");
        dlg.setFitContent();
        VueComponentUtil.loadComponent('/vue/components/common/WidgetContainerWithTaskbar.js').then((component) => {
            const vm = dlg.setContentComponent(component, {});
            vm.setControls(false, false, false, true);
            vm.setWidget(widgetUrl, {
                    [widgetSelectMethodName]: function (args) {
                        selectCallback(args);
                        dlg.close();
                    }
                }
            );
        });
    }
};