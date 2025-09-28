(function () {
    return {

        template: `
            <div :class="cssClass">
                <div @dragover="dragOver" @drop="drop">
                    <label class="btn btn-secondary" role="button" for="formFileMultiple" v-html="labelCaption"></label>
                    <input class="d-none" type="file" @change="fileChanged" id="formFileMultiple" multiple>
                    <button v-if="showUploadButton" @click="upload" class="btn btn-primary">Upload</button>
                </div>
                                
                <div v-if="showItemsList && object && object.items" class="mt-3">
                    <div class="assets row row-cols-auto g-3">
                        <div class="asset col" v-for="asset in object.items">
                            <div class="card bg-light">
                                <img @click="editMetaData(asset)" :src="'/assets/preview?HASH=' + asset.HASH" />
                                <div class="card-body">
                                    {{ asset.meta.name }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
            </div>`,

        props: {
            showItemsList: {
                type: Boolean,
                default: true
            },
            labelCaption: {
                type: String,
                default: "Datei hochladen oder hierher ziehen"
            },
            cssClass: {
                type: String,
                default: "m-2"
            }
        },

        data() {
            return {
                object: {},
                fileNames: [],
                fileData: [],
                showUploadButton: false
            }
        },
        methods: {
            fileChanged(e) {
                const self = this;
                const files = e.target.files || e.dataTransfer.files;
                if (!files.length)
                    return;

                //Remove previously selected files
                self.fileNames.splice(0, self.fileNames.length);
                self.fileData.splice(0, self.fileData.length);

                self.showUploadButton = true;

                for (const file of files) {
                    const reader = new FileReader();
                    reader.onload = function () {
                        self.fileNames.push(file.name);
                        self.fileData.push(this.result);
                        if (self.fileData.length === files.length)
                            self.upload();
                    }
                    reader.readAsDataURL(file);
                }
            },
            upload() {

                const self = this;

                if (!self.fileData)
                    return;

                self.showUploadButton = false;

                const query = {
                    names: self.fileNames,
                    dataUrls: self.fileData
                }
                client.post('/assets/upload', query).then((json) => {
                    self.object = json;
                    self.$emit('uploadFinished', json);
                })
            },
            editMetaData(asset) {
                const self = this;
                const dlg = dialog.create();
                dlg.setTitle("Edit Meta-Data");
                dlg.setConfirmText("Save");
                VueComponentUtil.loadComponent('/vue/components/AssetMetaDataEdit.js', {hash: asset.HASH}).then((component) => {
                    const vm = dlg.setContentComponent(component);
                    vm.showSaveButton(false);
                    dlg.addOnConfirmListener(() => vm.updateMetaData().then((metaData) => {
                        asset.meta = metaData;
                    }));
                });
            },
            dragOver(e) {
                e.preventDefault();
            },
            drop(e) {
                e.preventDefault();
                this.fileChanged(e);
            }
        },
        emits: [
            'uploadFinished'
        ],
    }
})();