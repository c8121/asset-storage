import { nextTick } from 'vue'

export default {
    template: `
        <div>
            <div v-for="item in list" class="text-nowrap text-start">
                
                <div>
                    <button @click="item.showChildren = !item.showChildren" class="btn m-0 p-0">
                        <span v-if="!item.showChildren"><i class="bi bi-folder-plus"></i></span>
                        <span v-else><i class="bi bi-folder-minus"></i></span>
                    </button>
                    <button @click="item.showChildren = true; itemClicked(item)" class="btn m-0 p-0">
                        <span v-if="!item.Name">({{ item.Id }})</span>{{ item.Name }}
                    </button>
                </div>
                
                <div class="ms-2">
                <path-tree ref="childTrees" v-if="item.showChildren" :value="item.Id"
                        @componentEvent="childClicked"></path-tree>
                </div>
            </div>
        </div>
    `,

    props: {
        value: "",
    },

    data() {
        return {
            parentId: 0,
            list: []
        }
    },

    watch: {
        value() {
            this.loadItems();
        }
    },

    methods: {

        loadItems() {
            const self = this;

            const showChildren = [];
            for (const item of self.list) {
                if (item.showChildren)
                    showChildren.push(item.Id);
            }

            const requestOptions = {
                method: 'GET'
            }
            fetch('/pathitems/list/' + self.value, requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.list.splice(0, self.list.length);
                    for (const item of json) {
                        item.showChildren = false;
                        self.list.push(item);
                    }

                    //Update showChildren after populating list to trigger reload
                    nextTick(() => {
                        for (const item of self.list)
                            item.showChildren = showChildren.includes(item.Id);
                    });
                });
        },

        itemClicked(item) {
            this.$emit('componentEvent', 'pathItemClick', 'path-tree', item);
        },

        childClicked() {
            this.$emit('componentEvent', 'pathItemClick', 'path-tree', arguments[2]);
        }
    },

    emits: ['componentEvent'],

    created() {
        this.loadItems();
    }
}