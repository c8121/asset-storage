(function () {
    return {

        template: `
            <div>
                Template not implemented.
            </div>
            `,

        data() {
            return {}
        },
        methods: {
            loadDataObject(id) {
                console.log("Not implemented: loadDataObject(" + id + ")");
            },
            /*saveDataObject() {
                console.log("Not implemented: saveDataObject()");
            },
            deleteDataObject() {
                console.log("Not implemented: deleteDataObject()");
            },
            addNewDataObject() {
                console.log("Not implemented: addNewDataObject()");
            }*/
        },

        created() {
            const self = this;
            self.$watch(
                () => this.$route.params,
                (params) => {
                    this.loadDataObject(params.id && params.id !== "null" ? Number(params.id) : null);
                }
            )
        },

        mounted() {
            const params = this.$route.params;
            this.loadDataObject(params.id && params.id !== "null" ? Number(params.id) : null);
        }
    }
})();