package metadata_db_entity

type Owner struct {
	Id   int64
	Name string `dbe:"Index"`
}

func (o *Owner) GetEntityName() string {
	return "Owner"
}

func (o *Owner) GetId() int64 {
	return o.Id
}

func (o *Owner) SetId(id int64) {
	o.Id = id
}

var (
	OwnerCache map[string]*Owner
)

func init() {
	OwnerCache = make(map[string]*Owner)
}

func LoadOwner(name string, createIfNotExists bool) (*Owner, error) {

	/*owner, ok := OwnerCache[name]
	if ok {
		return owner, nil
	}

	owner = &Owner{Name: name}
	err := db_entity.LoadEntity(createIfNotExists, owner)
	if err == nil {
		OwnerCache[name] = owner
	}

	return owner, err*/
	panic("implement me")
}
