import AssetListView from './asset-list.vue.js'
import AssetView from './asset-view.vue.js'
import AssetSelectionListView from "./asset-selection-list.vue.js";
import PathTree from "./path-tree.vue.js";
import Upload from "./upload.vue.js";
import CollectionListView from './collection-list.vue.js'
import CollectionView from './collection-view.vue.js'

export default {

    components: {
        'asset-list-view': AssetListView,
        'asset-view': AssetView,
        'asset-selection-list': AssetSelectionListView,
        'path-tree': PathTree,
        'upload': Upload,
        'collection-view': CollectionView,
        'collection-list-view': CollectionListView
    },

    data() {
        return {
            listFilter: {
                MimeType: null,
                FileName: null
            },

            activeAsset: null,
            activeAssetToast: null,
            selectedAssets: {},
            selectedAssetsToast: null,
            activeCollection: null,
            activeCollectionToast: null,

            mimeTypes: [],

            pathIdFilter: '',

            notifications: [],

            downloadCallback: (url) => {window.open(url);}
        }
    },

    methods: {

        componentEvent(eventName, sender) {
            switch (eventName) {
                case 'assetClick': {
                    this.activeAsset = arguments[2];
                    this.setToastVisible(this.activeAssetToast, this.activeAsset !== null);
                    if(sender === 'asset-view') {
                        const url = '/assets/' + this.activeAsset.Hash;
                        this.downloadCallback(url);
                    }
                    break;
                }
                case 'assetSelect': {
                    const asset = arguments[2];
                    this.selectedAssets[asset.Hash] = asset;
                    this.setToastVisible(this.selectedAssetsToast, Object.keys(this.selectedAssets).length > 0);
                    break;
                }
                case 'assetUnselect': {
                    const asset = arguments[2];
                    delete this.selectedAssets[asset.Hash];
                    this.setToastVisible(this.selectedAssetsToast, Object.keys(this.selectedAssets).length > 0);
                    break;
                }
                case 'assetDownloadClick': {
                    const asset = arguments[2];
                    const url = '/assets/' + asset.Hash;
                    this.downloadCallback(url);
                    break;
                }
                case 'assetFilteredDownloadClick': {
                    const req = arguments[2];
                    const url = '/assets/filter/' + this.createDownloadFilterReuqestUrl(req) + "/" + req.asset.Hash;
                    this.downloadCallback(url);
                    break;
                }
                case 'assetPathClick': {
                    const path = arguments[2];
                    this.filterPath(path);
                    break;
                }
                case 'assetFaceClick': {
                    const faceId = arguments[2];
                    this.resetListFilter();
                    this.listFilter.Face = faceId;
                    this.$refs.listView.loadAssets();
                    break;
                    break;
                }
                case 'pathItemClick': {
                    const path = arguments[2];
                    this.resetListFilter();
                    this.listFilter.PathId = path.Id;
                    this.pathIdFilter = path.Name;
                    this.$refs.listView.loadAssets();
                    break;
                }
                case 'collectionClick': {
                    const collection = arguments[2];
                    this.resetListFilter();
                    this.listFilter.Collection = collection.UUID;
                    this.$refs.listView.loadAssets();
                    this.activeCollection = collection;
                    this.setToastVisible(this.activeCollectionToast, true)
                    break;
                }
                case 'collectionSaved': {
                    this.$refs.collectionList.loadCollections();
                    break;
                }
                case 'uploadFinished': {
                    this.resetListFilter();
                    this.$refs.listView.loadAssets();
                    break;
                }
                case 'createCollectionClick': {
                    const name = arguments[2];
                    this.createCollection(name);
                    break;
                }
                default: {
                    console.log("Event not handled: '" + eventName + "', " + arguments.length + " arguments, sender: " + sender);
                    console.log(arguments);
                }
            }
        },

        resetListFilter() {
            this.listFilter.FileName = null;
            this.listFilter.MimeType = null;
            this.listFilter.PathId = null;
            this.listFilter.Face = null;
            this.listFilter.Collection = null;
            this.listFilter.Fulltext = null;

            this.pathIdFilter = null;
        },

        setToastVisible(el, visible) {
            if (!el) return;
            if (visible) {
                el.classList.add('show');
                el.classList.remove('hide');
            } else {
                el.classList.add('hide');
                el.classList.remove('show');
            }
        },

        loadMimeTypes() {
            const self = this;
            const requestOptions = {
                method: 'GET'
            }
            fetch('/mimetypes/list', requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.mimeTypes.splice(0, self.mimeTypes.length);
                    for (const item of json) {
                        self.mimeTypes.push(item);
                    }
                });
        },

        filterPath(path) {
            const self = this;
            const requestOptions = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: "path=" + encodeURI(path)
            }
            fetch('/pathitem/', requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.resetListFilter();
                    self.listFilter.PathId = json.Id;
                    self.pathIdFilter = path;
                    self.$refs.listView.loadAssets();
                });
        },

        loadNotifications() {
            const self = this;
            const requestOptions = {
                method: 'GET'
            }
            fetch('/notifications', requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.notifications.splice(0, self.notifications.length);
                    for (const msg of json) {
                        self.notifications.push(msg);
                    }
                });
        },

        toggleToastPosition() {
            const el = document.getElementById("toasts");
            if(el) {
                if(el.classList.contains("end-0")) {
                    el.classList.remove("end-0");
                    el.classList.add("start-0")
                } else {
                    el.classList.remove("start-0")
                    el.classList.add("end-0");
                }
            }
        },

        createCollection(name) {
            const self = this;

            const hashes = [];
            for( const hash of Object.keys(self.selectedAssets)) {
                hashes.push(hash);
            }

            if(!name || name.trim() === "") {
                name = "Collection " + new Date()
            }

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    Name: name,
                    Owner: "spa",
                    Description: "Collection created at " + new Date(),
                    Assets: hashes
                })
            }

            fetch('/collections/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    self.selectedAssets = {};
                    self.setToastVisible(this.selectedAssetsToast, false);
                    self.$refs.collectionList.loadCollections();

                    self.activeCollection = json;
                    self.setToastVisible(this.activeCollectionToast, true);

                    self.resetListFilter();
                    self.listFilter.Collection = json.UUID;
                    self.$refs.listView.loadAssets();
                });

        },

        removeCollection(collection) {
            const self = this;
            if(!confirm("Delete collection?"))
                return

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: "UUID=" + encodeURI(collection.UUID)
            }

            fetch('/collections/remove', requestParams)
                .then(res => res.json())
                .then(json => {
                    console.log(json);
                    self.activeCollection = null;
                    self.setToastVisible(this.activeCollectionToast, false);
                    self.$refs.collectionList.loadCollections();

                    self.resetListFilter();
                    self.$refs.listView.loadAssets();
                });
        },

        createDownloadFilterReuqestUrl(request) {
            let paramsString = "filterName:" + request.filter;
            if(request.params) {
                for(const param of request.params) {
                    const v = ("" + param.value).trim()
                    if(param.value !== "") {
                        paramsString += "|" + param.name + ":" + param.value
                    }
                }
            }
            paramsString = encodeURI(paramsString);
            return paramsString;
        },
    },

    mounted() {

        console.log("Attributes", this.$attrs);
        if(this.$attrs.downloadCallback) {
            this.downloadCallback = this.$attrs.downloadCallback;
            console.log("Replace downloadCallback", this.downloadCallback);
        }

        this.activeAssetToast = this.$refs["activeAssetToast"];
        this.setToastVisible(this.activeAssetToast, false);
        this.selectedAssetsToast = this.$refs["selectedAssetsToast"];
        this.setToastVisible(this.selectedAssetsToast, false);
        this.activeCollectionToast = this.$refs["activeCollectionToast"];
        this.setToastVisible(this.activeCollectionToast, false);

        this.loadMimeTypes();
        this.loadNotifications();
    }
}