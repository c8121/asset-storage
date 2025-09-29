(function () {
    return {

        template: `
            <div :class="cssClass">
                <div @dragover="dragOver" @drop="drop" :class="innerCssClass">
                    <label class="btn btn-secondary" role="button" for="formFileMultiple" v-html="labelCaption"></label>
                    <input class="d-none" type="file" @change="fileChanged" id="formFileMultiple" multiple>
                </div>
            </div>`,

        props: {
            labelCaption: {
                type: String,
                default: "Datei hochladen oder hierher ziehen"
            },
            cssClass: {
                type: String,
                default: "m-2"
            },
            innerCssClass: {
                type: String,
                default: "p-5 border border-primary-subtle bg-body-tertiary"
            }
        },

        data() {
            return {}
        },
        methods: {
            fileChanged(e) {
                const self = this;
                const files = e.target.files || e.dataTransfer.files;
                if (!files.length)
                    return;

                for (const file of files) {

                    console.log(file)

                    const request = new XMLHttpRequest();
                    request.open('POST', '/assets/upload');
                    request.onprogress = function() {
                        console.log("Upload progress: " + arguments)
                        console.log(arguments)
                    }
                    request.onreadystatechange = function () {
                        if (request.readyState === 4) {
                            if (request.status === 200) {
                                
                                const json = JSON.parse(request.responseText)
                                console.log(json)

                                self.addUploadedFile(json, file)

                            } else if (this.status !== 0)
                                console.error("Upload failed", request)
                        }
                    }
                    request.send(file);
                }
            },
            addUploadedFile(json, file) {
                const query = { 
                    TempName: json.tempName,
                    Name: file.name,
                    Owner: "spa", //TODO
                    FileTime: new Date(file.lastModified).toJSON()
                }
                client.post("/assets/upload/add", query).then((json) => {
                    console.log(json);
                })
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