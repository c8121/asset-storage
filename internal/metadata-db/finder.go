package metadata_db

type Finder interface {
	//Find searches the metadata database
	Find(query any) (ScoredIdMap, error)
}
