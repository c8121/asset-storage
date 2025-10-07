(function () {
    return {

        template: `
            <ul class="pathitem-tree">
                <li v-for="item in list" class="text-nowrap">
                    
                    <div>
                        <button @click="item.showChildren = !item.showChildren">
                            <span v-if="!item.showChildren">(+)</span>
                            <span v-else>(-)</span>
                        </button>
                        <button @click="itemClicked(item)">
                            <span v-if="!item.Name">({{ item.Id }})</span>{{ item.Name }}
                        </button>
                    </div>
                    
                    <PathItemTree ref="childTrees" v-if="item.showChildren" :value="item.Id"
                            @click="itemClicked"></PathItemTree>
                </li>
            </ul>`,

        props: {
            value: "",
        },

        data() {
            return {
                parentId: 0,
                list: []
            }
        },
        methods: {
            loadItemList(parentId) {

                const self = this;
                self.parentId = parentId;

                const showChildren = [];
                for (const item of self.list) {
                    if (item.showChildren)
                        showChildren.push(item.Id);
                }

                client.get('/pathitems/list/' + self.parentId).then((json) => {
                    self.list.splice(0, self.list.length);
                    for (const item of json) {
                        item.showChildren = false;
                        self.list.push(item);
                    }

                    //Update showChildren after populating list to trigger reload
                    Vue.nextTick(() => {
                        for (const item of self.list)
                            item.showChildren = showChildren.includes(item.Id);
                    });
                });
            },
            treeItemChanged(item) {
                const self = this;

                for (const i of self.list) {
                    if (i.id === item.Id) {
                        self.loadItemList(self.parentId);
                        return true;
                    } else if (i.id === item.parent) {
                        i.showChildren = true;
                        self.loadItemList(self.parentId);
                        return true;
                    }
                }

                if (self.$refs.childTrees) {
                    for (const childTree of self.$refs.childTrees) {
                        if (childTree.treeItemChanged(item))
                            return true;
                    }
                }
                return false;
            },
            itemClicked(item) {
                item.showChildren = true;
                this.$emit('click', item);
            }
        },
        emits: [
            'click'
        ],

        mounted() {
            if (this.value) {
                this.loadItemList(this.value);
            } else {
                this.loadItemList(null);
            }
        },

        componentLoaded() {
            //Use this component recursively in template:
            this.components = {'PathItemTree': this};
        }
    }
})();